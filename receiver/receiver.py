import socket
import argparse
import json
import xml.etree.ElementTree as ET

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
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as client_socket:
        client_socket.connect((HOST, PORT))

        subscribe_message = create_message(command, output_format, topic)
        client_socket.sendall(subscribe_message.encode())

        while True:
            data = client_socket.recv(1024)
            if not data:
                break
            content = data.decode()
            formatted_message = create_message(command, output_format, topic, content)
            print(formatted_message)

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Subscribe to a topic.')
    parser.add_argument('topic', help='The topic to subscribe to')
    parser.add_argument('--format', choices=['json', 'xml'], default='json', help='Output format for messages (default: json)')
    parser.add_argument('--unsubscribe', action='store_true', help='Unsubscribe from the topic instead of subscribing')
    
    args = parser.parse_args()

    command = "unsubscribe" if args.unsubscribe else "subscribe"
    
    send_command_to_broker(command, args.topic, args.format)