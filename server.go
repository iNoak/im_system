package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}

	return server
}

func (this *Server) Handler(conn net.Conn) {
	fmt.Println("connection succeed!")
}

func (this *Server) Start() {
	// a socket listener
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port)) // connection type, address

	if err != nil {
		fmt.Println("Server Listen err:", err)
		return
	}
	defer listener.Close()

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
