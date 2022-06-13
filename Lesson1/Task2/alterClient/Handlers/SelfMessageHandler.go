package Handlers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
)

func GetSelfMessage(name string) {
	message := new(Message)
	message.Name = name
	for {
		data, _, err := bufio.NewReader(os.Stdin).ReadLine()
		if err != nil {
			ErrorChan <- err
		}

		message.Msg = string(data)
		reqBodyBytes := new(bytes.Buffer)

		if err = json.NewEncoder(reqBodyBytes).Encode(message); err != nil {
			ErrorChan <- err
		}
		BytesFromSelfChan <- reqBodyBytes.Bytes()
	}
}
