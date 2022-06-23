# v1 - Basic Server
Build the server, and start listening:
```shell
go build -o server main.go server.go
./server
```
Start a new terminal, and use `nc` to test the connection:
```
nc 127.0.0.1 8888
```

# v2 - user on line and broadcasting
Build:
```
go build -o server main.go server.go user.go
```
Start a new ternimal to be on line:
```
nc 127.0.0.1 8888
```
the output should be:
```
[127.0.0.1:55682]127.0.0.1:55682:  ONLINE
```

# v3 - user broadcast message
> Now, user can input the message and broadcast it.

Build:
```
go build -o server main.go server.go user.go
```
Use `nc` to test:
```
nc 127.0.0.1 8888
Hi, guys!
```
