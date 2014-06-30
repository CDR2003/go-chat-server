package server

import "net"
import "fmt"

type user_sendMessage struct {
	p0 string
}

type User struct {
	obj         *user
	close       chan interface{}
	sendMessage chan *user_sendMessage
}

func (u *User) loop() {
	for {
		select {
		case <-u.close:
			return
		case params := <-u.sendMessage:
			u.obj.sendMessage(params.p0)
		}
	}
}

func NewUser(connection net.Conn) *User {
	u := new(User)
	u.obj = newUser(connection)
	u.close = make(chan interface{})
	u.sendMessage = make(chan *user_sendMessage)
	go u.loop()
	fmt.Println("User " + connection.RemoteAddr().String() + " started")
	return u
}

func (u *User) Close() {
	u.close <- nil
	fmt.Println("User " + u.obj.connection.RemoteAddr().String() + " closed")
}

func (u *User) SendMessage(message string) {
	u.sendMessage <- &user_sendMessage{message}
}
