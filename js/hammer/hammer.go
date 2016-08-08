//Package hammer provides implmentation of the hammerjs.github.io library
package hammer

import (
	"github.com/gopherjs/gopherjs/js"
	// "honnef.co/go/js/dom"
)

const (
	HAMMER = "Hammer"

	//keys
	MANAGER = "Manager"
	ON      = "on"
)

type Event struct {
	*js.Object
	AdditionalEvent   string       `js:"additionalEvent"`
	Angle             float64      `js:"angle"`
	Center            *js.Object   `js:"center"`
	ChangedPointers   []*js.Object `js:"changedPointers"`
	DeltaTime         int          `js:"deltaTime"`
	DeltaY            int          `js:"deltaY"`
	DeltaX            int          `js:"deltaX"`
	Direction         int          `js:"direction"`
	Distance          float64      `js:"distance"`
	EventType         int          `js:"eventType"`
	IsFinal           bool         `js:"isFinal"`
	IsFirst           bool         `js:"isFirst"`
	MaxPointers       int          `js:"maxPointers"`
	OffsetDirection   int          `js:"offsetDirection"`
	OverallVelocity   float64      `js:"overallVelocity"`
	OverallVelocityX  float64      `js:"overallVelocityX"`
	OverallVelocityY  float64      `js:"overallVelocityY"`
	PointerType       string       `js:"pointerType"`
	Pointers          []*js.Object `js:"pointers"`
	PreventionDefault *js.Object   `js:"preventionDefault"`
	Rotation          int          `js:"rotation"`
	Scale             int          `js:"scale"`
	SrcEvent          *js.Object   `js:"srcEvent"`
	Target            *js.Object   `js:"target"`
	Timestamp         int          `js:"timeStamp"`
	Type              string       `js:"type"`
	Velocity          float64      `js:"velocity"`
	VelocityX         float64      `js:"velocityX"`
	VelocityY         float64      `js:"velocityY"`
}

type Hammer struct {
	manager *js.Object
}

type EventCallback func(Event)

func (self *Hammer) New(i interface{}) {
	self.manager = js.Global.Get(HAMMER).New(i)
}

func (self *Hammer) On(i ...interface{}) *Hammer {

	self.manager.Call(ON, i...)
	return self
}
