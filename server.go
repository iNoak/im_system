package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string // IP address
	Port int    // port

	OnlineMap map[string]*User // store the online users
	mapLock   sync.RWMutex     // the mutex for online map

	Message chan string // the channel for broadcasting the message
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

// listen to the Message
func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap { // send it to all clients
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}

}

// broadcast the message to all the user online !
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":  " + msg
	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	user := NewUser(conn)

	// add user to the online map
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()

	// broadcast that a user is online!
	this.BroadCast(user, "ONLINE")

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				this.BroadCast(user, "OFFLINE")
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}

			msg := string(buf[:n-1])
			this.BroadCast(user, msg)
		}

	}()

	select {}
}

func (this *Server) Start() {
	// a socket listener
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port)) // connection type, address

	if err != nil {
		fmt.Println("Server Listen err: ", err)
		return
	}
	defer listener.Close()

	go this.ListenMessage()

	for {
		// accpet
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}

		// start a go routine, to do the handler
		go this.Handler(conn)
	}
}
