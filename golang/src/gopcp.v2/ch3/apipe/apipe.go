package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

/****************************************************

说明:
	1. 模拟管道方式执行命令.
		- 前一个命令执行结果为后一个的输入
	2. 重构了原操作代码, 注意传参格式.
	3. 指针型形参, 传入实参时, 要取地址&
	4. 参数传递:
		- 注意区别 指针型和 传值型, 后者为值复制.

****************************************************/

func main() {
	runCmd()
	fmt.Println("----------------------------------")
	runCmdWithPipe()

}

func runCmd() {
	useBufferIO := false
	fmt.Println("Run command `echo -n \"My first command comes from golang.\"`: ")
	cmd0 := exec.Command("echo", "-n", "My first command comes from golang.")

	stdout0, err := cmd0.StdoutPipe()
	if err != nil {
		fmt.Printf("Error: Couldn't obtain the stdout pipe for command No.0: %s\n", err)
		return
	}

	if err := cmd0.Start(); err != nil {
		fmt.Printf("Error: The command No.0 can not be startup: %s\n", err)
		return
	}

	if !useBufferIO { // 不使用 bufio 模块
		var outputBuf0 bytes.Buffer

		// 死循环:
		//     - 每次读取5字节, 自行拼接最终的结果字符串.
		for {
			tempOutput := make([]byte, 5)
			n, err := stdout0.Read(tempOutput) // 从 stdout 读取内容到 tempOutput

			if err != nil {
				if err == io.EOF { // 读取到结束符
					break
				} else {
					fmt.Printf("Error: Couldn't read data from the pipe: %s\n", err)
					return
				}
			}

			if n > 0 { // 读取内容不为空
				outputBuf0.Write(tempOutput[:n])                    // 将读取到的内容, 拼接到最终结果字符串
				fmt.Println("[Mid value]:\t" + outputBuf0.String()) // 自定义拼接过程
			}
		}
		fmt.Printf("result:\n%s\n", outputBuf0.String()) // 输出最终结果
	} else { // 使用 bufio 模块
		outputBuf0 := bufio.NewReader(stdout0)
		output0, _, err := outputBuf0.ReadLine()
		if err != nil {
			fmt.Printf("Error: Coundn't read data from the pipe: %s\n", err)
			return
		}
		fmt.Printf("result:\n%s\n", string(output0)) // 输出最终结果
	}

}

//
// 模拟管道方式执行命令:
// 	   - 前一条命令的执行结果, 为后一条命令的输入
//	   - 串行执行
//
func runCmdWithPipe() {
	fmt.Println("Run command `ps aux | grep apipe`: ")

	// 命令初始化:
	cmd1 := exec.Command("ps", "aux")
	cmd2 := exec.Command("grep", "apipe")

	// 命令1执行:
	var outputBuf1 bytes.Buffer
	do(cmd1, &outputBuf1) // 执行命令, 并将命令结果写到 outputbuf

	cmd2.Stdin = &outputBuf1 // 更改cmd2命令的 stdin 输入源: 将cmd1命令执行结果, 当成 cmd2 的输入

	// 命令2执行:
	var outputBuf2 bytes.Buffer
	do(cmd2, &outputBuf2) // 执行命令, 并将命令结果写到 outputbuf

	// 输出命令2执行结果:
	fmt.Printf("result:\n%s\n", outputBuf2.Bytes())

}

//
// 执行命令:
//     - 更改默认的stdout输出方式
//	   - 将输出结果, 写出到 outbuf 中
//
func do(cmd *exec.Cmd, outbuf *bytes.Buffer) {
	cmd.Stdout = outbuf // 更改标准输出到: outbuf

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error: the first command can not be startup %s\n", err)
		return
	}

	if err := cmd.Wait(); err != nil {
		fmt.Printf("Error: Couldn't wait for the first command: %s\n", err)
		return
	}

	// 输出执行结果:
	//fmt.Printf("result:\n%s\n", outbuf.Bytes())
}
