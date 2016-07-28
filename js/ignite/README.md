#Ignite

Ignite is a go package that is dependent on the gopherjs project.  Ignite allows you to easily mark up your html DOM elements with class and data attributes to bind either one way or two way updates to your go structs.  

##Installation

	go get -u github.com/gopherjs/gopherjs
	go get github.com/cloud-ignite/go/js/ignite

##Design
Ignite is based on a virtual dom design.  The virtual dom implements an observer pattern on golang objects and uses reflection to access and update the objects.  

Observable objects must be implemented by the caller and added to the virtual dom before binding to dom elements occur.  Subsequent additions of observable objects and bindings can be performed at runtime as well.

##Requirements

Observable golang objects must implement the following interface to adhere to the virtual dom's requirements.

	type ObservableHandler func(string)

	type Observable interface {
		AddObserver(ObservableHandler)
		Commit()
		Key() string
	}


##HTML Elements

To bind to a golang observable object you must mark up your elements with the following class:  

	class = "ignite"
	
This tells the virtual dom which dom elememts to store in memory.  Ignite then looks for specific bindings based on data attributes of the element.  The following binds to the innerhtml of a golang object with a Key of "menu".  Only public properties with starting capital letters are available for binding based on golang reflection rules. 

	data-bind-innerhtml = "menu.Home.Title"

The following binding types are supported.

###Binding Types

	data-bind-innerhtml  :  Sets the innerhtml of the element
	
###Examples

See example directory.
	

