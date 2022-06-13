package Handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
)

var (
	ErrorChan           = make(chan error)
	BytesFromServerChan = make(chan []byte)
	BytesFromSelfChan   = make(chan []byte)
	MessageChan         = make(chan Message)
)

func ChanListener(conn net.Conn) {
	for {
		select {

		case selfMessage := <-BytesFromSelfChan:

			_, err := conn.Write(selfMessage)
			if err != nil {
				ErrorChan <- err
			}

		case dataForDecoding := <-BytesFromServerChan:
			go ReadMessage(dataForDecoding)

		case Msg := <-MessageChan:
			_, err := fmt.Fprintf(os.Stdout, "from %v: %v\n", Msg.Name, Msg.Msg)
			if err != nil {
				ErrorChan <- err
			}

		}
	}
}

func ReadMessage(data []byte) {

	NewMessage := new(Message)

	reqBodyBytes := new(bytes.Buffer)

	if _, err := reqBodyBytes.Write(data); err != nil {
		ErrorChan <- err
	}

	if err := json.NewDecoder(reqBodyBytes).Decode(NewMessage); err != nil {
		if err != io.EOF {
			ErrorChan <- err
		}
	}

	MessageChan <- *NewMessage
}
