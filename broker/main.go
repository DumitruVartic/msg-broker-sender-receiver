package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "message_broker/proto"

	"google.golang.org/grpc"
)

const PORT = ":65432"

type server struct {
	pb.UnimplementedMessageBrokerServer
}

var (
	subscribers = make(map[string][]chan *pb.Message) // Map of topic to a list of channels
	mu          sync.Mutex
)

func (s *server) Publish(ctx context.Context, metadata *pb.MessageMetadata) (*pb.Response, error) {
	mu.Lock()
	defer mu.Unlock()

	topic := metadata.Topic
	message := metadata.Message.Content
	log.Print("Message", message, "Published to topic: ", topic)

	// Broadcast to all subscribers of the topic
	for _, sub := range subscribers[topic] {
		sub <- &pb.Message{Content: message} // Send the message to the channel
	}

	return &pb.Response{Success: true, Message: "Message published"}, nil
}

func (s *server) Subscribe(req *pb.TopicRequest, stream pb.MessageBroker_SubscribeServer) error {
	topic := req.Topic
	mu.Lock()
	ch := make(chan *pb.Message)                        // Create a new channel for the subscriber
	subscribers[topic] = append(subscribers[topic], ch) // Add the subscriber
	mu.Unlock()
	log.Print("Subscribed to topic", topic)
	// Ensure we clean up the subscriber on function exit
	defer func() {
		mu.Lock()
		defer mu.Unlock()
		for i, sub := range subscribers[topic] {
			if sub == ch {
				subscribers[topic] = append(subscribers[topic][:i], subscribers[topic][i+1:]...) // Remove subscriber
				break
			}
		}
	}()

	// Wait for messages and send them to the stream
	for {
		msg, ok := <-ch // Wait for a message
		if !ok {
			return nil // Channel closed
		}
		if err := stream.Send(msg); err != nil {
			return err // Stream error
		}
	}
}

func (s *server) Unsubscribe(ctx context.Context, req *pb.TopicRequest) (*pb.Response, error) {
	mu.Lock()
	defer mu.Unlock()

	topic := req.Topic
	if _, exists := subscribers[topic]; exists {
		delete(subscribers, topic)
		return &pb.Response{Success: true, Message: "Unsubscribed from topic"}, nil
	}

	return &pb.Response{Success: false, Message: "Topic not found"}, nil
}

func main() {
	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMessageBrokerServer(grpcServer, &server{})

	log.Printf("Server is listening on port %s...", PORT)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
