package main

import (
	"encoding/json"
	"fmt"
)

type WalletConnectMessage struct {
	Silent      bool   `json:"silent"`
	Topic       string `json:"topic"`
	Payload     string `json:"payload"`
	MessageType string `json:"type"`

	Sender *Client `json:"-"` // ignored
}

func ParseMessage(client *Client, content string) *WalletConnectMessage {
	message := new(WalletConnectMessage)
	message.Sender = client
	err := json.Unmarshal([]byte(content), &message)
	if err != nil {
		fmt.Printf("Error on json decode: %s", err)
		return message
	}

	return message
}

func SerializeMessage(message *WalletConnectMessage) string {
	bytes, err := json.Marshal(message)
	if err != nil {
		fmt.Printf("Error on json encode: %s", err)
		return ""
	}
	return string(bytes)
}
