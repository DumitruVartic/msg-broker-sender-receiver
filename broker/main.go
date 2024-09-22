package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

const (
	PORT = ":65432"
)

type Message struct {
	Topic   string `json:"topic" xml:"topic"`
	Content string `json:"content" xml:"content"`
}

var subscribers = make(map[string][]net.Conn) // Key: topic name, Value: slice of connections

func main() {
	ln, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer ln.Close() // Ensure the listener is closed when exiting main

	fmt.Println("Message Broker is listining on port", PORT)

	// Main loop for incoming conns
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close() // Ensure conn close on func exit
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error readin from connection:", err)
		return
	}

	message := string(buffer[:n])            // Converting bytes read to a string
	parts := strings.SplitN(message, ":", 3) // Splits message in parts by ":"

	if len(parts) < 2 { // if the content part is missing
		return // Not enough parts
	}

	// Handling command from sender
	cmd, topic, content := parts[0], parts[1], ""
	if len(parts) == 3 {
		content = parts[2] // if the contntent/message is present
	}

	switch cmd {
	case "publish":
		publishMessage(topic, content)
	case "subscribe":
		subscribe(topic, conn)
	}

}

func publishMessage(topic, content string) {
	message := Message{Topic: topic, Content: content}
	jsonData, _ := json.Marshal(message)

	// if there are subscribers for topic, send them mesage
	if conns, found := subscribers[topic]; found {
		for _, subConn := range conns { // for each sub conn
			subConn.Write(jsonData)
		}
	}
	fmt.Printf("Published to topic \"%s\": %s\n", topic, content)
}

func subscribe(topic string, conn net.Conn) {
	subscribers[topic] = append(subscribers[topic], conn)
	fmt.Printf("Subscribed to topic \"%s\"\n", topic)
}
