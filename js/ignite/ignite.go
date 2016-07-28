//Package ignite provides implmentation of a virtual dom to interact with html dom elements.
package ignite

import (
	// "fmt"
	"honnef.co/go/js/dom"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

const VIRTUAL_DOM_INNERHTML = "data-bind-innerhtml"

type ObservableHandler func(string)

//Observable methods for your structs are required for the VirtualDom to operate properly.
//During the binding process to Dom elements the Key() is the first returned element.
type Observable interface {
	AddObserver(ObservableHandler)
	Commit()
	Key() string
}

type Attribute struct {
	M    sync.RWMutex
	Keys map[string]string
}

type Binding struct {
	M          sync.RWMutex
	Element    dom.Element
	Attributes map[string]Attribute
}

type VirtualDom struct {
	M                 sync.RWMutex
	ObservableObjects map[string]Observable
	Bindings          []Binding
}

func (self *VirtualDom) AddObservable(object Observable) {
	if self.ObservableObjects == nil {
		self.ObservableObjects = make(map[string]Observable)
	}
	self.M.Lock()
	self.ObservableObjects[object.Key()] = object
	self.M.Unlock()

	object.AddObserver(self.handleObservableCallback)
}

func (self *VirtualDom) handleObservableCallback(key string) {

	//First Get the Object observed to change
	obj := self.ObservableObjects[key]

	if obj == nil {
		return
	}

	//Do this concurrently

	for _, binding := range self.Bindings {
		attribute, ok := binding.Attributes[key]
		if ok {
			for key, value := range attribute.Keys {

				// fmt.Println("printing the Property")
				// fmt.Printf("%v\n", value)

				//Use Reflection to get the value of the obj
				propertyValue := getReflectionValue(value, obj)

				switch key {
				case VIRTUAL_DOM_INNERHTML:
					binding.Element.SetInnerHTML(propertyValue)
					break
				}
			}
		}
	}

}

func (self *VirtualDom) BindDocument(document dom.Document) error {

	//Do this concurrently

	for _, element := range document.GetElementsByClassName("gopherjs") {

		var b Binding
		b.Element = element
		b.Attributes = make(map[string]Attribute)

		for key, value := range element.Attributes() {

			var add bool

			switch key {
			case VIRTUAL_DOM_INNERHTML:
				add = true
				break
			}

			if add {

				reflectionKeys := strings.Split(value, ".")
				observableObj := reflectionKeys[0]
				reflectionKey := strings.Replace(value, observableObj+".", "", 1)

				attribute, ok := b.Attributes[observableObj]

				if ok {
					attribute.Keys[key] = reflectionKey

				} else {
					attribute.Keys = make(map[string]string)
					attribute.Keys[key] = reflectionKey

					b.Attributes[observableObj] = attribute
				}
			}
		}

		self.M.Lock()
		self.Bindings = append(self.Bindings, b)
		self.M.Unlock()

	}

	// fmt.Printf("%v\n", self.Bindings[0].Attributes)

	return nil
}

func getReflectionValue(key string, x interface{}) string {

	splitKey := strings.Split(key, ".")
	propertyName := splitKey[0]

	val := reflect.ValueOf(x).Elem()

	for i := 0; i < val.NumField(); i++ {

		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		name := typeField.Name

		if name != propertyName {
			continue
		}

		if !valueField.CanInterface() {
			continue
		}

		f := valueField.Interface()
		val := reflect.ValueOf(f)

		switch val.Kind() {

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return strconv.FormatInt(val.Int(), 10)
		case reflect.Float32:
			return strconv.FormatFloat(val.Float(), 'E', -1, 32)
		case reflect.Float64:
			return strconv.FormatFloat(val.Float(), 'E', -1, 64)
		case reflect.String:
			return val.String()
		case reflect.Struct:
			if len(splitKey) > 1 {
				return getStructReflectionValue(strings.Replace(key, propertyName+".", "", 1), val)
			}

		}
	}
	return ""
}

func getStructReflectionValue(key string, val reflect.Value) string {

	splitKey := strings.Split(key, ".")
	propertyName := splitKey[0]

	for i := 0; i < val.NumField(); i++ {

		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		name := typeField.Name

		if name != propertyName {
			continue
		}

		if !valueField.CanInterface() {
			continue
		}

		f := valueField.Interface()
		val := reflect.ValueOf(f)

		switch val.Kind() {

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return strconv.FormatInt(val.Int(), 10)
		case reflect.Float32:
			return strconv.FormatFloat(val.Float(), 'E', -1, 32)
		case reflect.Float64:
			return strconv.FormatFloat(val.Float(), 'E', -1, 64)
		case reflect.String:
			return val.String()
		case reflect.Struct:
			if len(splitKey) > 1 {
				return getStructReflectionValue(strings.Replace(key, propertyName+".", "", 1), val)
			}

		}
	}
	return ""
}
