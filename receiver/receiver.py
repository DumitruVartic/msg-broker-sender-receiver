import socket
import argparse
import json
import xml.etree.ElementTree as ET

HOST = 'localhost'
PORT = 65432

def format_message(topic, content, output_format):
    if output_format == 'json':
        message = {
            "topic": topic,
            "content": content
        }
        return json.dumps(message, indent=4)
    
    elif output_format == 'xml':
        message = ET.Element("Message")
        topic_element = ET.SubElement(message, "Topic")
        topic_element.text = topic
        content_element = ET.SubElement(message, "Content")
        content_element.text = content
        return ET.tostring(message, encoding="unicode", method="xml")

def create_subscribe_message(topic, output_format):
    if output_format == 'json':
        subscribe_message = {
            "command": "subscribe",
            "topic": topic
        }
        return json.dumps(subscribe_message)
    
    elif output_format == 'xml':
        message = ET.Element("message")
        command_element = ET.SubElement(message, "command")
        command_element.text = "subscribe"
        topic_element = ET.SubElement(message, "topic")
        topic_element.text = topic
        return ET.tostring(message, encoding="unicode", method="xml")

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
