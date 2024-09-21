import socket

HOST = 'localhost'
PORT = 65432


def publish_to_topic(topic, message):
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as client_socket:
        # Connect to the broker
        client_socket.connect((HOST, PORT))

        # Send a publish request with the topic and message
        publish_message = f'publish:{topic}:{message}'
        client_socket.sendall(publish_message.encode())
        print(f'Sent to topic "{topic}": {message}')


if __name__ == "__main__":
    topic = input("Enter topic to publish to: ")
    message = input("Enter message to publish: ")
    publish_to_topic(topic, message)
