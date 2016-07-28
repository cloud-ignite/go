package main

import (
	"fmt"
	"github.com/cloud-ignite/go/js/ignite"
	"honnef.co/go/js/dom"
	"time"
)

type Menu struct {
	callbacks []ignite.ObservableHandler
	Home      struct {
		Title string
	}
}

func (self *Menu) AddObserver(handler ignite.ObservableHandler) {
	self.callbacks = append(self.callbacks, handler)
}

func (self *Menu) Commit() {
	if self == nil {
		return
	}
	for _, callback := range self.callbacks {
		callback("menu")
	}
}

func (self *Menu) Key() string {
	return "menu"
}

var window dom.Window
var document dom.Document
var menu Menu

var vDom ignite.VirtualDom

func main() {
	fmt.Println("Welcome to Ignite")

	window = dom.GetWindow()
	document = dom.GetWindow().Document()

	window.AddEventListener("load", true, onload)
	vDom.AddObservable(&menu)
	menu.Home.Title = "Home"
}

func onload(ev dom.Event) {

	vDom.BindDocument(document)

	go func() {
		for {
			time.Sleep(30 * time.Millisecond)
			menu.Home.Title = time.Now().String()
			menu.Commit()
		}
	}()
}
