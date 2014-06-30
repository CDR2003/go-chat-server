package server

import "net"

type user struct {
	connection net.Conn
}

func newUser(connection net.Conn) *user {
	return &user{connection}
}

func (u *user) sendMessage(message string) {
	u.connection.Write([]byte(message))
}
