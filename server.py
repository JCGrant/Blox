import socket

SERVER_ADDRESS = (HOST, PORT) = ('', 8888)
REQUEST_QUEUE_SIZE = 5

def handlere_request(client_connection):
    request = client_connection.recv(1024)
    print(request)

    http_response = b"""\
HTTP/1.1 200 OK

Hello, World!
"""
    client_connection.sendall(http_response)

def serve_forever():
    listen_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    listen_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    listen_socket.bind(SERVER_ADDRESS)
    listen_socket.listen(REQUEST_QUEUE_SIZE)
    print('Serving HTTP on port %s ...' % PORT)
    while True:
        client_connection, client_address = listen_socket.accept()
        handlere_request(client_connection)
        client_connection.close()

if __name__ == '__main__':
    serve_forever()
