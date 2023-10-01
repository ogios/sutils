import socket
import time

HOST = "127.0.0.1"
PORT = 15002

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.bind((HOST, PORT))
s.listen(1)
c, b = s.accept()

a = b""
start = time.time()
try:
    while 1:
        if time.time() - start > 3:
            break
        d = c.recv(1024)
        print(d)
        # if not d:
        #     break
        a += d
        time.sleep(0.1)
except:
    pass

print(a.decode())
