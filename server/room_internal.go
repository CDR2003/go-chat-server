package server

import "container/list"
import "fmt"

type room struct {
	members *list.List
}

func newRoom() *room {
	return &room{list.New()}
}

func (r *room) addMember(user *User) {
	r.members.PushBack(user)
	fmt.Println("User added.")
}

func (r *room) removeMember(user *User) {
	for e := r.members.Front(); e != nil; e = e.Next() {
		if e.Value == user {
			r.members.Remove(e)
			fmt.Println("User removed.")
			return
		}
	}
}

func (r *room) getMembers() *list.List {
	return r.members
}

func (r *room) broadcast(message string) {
	for e := r.members.Front(); e != nil; e = e.Next() {
		member := e.Value.(*User)
		member.SendMessage(message)
	}
}
