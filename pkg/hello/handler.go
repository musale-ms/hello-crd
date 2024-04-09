package hello

import (
	"fmt"

	"k8s.io/client-go/tools/cache"
)

type HelloEvent interface {
	AddEvent() func(obj interface{})
	UpdateEvent() func(obj interface{})
	DeleteEvent() func(obj interface{})
}

type HelloEventHandler struct{}

func (h HelloEventHandler) NewEvent() cache.ResourceEventHandlerFuncs {
	return cache.ResourceEventHandlerFuncs{
		AddFunc:    h.AddEvent(),
		UpdateFunc: h.UpdateEvent(),
		DeleteFunc: h.DeleteEvent(),
	}
}

func (h HelloEventHandler) AddEvent() func(obj interface{}) {
	return func(obevj interface{}) {
		fmt.Println("Added a ", obevj)
	}
}

func (h HelloEventHandler) UpdateEvent() func(oldobj, newobj interface{}) {
	return func(oldobj, newobj interface{}) {
		fmt.Println("Updated a ", oldobj, newobj)
	}
}

func (h HelloEventHandler) DeleteEvent() func(obj interface{}) {
	return func(obevj interface{}) {
		fmt.Println("Deleted a ", obevj)
	}
}
