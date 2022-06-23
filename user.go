package main

import "net"

type User struct {
	Name   string      // user name
	Addr   string      // address
	C      chan string // user channel
	conn   net.Conn    //	user connection
	server *Server     // belongs to which server
}

// user constructor
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	// start a goroutine, to listen to the user channel
	go user.ListenMessage()

	return user
}

// user online
func (this *User) Online() {

	// add user to the online map
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	// broadcast that a user is online!
	this.server.BroadCast(this, "ONLINE")
}

// user offline
func (this *User) Offline() {
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	// broadcast that a user is online!
	this.server.BroadCast(this, "OFFLINE")
}

// user handle the message
func (this *User) DoMessage(msg string) {
	this.server.BroadCast(this, msg)
}

// listen to the user channel, once receiving a message, send it to the client
func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		// write the message to the connection, after transforming it to the bytes
		this.conn.Write([]byte(msg + "\n"))
	}
}
