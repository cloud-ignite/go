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

const (
	VIRTUAL_DOM_INNERHTML = "data-bind-innerhtml"
	VIRTUAL_DOM_TOP       = "data-bind-top"
	VIRTUAL_DOM_LEFT      = "data-bind-left"
)

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

func (self *VirtualDom) BindDocument(document dom.Document) error {

	//Do this concurrently

	for _, element := range document.GetElementsByClassName("ignite") {

		var b Binding
		b.Element = element
		b.Attributes = make(map[string]Attribute)

		for key, value := range element.Attributes() {

			var add bool

			switch key {
			case VIRTUAL_DOM_INNERHTML, VIRTUAL_DOM_LEFT, VIRTUAL_DOM_TOP:
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

	//Loop through the observable elemements and bind.

	for _, obj := range self.ObservableObjects {
		self.bindDOM(obj.Key())
	}
	// fmt.Printf("%v\n", self.Bindings[0].Attributes)

	return nil
}

func (self *VirtualDom) handleObservableCallback(key string) {
	self.bindDOM(key)
}

func (self *VirtualDom) bindDOM(key string) {

	// fmt.Println("BindingDom")
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

				//Use Reflection to get the value of the obj
				propertyValue := getReflectionValue(value, obj)
				htmlElement := binding.Element.(dom.HTMLElement)

				switch key {
				case VIRTUAL_DOM_INNERHTML:
					binding.Element.SetInnerHTML(propertyValue)
					break
				case VIRTUAL_DOM_LEFT:
					cssUnit := getCSSUnit(htmlElement.Style().GetPropertyValue("left"))
					htmlElement.Style().SetProperty("left", propertyValue+cssUnit, "important")
					break
				case VIRTUAL_DOM_TOP:
					cssUnit := getCSSUnit(htmlElement.Style().GetPropertyValue("top"))
					htmlElement.Style().SetProperty("top", propertyValue+cssUnit, "important")
					break
				}
			}
		}
	}
}

func getReflectionValue(key string, x interface{}) string {

	splitKey := strings.Split(key, ".")

	propertyName := splitKey[0]
	originalPropertyName := propertyName

	arrayIndex := strings.Index(propertyName, "[")
	if arrayIndex != -1 {
		propertyName = propertyName[:arrayIndex]
	}

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
		case reflect.Array, reflect.Slice:

			lastIndex := strings.Index(originalPropertyName, "]")
			index := StringToInt(originalPropertyName[arrayIndex+1 : lastIndex])
			arrayItem := val.Index(index)

			switch arrayItem.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return strconv.FormatInt(arrayItem.Int(), 10)
			case reflect.Float32:
				return strconv.FormatFloat(arrayItem.Float(), 'E', -1, 32)
			case reflect.Float64:
				return strconv.FormatFloat(arrayItem.Float(), 'E', -1, 64)
			case reflect.String:
				return arrayItem.String()

			case reflect.Struct:
				if len(splitKey) > 1 {
					return getStructReflectionValue(strings.Replace(key, originalPropertyName+".", "", 1), arrayItem)
				}
			}

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
	originalPropertyName := propertyName

	arrayIndex := strings.Index(propertyName, "[")
	if arrayIndex != -1 {
		propertyName = propertyName[:arrayIndex]
	}

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
		case reflect.Array, reflect.Slice:

			lastIndex := strings.Index(originalPropertyName, "]")
			index := StringToInt(originalPropertyName[arrayIndex+1 : lastIndex])
			arrayItem := val.Index(index)

			switch arrayItem.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return strconv.FormatInt(arrayItem.Int(), 10)
			case reflect.Float32:
				return strconv.FormatFloat(arrayItem.Float(), 'E', -1, 32)
			case reflect.Float64:
				return strconv.FormatFloat(arrayItem.Float(), 'E', -1, 64)
			case reflect.String:
				return arrayItem.String()

			case reflect.Struct:
				if len(splitKey) > 1 {
					return getStructReflectionValue(strings.Replace(key, originalPropertyName+".", "", 1), arrayItem)
				}
			}
		case reflect.Struct:
			if len(splitKey) > 1 {
				return getStructReflectionValue(strings.Replace(key, propertyName+".", "", 1), val)
			}

		}
	}
	return ""
}

func getCSSUnit(s string) string {
	if strings.Contains(s, "px") {
		return "px"
	} else if strings.Contains(s, "em") {
		return "em"
	} else if strings.Contains(s, "ex") {
		return "ex"
	} else if strings.Contains(s, "ch") {
		return "ch"
	} else if strings.Contains(s, "rem") {
		return "rem"
	} else if strings.Contains(s, "vw") {
		return "vw"
	} else if strings.Contains(s, "vh") {
		return "vh"
	} else if strings.Contains(s, "vmin") {
		return "vmin"
	} else if strings.Contains(s, "vmax") {
		return "vmax"
	} else if strings.Contains(s, "%") {
		return "%"
	} else if strings.Contains(s, "cm") {
		return "cm"
	} else if strings.Contains(s, "mm") {
		return "mm"
	} else if strings.Contains(s, "in") {
		return "in"
	} else if strings.Contains(s, "pt") {
		return "pt"
	} else if strings.Contains(s, "pc") {
		return "pc"
	}

	return "px"
}

func StringToInt(val string) int {

	r, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return r
}

func IntToString(val int) string {
	return strconv.Itoa(val)
}
