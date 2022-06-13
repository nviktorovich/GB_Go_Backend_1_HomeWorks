package Handlers

import (
	"net"
)

func GetIncomeMessage(conn net.Conn) {
	buff := make([]byte, 1024)
	for {
		_, err := conn.Read(buff)
		if err != nil {
			ErrorChan <- err
		}

		BytesFromServerChan <- buff

	}
}
