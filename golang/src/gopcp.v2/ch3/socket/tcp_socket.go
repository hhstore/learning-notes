package main

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

/****************************************************

        TCP 通信 : 服务端, 客户端

说明:
	- 基于 go 例程实现:
		- 一个 go server() 例程, 一个 client() 例程.
	- 注意TCP连接处理过程: < 服务端, 客户端>
	- go 例程 调度方式:
		- 利用 sync.WaitGroup

****************************************************/

const (
	// server:
	SERVER_NETWORK = "tcp"
	SERVER_ADDRESS = "127.0.0.1:8005"

	// 分隔符
	DELIMITER = '\t'
)

var wg sync.WaitGroup

func main() {
	wg.Add(2)

	go serverGo()                      // 服务端
	time.Sleep(500 * time.Millisecond) // 等待
	go clientGo(1)                     // 客户端

	wg.Wait()
}

//========================================================

// 服务器端:
func serverGo() {
	var listener net.Listener
	listener, err := net.Listen(SERVER_NETWORK, SERVER_ADDRESS)
	if err != nil {
		printServerLog("Listen Error: %s", err)
		return
	}
	defer listener.Close()

	printServerLog("Got listener for the server. (local address: %s)", listener.Addr())

	// 死循环:
	for {
		conn, err := listener.Accept() // 监听, 阻塞等请求.
		if err != nil {
			printServerLog("Accept Error: %s", err)
		}
		printServerLog("####Established a connection with a client application. (remote address: %s)", conn.RemoteAddr())

		// 启动一个服务例程单独处理:
		go handleConn(conn)
	}

}

// 客户端:
func clientGo(id int) {
	defer wg.Done()

	// 创建TCP连接, 连接服务器端:
	conn, err := net.DialTimeout(SERVER_NETWORK, SERVER_ADDRESS, 2*time.Second)
	if err != nil {
		printClientLog(id, "Dial Error: %s", err)
		return
	}
	defer conn.Close() // 请求结束, 关闭TCP连接

	//
	printClientLog(id, "Connected to server. (remote address: %s, local address: %s",
		conn.RemoteAddr(),
		conn.LocalAddr(),
	)

	time.Sleep(200 * time.Millisecond)                     // 休眠时间
	reqNum := 5                                            // 请求个数
	conn.SetDeadline(time.Now().Add(5 * time.Millisecond)) // 请求超时时间

	// 发送消息:
	for i := 0; i < reqNum; i++ {
		req := rand.Int31()
		//
		n, err := write(conn, fmt.Sprintf("%d", req))
		if err != nil {
			printClientLog(id, "Write Error: %s", err)
			continue
		}
		printClientLog(id, "Sent request (written %d bytes): %d.", n, req)
	}

	// 读取响应:
	for j := 0; j < reqNum; j++ {
		strResp, err := read(conn)
		if err != nil {
			if err == io.EOF {
				printClientLog(id, "The connection is closed by annother side.")
			} else {
				printClientLog(id, "Read Error: %s", err)
			}
			break
		}
		printClientLog(id, "Received response: %s.", strResp)
	}

}

//========================================================

func handleConn(conn net.Conn) {
	defer func() {
		conn.Close()
		wg.Done()
	}()

	// 处理客户端请求, 并返回响应给客户端
	for {
		// 设置连接超时时间:
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))

		//----------------------------------
		//         处理客户端请求
		//----------------------------------

		strReq, err := read(conn)
		if err != nil {
			// 客户端主动关闭TCP连接.
			if err == io.EOF {
				printServerLog("the connection is closed by another side.")
			} else {
				printServerLog("Read Error: %s", err)
			}
			break // 结束运行
		}
		printServerLog("###Received request: %s.", strReq)

		//
		intReq, err := strToInt32(strReq)
		if err != nil {
			//----------------------------------
			//         返回客户端 正常响应
			//----------------------------------
			n, err := write(conn, err.Error())
			printServerLog("send error message (Written %d bytes): %s.", n, err)
			continue
		}
		floatReq := cbrt(intReq)

		// 待返回响应数据:
		respMsg := fmt.Sprintf("the cube root of %d is %f.", intReq, floatReq)

		//----------------------------------
		//         返回客户端 正常响应
		//----------------------------------
		n, err := write(conn, respMsg)
		if err != nil {
			printServerLog("Write Error: %s", err)
		}
		printServerLog("sent response (written %d bytes): %s.", n, respMsg)

	}

}

//========================================================

// 千万不要使用这个版本的read函数！
//func read(conn net.Conn) (string, error) {
//	reader := bufio.NewReader(conn)
//	readBytes, err := reader.ReadBytes(DELIMITER)
//	if err != nil {
//		return "", err
//	}
//	return string(readBytes[:len(readBytes)-1]), nil
//}

func read(conn net.Conn) (string, error) {
	//
	readBytes := make([]byte, 1)
	var buffer bytes.Buffer

	// 逐个字节读取.
	for {
		_, err := conn.Read(readBytes)
		if err != nil {
			return "", err
		}

		readByte := readBytes[0]
		if readByte == DELIMITER { // 分隔符
			break
		}
		buffer.WriteByte(readByte)

	}

	return buffer.String(), nil
}

func write(conn net.Conn, content string) (int, error) {
	var buffer bytes.Buffer

	buffer.WriteString(content) // 内容
	buffer.WriteByte(DELIMITER) // 分隔符

	return conn.Write(buffer.Bytes())
}

//========================================================

func strToInt32(str string) (int32, error) {
	//
	num, err := strconv.ParseInt(str, 10, 0)
	if err != nil {
		return 0, fmt.Errorf("[%s] is not integer.", str)
	}

	// 范围越界
	if num > math.MaxInt32 || num < math.MinInt32 {
		return 0, fmt.Errorf("[%d] is not 32-bit integer.", num)
	}

	return int32(num), nil // 正常转换
}

func cbrt(param int32) float64 {
	//
	return math.Cbrt(float64(param))
}

//========================================================

func printLog(role string, sn int, format string, args ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"

	}
	fmt.Printf("%s[%d]: %s", role, sn, fmt.Sprintf(format, args...))
}

func printServerLog(format string, args ...interface{}) {
	printLog("Server", 0, format, args...)
}

func printClientLog(sn int, format string, args ...interface{}) {
	printLog("Client", sn, format, args...)
}
