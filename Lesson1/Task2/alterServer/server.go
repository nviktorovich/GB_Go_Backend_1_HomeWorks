package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type Message struct {
	Name string `json:"Name"`
	Msg  string `json:"Msg"`
}

var Clients = make(map[string]net.Conn)
var (
	clientMessageChannel = make(chan []byte)
	serverMessageChannel = make(chan []byte)
	errorChannel         = make(chan error)
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = listener.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	go Broadcast()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go HandleConn(conn)
	}

}

// HandleConn обрабатывает соединение. При новом подключении выводит адрес
// клиента, затем, в непрерывном цикле считывает в буффер, размером 1024 байта
// данные, которые с помошью json.Encoder распаковываются в объект структуры Message
func HandleConn(conn net.Conn) {
	fmt.Println(conn.RemoteAddr(), "was connected")
	Clients[conn.LocalAddr().String()] = conn
	for {
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		MessageProcessor(buf)
	}
}

func SendMessageToClient(conn net.Conn, data []byte) {
	if _, err := conn.Write(data); err != nil {
		log.Fatal(err)
	}
}

// Broadcast слушает каналы с сообщениями, в случае, если приходит сообщение от
// клиента, либо от сервера, дублирует его все подключенным клиентам.
func Broadcast() {
	for {
		select {
		case cliData := <-clientMessageChannel:
			for _, cli := range Clients {
				SendMessageToClient(cli, cliData)
			}
		case servData := <-serverMessageChannel:
			for _, cli := range Clients {
				SendMessageToClient(cli, servData)
			}
		case err := <-errorChannel:
			log.Fatal(err)
		}

	}
}

func MessageProcessor(buff []byte) {
	NewMessage := new(Message)

	clientMessageChannel <- buff

	reqBodyBytes := new(bytes.Buffer)

	if _, err := reqBodyBytes.Write(buff); err != nil {
		errorChannel <- err
	}

	if err := json.NewDecoder(reqBodyBytes).Decode(NewMessage); err != nil {
		if err != io.EOF {
			errorChannel <- err
		}
	}
	if _, err := fmt.Fprint(
		os.Stdout,
		fmt.Sprintf("%s say: %s\n", NewMessage.Name, NewMessage.Msg)); err != nil {
		log.Fatal(err)
	}
}
