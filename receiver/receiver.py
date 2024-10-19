import grpc
import argparse
import json
import xml.etree.ElementTree as ET
import message_broker_pb2
import message_broker_pb2_grpc

HOST = 'localhost'
PORT = 65432

def create_message(command, message_type, topic, content=None):
    """Creates a message in the specified format."""
    if message_type == 'json':
        message = {
            "topic": topic,
            "content": content if content else None
        }
        return json.dumps(message, indent=4) if content else json.dumps({"command": command, "topic": topic})

    elif message_type == 'xml':
        message_element = ET.Element("Message" if content else "message")
        if content:
            ET.SubElement(message_element, "Topic").text = topic
            ET.SubElement(message_element, "Content").text = content
        else:
            ET.SubElement(message_element, "Command").text = command
            ET.SubElement(message_element, "Topic").text = topic
        return ET.tostring(message_element, encoding="unicode")

def send_command_to_broker(command, topic, output_format):
    """Handles subscription to a topic and receives messages."""
    # Establish a gRPC channel
    with grpc.insecure_channel(f'{HOST}:{PORT}') as channel:
        stub = message_broker_pb2_grpc.MessageBrokerStub(channel)

        if command == "subscribe":
            request = message_broker_pb2.TopicRequest(topic=topic)
            response_iterator = stub.Subscribe(request)

            print(f"Subscribed to topic: {topic}")
            for message in response_iterator:
                formatted_message = create_message(command, output_format, topic, message.content)
                print(formatted_message)

        elif command == "unsubscribe":
            request = message_broker_pb2.TopicRequest(topic=topic)
            response = stub.Unsubscribe(request)
            print(response.message)

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Subscribe to a topic.')
    parser.add_argument('topic', help='The topic to subscribe to')
    parser.add_argument('--format', choices=['json', 'xml'], default='json', help='Output format for messages (default: json)')
    parser.add_argument('--unsubscribe', action='store_true', help='Unsubscribe from the topic instead of subscribing')
    
    args = parser.parse_args()

    command = "unsubscribe" if args.unsubscribe else "subscribe"
    
    send_command_to_broker(command, args.topic, args.format)