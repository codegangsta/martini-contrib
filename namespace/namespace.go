// Package namespace allows you to register namespaces for your martini applications
package namespace

import (
	"reflect"

	"github.com/codegangsta/martini"
)

type MartiniNamespace struct {
	*martini.ClassicMartini
	Namespace string
	Handlers []martini.Handler
}

func Namespace(m *martini.ClassicMartini, namespace string, handlers ...martini.Handler) {
	namespaceFunc := handlers[len(handlers)-1]
	handlers       = handlers[:len(handlers)-1]
	ns            := newNamespace(m, namespace, handlers)
	arguments     := []reflect.Value{reflect.ValueOf(ns)}
	
	reflect.ValueOf(namespaceFunc).Call(arguments)
}

func (ns *MartiniNamespace) Get(path string, handlers ...martini.Handler) martini.Route {
	return ns.ClassicMartini.Get(ns.Namespace + path, append(ns.Handlers, handlers...)...)
}

func (ns *MartiniNamespace) Patch(path string, handlers ...martini.Handler) martini.Route {
	return ns.ClassicMartini.Patch(ns.Namespace + path, append(ns.Handlers, handlers...)...)
}

func (ns *MartiniNamespace) Post(path string, handlers ...martini.Handler) martini.Route {
	return ns.ClassicMartini.Post(ns.Namespace + path, append(ns.Handlers, handlers...)...)
}

func (ns *MartiniNamespace) Put(path string, handlers ...martini.Handler) martini.Route {
	return ns.ClassicMartini.Put(ns.Namespace + path, append(ns.Handlers, handlers...)...)
}

func (ns *MartiniNamespace) Delete(path string, handlers ...martini.Handler) martini.Route {
	return ns.ClassicMartini.Delete(ns.Namespace + path, append(ns.Handlers, handlers...)...)
}

func (ns *MartiniNamespace) Options(path string, handlers ...martini.Handler) martini.Route {
	return ns.ClassicMartini.Options(ns.Namespace + path, append(ns.Handlers, handlers...)...)
}

func (ns *MartiniNamespace) Head(path string, handlers ...martini.Handler) martini.Route {
	return ns.ClassicMartini.Head(ns.Namespace + path, append(ns.Handlers, handlers...)...)
}

func (ns *MartiniNamespace) Any(path string, handlers ...martini.Handler) martini.Route {
	return ns.ClassicMartini.Any(ns.Namespace + path, append(ns.Handlers, handlers...)...)
}

func (ns *MartiniNamespace) NotFound(handlers ...martini.Handler) martini.Route {
	return ns.ClassicMartini.Any(ns.Namespace + "**", handlers...)
}

func newNamespace(m *martini.ClassicMartini, namespace string, handlers []martini.Handler) *MartiniNamespace {
	return &MartiniNamespace{
		m,
		namespace,
		handlers,
	}
}