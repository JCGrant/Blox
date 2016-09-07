import errno
import os
import signal
import socket
import subprocess

MC_SERVER_DIR = 'mc-server'
MC_SERVER_PROPERTIES_PATH = os.path.join(MC_SERVER_DIR, 'server.properties')

WRAPPER_ADDRESS = (_, PORT) = '', 25565
REQUEST_QUEUE_SIZE = 1024


def grim_reaper(signum, frame):
    while True:
        try:
            pid, status = os.waitpid(
                -1,          # Wait for any child process
                 os.WNOHANG  # Do not block and return EWOULDBLOCK error
            )
        except OSError:
            return

        if pid == 0:  # no more zombies
            return


def handle_request(client_connection, mc_server_address):
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as server_connection:
        server_connection.connect(mc_server_address)

        request = client_connection.recv(1024)
        print(b'req:' + request)
        server_connection.send(request)

        response = server_connection.recv(1024)
        print(b'res:' + response)
        client_connection.send(response)
        print()


def serve_forever(mc_server_address):
    listen_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    listen_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    listen_socket.bind(WRAPPER_ADDRESS)
    listen_socket.listen(REQUEST_QUEUE_SIZE)
    print('Wrapper running on port {port}...'.format(port=PORT))

    signal.signal(signal.SIGCHLD, grim_reaper)

    while True:
        try:
            client_connection, client_address = listen_socket.accept()
        except IOError as e:
            code, msg = e.args
            # restart 'accept' if it was interrupted
            if code == errno.EINTR:
                continue
            else:
                raise

        pid = os.fork()
        if pid == 0:  # child
            listen_socket.close()  # close child copy
            handle_request(client_connection, mc_server_address)
            client_connection.close()
            os._exit(0)
        else:  # parent
            client_connection.close()  # close parent copy and loop over


def run_mc_server():
    run_command = 'java -Xmx1024M -Xms1024M -jar minecraft_server.1.10.2.jar nogui'
    mc = subprocess.Popen(run_command.split(), cwd=MC_SERVER_DIR)
    return mc


def read_properties_file(filepath):
    with open(filepath) as f:
        contents = f.read()
    key_values = [row.split('=') for row in contents.split('\n')[2:-1]]
    return {key: value for (key, value) in key_values}


if __name__ == '__main__':
    try:
        print('Starting Minecraft Server...')
        mc = run_mc_server()
        while True:
            try:
                mc_port = int(read_properties_file(MC_SERVER_PROPERTIES_PATH)['server-port'])
                break
            except (FileNotFoundError, KeyError):
                continue
        mc_server_address = '', mc_port
        serve_forever(mc_server_address)
    except:
        mc.send_signal(signal.SIGINT)
