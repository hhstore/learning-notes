package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

/****************************************************

说明:
	1. 多进程编程 - IPC 通信 - 基于信号
		- 通过信号+chanel, 控制通信过程
	2. 代码编写思路:
		- 自顶向下写代码.
		- 子模块, 子函数, 先声明为空函数, 返回nil.
		- 把主干逻辑实现好, 再填子逻辑部分.

****************************************************/

var wg sync.WaitGroup

func main() {
	go func() {
		time.Sleep(3 * time.Second)
		sendSignal() // 信号发送
	}()

	// 信号处理
	handleSignal()
}

//
func handleSignal() {
	// 信号1:
	sigRecv1 := make(chan os.Signal, 1)
	sigs1 := []os.Signal{
		syscall.SIGINT,
		syscall.SIGQUIT,
	}

	fmt.Printf("Set notification for %s... [sigRecv1]\n", sigs1)
	signal.Notify(sigRecv1, sigs1...)

	// 信号2:
	sigRecv2 := make(chan os.Signal, 1)
	sigs2 := []os.Signal{
		syscall.SIGQUIT,
	}

	fmt.Printf("Set notification for %s... [sigRecv2]\n", sigs2)
	signal.Notify(sigRecv2, sigs2...)

	//

	wg.Add(2)

	// 例程1:
	go func() {
		for sig := range sigRecv1 {
			fmt.Printf("Received a signal from sigRecv1: %s\n", sig)
		}
		//
		fmt.Printf("End. [sigRecv1]\n")
		wg.Done()
	}()

	// 例程2:
	go func() {
		for sig := range sigRecv2 {
			fmt.Printf("Received a signal from sigRecv1: %s\n", sig)
		}
		//
		fmt.Printf("End. [sigRecv2]\n")
		wg.Done()
	}()

	//
	fmt.Println("Wait for 2 seconds...")
	time.Sleep(2 * time.Second)
	//
	signal.Stop(sigRecv1)
	close(sigRecv1)
	//
	fmt.Printf("Done. [sigRecv1]\n")
	wg.Wait()

}

//
func sendSignal() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Fatal Error: %s\n", err)
			debug.PrintStack()
		}
	}()

	// 命令集
	// ps aux | grep "signal" | grep -v "grep" | grep -v "go run" | awk '{print $2}'
	cmds := []*exec.Cmd{
		exec.Command("ps", "aux"),
		exec.Command("grep", "signal"),
		exec.Command("grep", "-v", "grep"),
		exec.Command("grep", "-v", "go run"),
		exec.Command("awk", "{print $2}"),
	}

	// 命令执行:
	output, err := runCmds(cmds) // 返回值为进程ID列表
	if err != nil {
		fmt.Printf("Command Execution Error: %s\n", err)
		return
	}

	pids, err := getPids(output)
	if err != nil {
		fmt.Printf("PID Parsing Error: %s\n", err)
		return
	}

	fmt.Printf("Target PID(s):\n%v\n", pids)

	for _, pid := range pids {
		proc, err := os.FindProcess(pid)
		if err != nil {
			fmt.Printf("Process Finding Error: %s\n", err)
			return
		}

		sig := syscall.SIGQUIT
		fmt.Printf("Send signal '%s' to the process (pid=%d)...\n", sig, pid)

		err = proc.Signal(sig)
		if err != nil {
			fmt.Printf("Signal Sending Error: %s\n", err)
			return
		}

	}

}

// 获取PIDs
func getPids(strs []string) ([]int, error) {
	var pids []int

	for _, str := range strs {
		// 转换:
		pid, err := strconv.Atoi(strings.TrimSpace(str))
		if err != nil {
			return nil, err
		}

		// 插入到数组
		pids = append(pids, pid)
	}
	//
	return pids, nil
}

// 执行命令:
func runCmds(cmds []*exec.Cmd) ([]string, error) {
	if cmds == nil || len(cmds) == 0 {
		return nil, errors.New("The cmd slice is invalid!")
	}

	first := true
	var output []byte
	var err error

	// 迭代命令集, 逐个执行:
	for _, cmd := range cmds {
		// 执行过程:
		//		- 前一条命令执行的结果, 为后一条命令的输入, (管道方式执行)
		//		- 不断提取
		//		- 不断刷新 stdoutBuf 的内容
		//
		fmt.Printf("[CMD]: %v\n", getCmdPlaintext(cmd))

		if !first {
			// 输入:
			var stdinBuf bytes.Buffer
			stdinBuf.Write(output)
			cmd.Stdin = &stdinBuf
		}

		// 输出:
		var stdoutBuf bytes.Buffer
		cmd.Stdout = &stdoutBuf

		// 执行命令:
		if err = cmd.Start(); err != nil {
			return nil, getError(err, cmd)
		}
		if err = cmd.Wait(); err != nil {
			return nil, getError(err, cmd)
		}

		// 输出结果:
		output = stdoutBuf.Bytes()
		//fmt.Println("CMD Result:\n" + string(output))

		if first {
			first = false
		}
	}

	var lines []string
	var outputBuf bytes.Buffer

	// 数据写入到buf:
	outputBuf.Write(output) // output为长字符串

	// 解析提取 buf 内容:
	for {
		// 根据换行符切割字符串, 再依次存入数组
		line, err := outputBuf.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, getError(err, nil)
		}

		// 存入数组
		lines = append(lines, string(line))
	}

	//
	return lines, nil
}

func getCmdPlaintext(cmd *exec.Cmd) string {
	var buf bytes.Buffer
	buf.WriteString(cmd.Path)

	for _, arg := range cmd.Args[1:] {
		buf.WriteRune(' ')
		buf.WriteString(arg)
	}

	return buf.String()
}

func getError(err error, cmd *exec.Cmd, extraInfo ...string) error {
	var errMsg string

	if cmd != nil {
		errMsg = fmt.Sprintf("%s [%s %v]", err, (*cmd).Path, (*cmd).Args)
	} else {
		errMsg = fmt.Sprintf("%s", err)
	}

	if len(extraInfo) > 0 {
		errMsg = fmt.Sprintf("%s (%v)", errMsg, extraInfo)
	}

	return errors.New(errMsg)
}
