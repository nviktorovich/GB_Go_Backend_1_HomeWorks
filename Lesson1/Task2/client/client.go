package main

import (
	"fmt"
	"io"
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
	defer conn.Close()

	go func() {
		//data, _, err := bufio.NewReader(conn).ReadLine()
		//if err != nil {
		//	log.Fatal(err)
		//}
		//fmt.Fprint(os.Stdout, string(data))
		io.Copy(os.Stdout, conn)
	}()

	for {
		var data = make([]byte, 1024)
		_, err = fmt.Fscan(os.Stdin, &data)
		if err != nil {
			log.Fatal(err)
		}
		message := name + string(data) + "\n"
		io.WriteString(conn, message)

		//io.Copy(conn, os.Stdin) // until you send ^Z
		fmt.Printf("%s: exit", conn.LocalAddr())
	}

}

func getName() (string, error) {
	var name string
	_, err := fmt.Fscanf(os.Stdin, "%s", &name)
	if err != nil {
		return "", err
	}
	return " " + name + ": ", nil
}
