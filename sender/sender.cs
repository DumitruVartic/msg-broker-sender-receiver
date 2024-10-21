using System;
using System.Threading.Tasks;
using Grpc.Net.Client;
using Broker; // Generated from the .proto file
using Google.Protobuf.WellKnownTypes;

class Program
{
    const string HOST = "localhost:50051"; // gRPC server address

    static async Task Main(string[] args)
    {
        // Validate that command, topic, and optionally message are provided
        if (args.Length < 2 || string.IsNullOrWhiteSpace(args[0]) || string.IsNullOrWhiteSpace(args[1]))
        {
            Console.WriteLine("Usage: Program <command> <topic> [message] [--format=json|xml]");
            Console.WriteLine("Both command and topic must be provided and not be empty.");
            return;
        }

        // Extract the command, topic, and optional message
        string command = args[0].ToLower();
        string topic = args[1];
        string messageContent = args.Length > 2 ? args[2] : null;

        // Establish gRPC connection
        using var channel = GrpcChannel.ForAddress($"http://{HOST}");
        var client = new MessageBroker.MessageBrokerClient(channel);

        switch (command)
        {
            case "publish":
                if (string.IsNullOrWhiteSpace(messageContent))
                {
                    Console.WriteLine("Message content is required for publishing.");
                    return;
                }
                await PublishMessage(client, topic, messageContent);
                break;

            case "subscribe":
                await SubscribeToTopic(client, topic);
                break;

            case "unsubscribe":
                await UnsubscribeFromTopic(client, topic);
                break;

            default:
                Console.WriteLine($"Unknown command: {command}. Use publish, subscribe, or unsubscribe.");
                break;
        }
    }

    static async Task PublishMessage(MessageBroker.MessageBrokerClient client, string topic, string messageContent)
    {
        var message = new Message { Content = messageContent };
        var metadata = new MessageMetadata
        {
            Message = message,
            Command = "publish",
            Topic = topic
        };

        try
        {
            var response = await client.PublishAsync(metadata);
            Console.WriteLine($"Response: {response.Message}, Success: {response.Success}");
        }
        catch (Exception ex)
        {
            Console.WriteLine($"Error during publish: {ex.Message}");
        }
    }

    static async Task SubscribeToTopic(MessageBroker.MessageBrokerClient client, string topic)
    {
        var request = new TopicRequest { Topic = topic };
        try
        {
            using var call = client.Subscribe(request);

            Console.WriteLine($"Subscribed to topic: {topic}");
            await foreach (var message in call.ResponseStream.ReadAllAsync())
            {
                Console.WriteLine($"Received message: {message.Content}");
            }
        }
        catch (Exception ex)
        {
            Console.WriteLine($"Error during subscription: {ex.Message}");
        }
    }

    static async Task UnsubscribeFromTopic(MessageBroker.MessageBrokerClient client, string topic)
    {
        var request = new TopicRequest { Topic = topic };

        try
        {
            var response = await client.UnsubscribeAsync(request);
            Console.WriteLine($"Unsubscribe response: {response.Message}, Success: {response.Success}");
        }
        catch (Exception ex)
        {
            Console.WriteLine($"Error during unsubscribe: {ex.Message}");
        }
    }
}
