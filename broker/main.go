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
	"time"
)

const PORT = ":65432"

type Message struct {
	Content string `json:"content" xml:"Content"`
}

type MessageMetadata struct {
	Message Message `json:"message" xml:"Message"`
	Command string  `json:"command" xml:"Command"`
	Topic   string  `json:"topic" xml:"Topic"`
	Format  string
}

type Subscriber struct {
	Conn   net.Conn
	Format string
}

var subscribers = make(map[string][]Subscriber)
var messages = make(map[string][]Message) // Storing messages by topic
var mu sync.Mutex
var shutdown = false

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
		shutdown = true
		time.Sleep(1 * time.Second)
		ln.Close()
		os.Exit(0)
	}()

	fmt.Println("Message Broker is listening on port", PORT)

	for {
		conn, err := ln.Accept()
		if err != nil {
			if shutdown {
				fmt.Println("Server shutting down, no longer accepting connections.")
				return
			}
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

	var metadata MessageMetadata
	format, err := determineMessageFormat(buffer[:n], &metadata)
	if err != nil {
		fmt.Println("Failed to parse message:", err)
		return
	}

	metadata.Format = format

	fmt.Printf("Parsed message as %s\n", metadata.Format)
	fmt.Printf("Message %s\n", metadata.Message.Content)
	fmt.Printf("Meta %s\n", metadata)
	handleCommand(metadata, conn)
}

func determineMessageFormat(data []byte, metadata *MessageMetadata) (string, error) {
	if json.Valid(data) {
		if err := json.Unmarshal(data, metadata); err != nil {
			return "", fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
		return "json", nil
	}

	if err := xml.Unmarshal(data, metadata); err == nil {
		return "xml", nil
	}

	return "", fmt.Errorf("message could not be parsed as JSON or XML")
}

func handleCommand(metadata MessageMetadata, conn net.Conn) {
	switch metadata.Command {
	case "subscribe":
		handleSubscribe(metadata.Topic, conn, metadata.Format)
	case "publish":
		handlePublish(metadata.Topic, metadata.Message.Content) // Use the content of the Message
	case "unsubscribe":
		handleUnsubscribe(metadata.Topic, conn)
	}
}

func handleUnsubscribe(topic string, conn net.Conn) {
	mu.Lock()
	defer mu.Unlock()

	subscribers[topic] = removeSubscriber(subscribers[topic], conn)
	fmt.Printf("Client unsubscribed from topic \"%s\"\n", topic)
}

func removeSubscriber(subs []Subscriber, conn net.Conn) []Subscriber {
	var updatedSubs []Subscriber
	for _, sub := range subs {
		if sub.Conn != conn {
			updatedSubs = append(updatedSubs, sub)
		}
	}
	return updatedSubs
}

func marshalMessage(content string, format string) ([]byte, error) {
	message := Message{Content: content}
	if format == "json" {
		return json.Marshal(message)
	} else if format == "xml" {
		return xml.Marshal(message)
	}
	return nil, fmt.Errorf("unsupported format: %s", format)
}

func handleSubscribe(topic string, conn net.Conn, format string) {
	mu.Lock()
	defer mu.Unlock()

	subscribers[topic] = append(subscribers[topic], Subscriber{Conn: conn, Format: format})
	fmt.Printf("Client subscribed to topic \"%s\" with format \"%s\"\n", topic, format)

	for _, msg := range messages[topic] {
		if data, err := marshalMessage(msg.Content, format); err == nil {
			conn.Write(data)
		}
	}
}

func handlePublish(topic, content string) {
	mu.Lock()
	defer mu.Unlock()

	messages[topic] = append(messages[topic], Message{Content: content})

	for _, sub := range subscribers[topic] {
		if data, err := marshalMessage(content, sub.Format); err == nil {
			sub.Conn.Write(data)
		}
	}
	fmt.Printf("Published to topic \"%s\": %s\n", topic, content)
}
