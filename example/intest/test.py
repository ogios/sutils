import socket

HOST = "127.0.0.1"
PORT = 15002


def int_to_255(total) -> list[int]:
    l: list[int] = []
    while total >= 255:
        l.append(total % 255)
        total //= 255
    l.append(total)
    return l


def get_len(content: bytes) -> list[int]:
    l = int_to_255(len(content))
    l.append(255)
    return l


def add_bytes(old: bytes, add: bytes) -> bytes:
    tl = get_len(add)
    return old + bytes(tl) + add


def add_string(old: bytes, content: str) -> bytes:
    t = content.encode()
    return add_bytes(old, t)

b=b""
b+=bytes(int_to_255(0) + [255])
print(b)

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

s.connect((HOST, PORT))
s.send(b)
s.close()
