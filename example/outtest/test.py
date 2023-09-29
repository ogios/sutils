# a = 255**2
a = 255**2
b = []
count = 1
print(a)
while a >= 255:
    count += 1
    b.append(a % 255)
    a //= 255
    print(a)
b.append(a)
print(b)
print(count)

print((255**8)/(1024**3))
