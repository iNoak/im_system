package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
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
	user := NewUser(conn, this)
	user.Online()

	isLive := make(chan bool)

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 { // if use ctrl+C to close the connection, the return n is 0
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}

			msg := string(buf[:n-1])
			user.DoMessage(msg)

			isLive <- true
		}

	}()
	for {
		select {
		case <-isLive:
			// the user is live, reset the timer
		case <-time.After(time.Second * 100):
			// time out , kick out
			user.SendMsg("you are kicked out!\n")
			close(user.C)
			conn.Close()
		}

	}
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
