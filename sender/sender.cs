using System;
using System.Net.Sockets;
using System.Text;
using System.Text.Json;
using System.Xml.Serialization;
using System.IO;

class Program
{
    const string HOST = "localhost";
    const int PORT = 65432;

    static void Main(string[] args)
    {
        // Validate that both topic and message are provided and not empty
        if (args.Length < 2 || string.IsNullOrWhiteSpace(args[0]) || string.IsNullOrWhiteSpace(args[1]))
        {
            Console.WriteLine("Usage: Program <topic> <message> [--format=json|xml]");
            Console.WriteLine("Both topic and message must be provided and not be empty.");
            return;
        }

        // Get the topic and message from command-line arguments
        string topic = args[0];
        string message = args[1];

        // Default to JSON format if format is not specified
        string format = "json";
        if (args.Length >= 3 && args[2].StartsWith("--format="))
        {
            format = args[2].Split('=')[1].ToLower();
        }

        // Validate the format
        if (format != "json" && format != "xml")
        {
            Console.WriteLine("Invalid format. Supported formats are json and xml.");
            return;
        }

        // Publish to the topic with the selected format
        PublishToTopic(topic, message, format);
    }

    static void PublishToTopic(string topic, string message, string format)
    {
        // Create the publish message object
        var publishMessage = new PublishMessage
        {
            Command = "publish",
            Topic = topic,
            Message = new MessageContent { Content = message }
        };

        // Serialize the publish message based on the format
        string serializedMessage = format == "json" ? SerializeToJson(publishMessage) : SerializeToXml(publishMessage);

        // Send the message
        try
        {
            using (TcpClient client = new TcpClient(HOST, PORT))
            using (NetworkStream stream = client.GetStream())
            {
                byte[] data = Encoding.UTF8.GetBytes(serializedMessage);

                // Send the publish message to the broker
                stream.Write(data, 0, data.Length);
                Console.WriteLine($"Sent to topic \"{topic}\" in {format.ToUpper()} format: {message}");
            }
        }
        catch (Exception ex)
        {
            Console.WriteLine($"Error: {ex.Message}");
        }
    }

    static string SerializeToJson(PublishMessage message)
    {
        return JsonSerializer.Serialize(message, new JsonSerializerOptions { WriteIndented = true });
    }

    static string SerializeToXml(PublishMessage message)
    {
        var xmlSerializer = new XmlSerializer(typeof(PublishMessage));
        var xmlSettings = new XmlWriterSettings
        {
            Indent = true,
            Encoding = new UTF8Encoding(false), // Force UTF-8 without BOM
            OmitXmlDeclaration = true           // Remove XML declaration
        };

        var ns = new XmlSerializerNamespaces();
        ns.Add("", ""); // Remove namespaces

        using (var stringWriter = new StringWriter())
        using (var xmlWriter = XmlWriter.Create(stringWriter, xmlSettings))
        {
            xmlSerializer.Serialize(xmlWriter, message, ns);
            return stringWriter.ToString();
        }
    }
}

// Define a class for the publish message
public class PublishMessage
{
    public string Command { get; set; }
    public string Topic { get; set; }
    public MessageContent Message { get; set; }
}

public class MessageContent
{
    public string Content { get; set; }
}
