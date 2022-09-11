package main

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	messages chan *ReceivedMessage

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// storing sub and pub messages.
	storedMessages map[string][]*ReceivedMessage // topic -> message
}

type ReceivedMessage struct {
	client  *Client
	message *WalletConnectMessage
}

func newHub() *Hub {
	return &Hub{
		messages:   make(chan *ReceivedMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),

		storedMessages: make(map[string][]*ReceivedMessage),
	}
}

func getPub(h *Hub, topic string) []*ReceivedMessage {
	pubKey := "socketMessage:" + topic
	return h.storedMessages[pubKey]
}

func getSub(h *Hub, topic string) []*ReceivedMessage {
	return h.storedMessages[topic]
}

func setPub(h *Hub, topic string, receivedMessage *ReceivedMessage) {
	pubKey := "socketMessage:" + topic
	arr := h.storedMessages[pubKey]
	arr = append(arr, receivedMessage)
	h.storedMessages[pubKey] = arr
}

func setSub(h *Hub, topic string, receivedMessage *ReceivedMessage) {
	subs := h.storedMessages[topic]
	subs = append(subs, receivedMessage)
	h.storedMessages[topic] = subs
}

func socketSendAll(socket chan []byte, messages []*ReceivedMessage) {
	if len(messages) > 0 {
		for _, m := range messages {
			socketSend(socket, m)
		}
	}
}

func socketSend(socket chan []byte, rm *ReceivedMessage) {
	socket <- []byte(SerializeMessage(rm.message))
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case receivedMessage := <-h.messages:
			topic := receivedMessage.message.Topic
			switch receivedMessage.message.MessageType {
			case "sub":
				setSub(h, topic, receivedMessage)
				pubs := getPub(h, topic)
				go socketSendAll(receivedMessage.client.send, pubs)
			case "pub":
				subs := getSub(h, topic)

				if !receivedMessage.message.Silent {
					// todo: push notification...
				}

				if len(subs) > 0 {
					for _, sub := range subs {
						go socketSend(sub.client.send, receivedMessage)
					}
				} else {
					setPub(h, topic, receivedMessage)
				}
			}
		}
	}
}
