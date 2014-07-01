package main

import "./server"
import "fmt"
import "net"
import "os"
import "os/signal"
import "sync"

func StartRead(connection net.Conn, ch chan []byte) {
	buffer := make([]byte, 1024)
	for {
		count, err := connection.Read(buffer)
		if err != nil {
			ch <- nil
			return
		}
		ch <- buffer[:count]
	}
}

func HandleConnection(room *server.Room, connection net.Conn, waitGroup *sync.WaitGroup, closeChan chan bool) {
	defer waitGroup.Done()
	defer connection.Close()

	user := server.NewUser(connection)
	defer user.Close()
	room.AddMember(user)
	defer room.RemoveMember(user)

	ch := make(chan []byte)
	go StartRead(connection, ch)

	for {
		select {
		case data := <-ch:
			if data == nil {
				return
			}
			room.Broadcast(string(data))
		case <-closeChan:
			return
		}
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

	closeChan := make(chan bool)

	var waitGroup sync.WaitGroup
	defer waitGroup.Wait()

	for {
		select {
		case conn := <-acceptChan:
			waitGroup.Add(1)
			go HandleConnection(room, conn, &waitGroup, closeChan)
		case <-interruptChan:
			closeChan <- true
			return
		}
	}
}

func main() {
	StartServer()
}
