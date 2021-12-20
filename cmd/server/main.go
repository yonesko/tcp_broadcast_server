package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"go.avito.ru/cmpl/tcp_broadcast_server/internal/tcpserver"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}

	l, err := net.Listen("tcp4", ":"+arguments[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	srv := tcpserver.New()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go srv.ReadRequests(readMsgs(c))
	}
}

func readMsgs(c net.Conn) chan tcpserver.Req {
	reqs := make(chan tcpserver.Req)
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	go func() {
		for {
			netData, err := bufio.NewReader(c).ReadString('\n')
			if err != nil {
				fmt.Println(err)
				reqs <- tcpserver.Req{
					C:      tcpserver.Conn{Conn: c},
					Text:   strings.TrimSpace(netData),
					Closed: true,
				}
				return
			}
			reqs <- tcpserver.Req{
				C:    tcpserver.Conn{Conn: c},
				Text: strings.TrimSpace(netData),
			}
		}
	}()
	return reqs
}
