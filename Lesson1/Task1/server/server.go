package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

var (
	messageCh = make(chan string)
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	go serviceMessage()

	for {
		con, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go handleCon(con)
	}

}

// handleCon с использованием конструкции select-case смотрим два канала.
// Полученную инфу кидаем в net.Conn
func handleCon(c net.Conn) {
	defer c.Close()
	for {
		select {
		case <-time.After(time.Second):
			_, err := io.WriteString(c, time.Now().Format("15:04:05\n\r"))
			if err != nil {
				return
			}
		case data := <-messageCh:
			fmt.Println("get new service message", data)
			_, err := io.WriteString(c, data)
			if err != nil {
				return
			}
		}
	}
}

// serviceMessage обрабатывает os.Stdin, в случае, если на ввод подается строка,
// она отправляется в канал messageCh
func serviceMessage() {
	go func() {
		for {
			buf, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(buf)
			messageCh <- buf
		}
	}()
}
