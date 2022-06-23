package main

import "net"

type User struct {
	Name string      // user name
	Addr string      // address
	C    chan string // user channel
	conn net.Conn    //	user connection
}

// user constructor
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}

	// start a goroutine, to listen to the user channel
	go user.ListenMessage()

	return user
}

// listen to the user channel, once receiving a message, send it to the client
func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		// write the message to the connection, after transforming it to the bytes
		this.conn.Write([]byte(msg + "\n"))
	}
}
