package main

import (
	"net"
	"strings"
)

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

func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

// user handle the message
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		// to see who is online
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":  " + "is there\n"
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// rename|newUsername
		newName := strings.Split(msg, "|")[1]
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMsg("Change username failed: the '" + newName + "' is unavailable.\n")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()

			this.Name = newName
			this.SendMsg("Change username succeed: '" + newName + "'\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		// to|who|message
		who := strings.Split(msg, "|")[1]
		if who == "" {
			this.SendMsg("unknown format!\n")
			return
		}
		user, ok := this.server.OnlineMap[who]
		if !ok {
			this.SendMsg("the user is not online!\n")
			return
		}
		content := strings.Split(msg, "|")[2]
		if content == "" {
			this.SendMsg("empty message!\n")
			return
		}
		user.SendMsg(this.Name + " said: " + content + "\n")

	} else {
		this.server.BroadCast(this, msg)

	}
}

// listen to the user channel, once receiving a message, send it to the client
func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		// write the message to the connection, after transforming it to the bytes
		this.conn.Write([]byte(msg + "\n"))
	}
}
