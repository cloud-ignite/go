//Package mobx provides implmentation of the http://mobxjs.github.io/ library
package mobx

import (
	"github.com/gopherjs/gopherjs/js"
	// "honnef.co/go/js/dom"
)

const (
	MOBX = "mobx"

	//keys
	OBSERVABLE  = "observable"
	AUTORUN     = "autorun"
	REACTION    = "reaction"
	TRANSACTION = "transaction"
	ACTION      = "action"
)

type MobX struct {
}

func (self *MobX) Observable(i interface{}) *js.Object {
	return self.mobx().Call(OBSERVABLE, i)
}

func (self *MobX) Autorun(i interface{}) {
	self.mobx().Call(AUTORUN, i)
}

func (self *MobX) mobx() *js.Object {

	return js.Global.Get("mobx")
}
