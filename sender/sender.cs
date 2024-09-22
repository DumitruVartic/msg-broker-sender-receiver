using System;
using System.Net.Sockets;
using System.Text;

class Program
{
    const string HOST = "localhost";
    const int PORT = 65432;

    static void Main(string[] args)
    {
        Console.Write("Enter topic to publish to: ");
        string topic = Console.ReadLine();
        
        Console.Write("Enter message to publish: ");
        string message = Console.ReadLine();

        PublishToTopic(topic, message);
    }

    static void PublishToTopic(string topic, string message)
    {
        using (TcpClient client = new TcpClient(HOST, PORT))
        using (NetworkStream stream = client.GetStream())
        {
            // Create the publish message
            string publishMessage = $"publish:{topic}:{message}";
            byte[] data = Encoding.UTF8.GetBytes(publishMessage);

            // Send the publish message to the broker
            stream.Write(data, 0, data.Length);
            Console.WriteLine($"Sent to topic \"{topic}\": {message}");
        }
    }
}