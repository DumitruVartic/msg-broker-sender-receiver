import socket
import argparse

HOST = 'localhost'
PORT = 65432

def subscribe_to_topic(topic):
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as client_socket:
        # Connect to the broker
        client_socket.connect((HOST, PORT))
        
        # Send subscription request
        client_socket.sendall(f'subscribe:{topic}'.encode())
        
        while True:
            # Receive messages
            data = client_socket.recv(1024)
            if not data:
                break
            print(f'Received from topic "{topic}": {data.decode()}')

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Subscribe to a topic.')
    parser.add_argument('topic', help='The topic to subscribe to')
    args = parser.parse_args()
    subscribe_to_topic(args.topic)

