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

Use `nc` to test:
```
nc 127.0.0.1 8888
Hi, guys!
```

# v4 - Integrate the user logic
> In fact, the `Server.Handler()` comprises the tasks: *user online*, *user offline*, *user message*. These task is user-specific, and we should package them to the `User`. That's what we'll do in this section.

We added three methods to `User`: `User.Online()`, `User.Offline()`, `User.DoMessage()`

# v5 - Query online users
>  Now, User can use 'who' command to query who is online.

Test:
```
nc 127.0.0.1 8888
who
```
# v6 - Change user name
> Now, user can change his name with input `rename|newUsername`.


Test:
```
nc 127.0.0.1 8888
rename|Jack
```

# v7 - Timeout
> If a user is not live ( don't send any message ) for some certain time, he will be kicked out and the connection will be closed.

We use a channel to simulate the timer.


# v8 - Private chat
> User can send the private messasge to specified user, through the command `to|username|message`.

Test:
```
to|jack|hi, Jack!
```
Note that the specified user must be online. 