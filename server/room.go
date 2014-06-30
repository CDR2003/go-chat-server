package server

import "container/list"
import "fmt"

type room_addMember struct {
	p0 *User
}

type room_removeMember struct {
	p0 *User
}

type room_getMembers struct {
	r *list.List
}

type room_broadcast struct {
	p0 string
}

type Room struct {
	obj          *room
	close        chan interface{}
	addMember    chan *room_addMember
	removeMember chan *room_removeMember
	getMembers   chan *room_getMembers
	broadcast    chan *room_broadcast
}

func (r *Room) loop() {
	for {
		select {
		case <-r.close:
			return
		case params := <-r.addMember:
			r.obj.addMember(params.p0)
		case params := <-r.removeMember:
			r.obj.removeMember(params.p0)
		case <-r.getMembers:
			r.getMembers <- &room_getMembers{r.obj.getMembers()}
		case params := <-r.broadcast:
			r.obj.broadcast(params.p0)
		}
	}
}

func NewRoom() *Room {
	r := new(Room)
	r.obj = newRoom()
	r.close = make(chan interface{})
	r.addMember = make(chan *room_addMember)
	r.removeMember = make(chan *room_removeMember)
	r.getMembers = make(chan *room_getMembers)
	r.broadcast = make(chan *room_broadcast)
	go r.loop()
	fmt.Println("Room started")
	return r
}

func (r *Room) Close() {
	r.close <- nil
	fmt.Println("Room closed")
}

func (r *Room) AddMember(user *User) {
	r.addMember <- &room_addMember{user}
}

func (r *Room) RemoveMember(user *User) {
	r.removeMember <- &room_removeMember{user}
}

func (r *Room) GetMembers() *list.List {
	r.getMembers <- nil
	returnValue := <-r.getMembers
	return returnValue.r
}

func (r *Room) Broadcast(message string) {
	r.broadcast <- &room_broadcast{message}
}
