package model

import "github.com/davyxu/cellnet"

type Backend struct {
	SvcName string
	Session cellnet.Session
}

type User struct {
	Targets []*Backend
}

func (self *User) AddBackend(svcName string, ses cellnet.Session) {
	self.Targets = append(self.Targets, &Backend{
		SvcName: svcName,
		Session: ses,
	})
}

func (self *User) SetBackend(svcName string, ses cellnet.Session) {

	for _, t := range self.Targets {
		if t.SvcName == svcName {
			t.Session = ses
			return
		}
	}
}

func (self *User) GetBackend(svcName string) cellnet.Session {

	for _, t := range self.Targets {
		if t.SvcName == svcName {
			return t.Session
		}
	}

	return nil
}

func NewUser() *User {
	return &User{}
}
