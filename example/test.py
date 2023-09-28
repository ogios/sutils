import socket

HOST = "127.0.0.1"
PORT = 15002

c = "腹蛤属的撒".encode("utf-8")
a = [len(c), 0]
b = bytes(a) + c
s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

s.connect((HOST, PORT))
s.send(b)
s.close()
