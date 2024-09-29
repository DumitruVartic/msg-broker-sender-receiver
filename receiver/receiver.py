import socket
import argparse
import json
import xml.etree.ElementTree as ET

HOST = 'localhost'
PORT = 65432

def create_message(message_type, topic, content=None:
    if message_type == 'json':
        message = {
            "topic": topic,
            "content": content if content else None
        }
        return json.dumps(message, indent=4) if content else json.dumps({"command": "subscribe", "topic": topic})

    elif message_type == 'Xml':
        message_element = ET.Element("Message" if content else "message")
        if content:
            ET.SubElement(message_element, "Topic").text = topic
            ET.SubElement(message_element, "Content").text = content
        else:
            ET.SubElement(message_element, "command").text = "subscribe"
            ET.SubElement(message_element, "topic").text = topic
        return ET.tostring(message_element, encoding="unicode")
)

def subscribe_to_topic(topic, output_format):
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as client_socket:
        client_socket.connect((HOST, PORT))
        
        subscribe_message = create_subscribe_message(topic, output_format)
        client_socket.sendall(subscribe_message.encode())
        
        while True:
            data = client_socket.recv(1024)
            if not data:
                break
            content = data.decode()
            
            formatted_message = format_message(topic, content, output_format)
            print(formatted_message)

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Subscribe to a topic.')
    parser.add_argument('topic', help='The topic to subscribe to')
    parser.add_argument('--format', choices=['json', 'xml'], default='json', help='The output format for sending subscription and receiving messages (default: json)')
    
    args = parser.parse_args()
    
    subscribe_to_topic(args.topic, args.format)
