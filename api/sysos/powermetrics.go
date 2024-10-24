package sysos

import (
	"bufio"
	"fmt"
	"os/exec"
)

var command = []string{"powermetrics", "-i", "10000", "-s", "cpu_power"}

// 获取系统信息
func GetPowermetrics(initial string) {
	if len(command) == 0 {
		return
	}

	if initial != "" {
		command[2] = initial
	}

	// 使用 sudo 执行 powermetrics 命令
	cmd := exec.Command("sudo", command...)

	// 获取命令的标准输出
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error getting stdout pipe:", err)
		return
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting command:", err)
		return
	}

	// 创建一个 bufio.Scanner 来逐行读取输出
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		handlerPower(line)
	}

	// 等待命令执行完成
	if err := cmd.Wait(); err != nil {
		fmt.Println("Error waiting for command:", err)
	}
}
