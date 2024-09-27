package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const PORT = ":65432"

type Message struct {
	Command string `json:"command" xml:"command"`
	Topic   string `json:"topic" xml:"topic"`
	Content string `json:"content" xml:"content"`
}

type Subscriber struct {
	Conn   net.Conn
	Format string
}

var subscribers = make(map[string][]Subscriber)
var mu sync.Mutex

func main() {
	ln, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer ln.Close()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		fmt.Println("Shutting down gracefully...")
		ln.Close()
		os.Exit(0)
	}()

	fmt.Println("Message Broker is listening on port", PORT)

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
	defer conn.Close()

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}

	var message Message
	// Try to unmarshal the message as JSON
	if err := json.Unmarshal(buffer[:n], &message); err == nil {
		fmt.Println("Parsed message as JSON")
		handleCommand(message, conn, "json")
		return
	}

	// Try to unmarshal the message as XML
	if err := xml.Unmarshal(buffer[:n], &message); err == nil {
		fmt.Println("Parsed message as XML")
		handleCommand(message, conn, "xml")
		return
	}

	fmt.Println("Failed to parse message")
}

func handleCommand(message Message, conn net.Conn, format string) {
	switch message.Command {
	case "subscribe":
		handleSubscribe(message.Topic, conn, format)
	case "publish":
		handlePublish(message.Topic, message.Content)
	case "unsubscribe":
		handleUnsubscribe(message.Topic, conn)
	}
}
func handleUnsubscribe(topic string, conn net.Conn) {
	mu.Lock()
	defer mu.Unlock()

	subscribers[topic] = removeSubscriber(subscribers[topic], conn)
	fmt.Printf("Client unsubscribed from topic \"%s\"\n", topic)

}

func removeSubscriber(subs []Subscriber, conn net.Conn) []Subscriber {
	for i, sub := range subs {
		if sub.Conn == conn {
			return append(subs[:i], subs[i+1:]...)
		}
	}
	return subs

}

func handleSubscribe(topic string, conn net.Conn, format string) {
	mu.Lock()
	defer mu.Unlock()

	subscribers[topic] = append(subscribers[topic], Subscriber{Conn: conn, Format: format})
	fmt.Printf("Client subscribed to topic \"%s\" with format \"%s\"\n", topic, format)
}

func handlePublish(topic, content string) {
	mu.Lock()
	defer mu.Unlock()

	message := Message{Command: "publish", Topic: topic, Content: content}

	for _, sub := range subscribers[topic] {
		var err error
		var data []byte

		// Convert the message to the subscriber's preferred format (detected on trying to parse)
		if sub.Format == "json" {
			data, err = json.Marshal(message)
		} else if sub.Format == "xml" {
			data, err = xml.Marshal(message)
		}

		if err == nil {
			sub.Conn.Write(data)
		}
	}
	fmt.Printf("Published to topic \"%s\": %s\n", topic, content)
}
