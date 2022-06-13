package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	con, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	defer con.Close()

	buf := make([]byte, 256) // Буфер для того, чтобы в него считывать информацию

	for {
		_, err := con.Read(buf)
		if err == io.EOF {
			break
		}
		io.WriteString(os.Stdout, fmt.Sprintf("get data from server: %s", string(buf)))
	}

}
