import subprocess
import sys
from time import strftime

start_command = 'java -jar minecraft_server.1.9.2.jar nogui'

COMMAND = 'COMMAND'
SHUTDOWN = 'SHUTDOWN'


class Wrapper:

    def __init__(self):
        self.proc = subprocess.Popen(start_command.split(),
                                     stdout=subprocess.PIPE,
                                     stdin=subprocess.PIPE,
                                     stderr=subprocess.PIPE)
        print("Started Minecraft server with pid: {}".format(self.proc.pid))
        self.app_name = 'MCSPyW'
        self.command_token = '!'

    def read_from_server(self):
        data = self.proc.stdout.readline()
        return data.decode('utf-8').rstrip('\n')

    def send_to_server(self, string):
        self.proc.stdin.write(bytes(string + '\n', 'utf-8'))
        self.proc.stdin.flush()

    def handle_command(self, command_text):
        self.log(command_text, COMMAND)
        self.send_to_server(command_text)

    def handle_text(self, server_text):
        print(server_text)
        player_text = ' '.join(server_text.split()[4:])
        if player_text.startswith(self.command_token):
            command_text = player_text.lstrip(self.command_token)
            self.handle_command(command_text)

    def log(self, message, message_type):
        time = strftime('%H:%M:%S')
        print('[{}] [{} - {}]:  {}'.format(time, self.app_name, message_type, message))

    def shutdown(self):
        self.log('Shutting down PMCSW.', SHUTDOWN)
        self.proc.kill()
        sys.exit(0)

    def run(self):
        while True:
            try:
                server_text = self.read_from_server()
                self.handle_text(server_text)
            except (KeyboardInterrupt):
                self.shutdown()

if __name__ == '__main__':
    w = Wrapper()
    w.run()
