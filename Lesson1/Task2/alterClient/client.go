package main

import (
	"client/Handlers"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {

	name, err := getName()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = conn.Close()
	}()

	go Handlers.ChanListener(conn)
	go Handlers.GetSelfMessage(name)
	go Handlers.GetIncomeMessage(conn)

	for {
		select {
		case err = <-Handlers.ErrorChan:
			log.Fatal(err)
		}
	}

}

// getName возвращает имя, пользователя (ввод с клавиатуры) или ошибку
func getName() (name string, err error) {
	_, err = fmt.Fprint(os.Stdout, "Введите имя пользователя\n")
	if err != nil {
		return
	}

	_, err = fmt.Fscanf(os.Stdin, "%s", &name)
	if err != nil {
		return
	}

	return
}
