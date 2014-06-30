package main

import "./server"
import "fmt"
import "net"
import "os"
import "os/signal"

func HandleConnection(room *server.Room, connection net.Conn) {
	buffer := make([]byte, 1024)
	for {
		count, err := connection.Read(buffer)
		if err != nil {
			return
		}
		room.Broadcast(string(buffer[:count]))
	}
}

func StartAccept(socket net.Listener, ch chan net.Conn) {
	for {
		conn, err := socket.Accept()
		if err != nil {
			continue
		}
		ch <- conn
	}
}

func StartServer() {
	room := server.NewRoom()
	defer room.Close()

	socket, err := net.Listen("tcp", ":12345")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer socket.Close()

	fmt.Println("Server started at :12345. Press Ctrl-C to stop.")

	interruptChan := make(chan os.Signal)
	signal.Notify(interruptChan, os.Interrupt)

	acceptChan := make(chan net.Conn)
	go StartAccept(socket, acceptChan)

	for {
		select {
		case conn := <-acceptChan:
			user := server.NewUser(conn)
			defer user.Close()
			room.AddMember(user)
			defer room.RemoveMember(user)
			go HandleConnection(room, conn)
		case <-interruptChan:
			return
		}
	}
}

func main() {
	StartServer()
}
