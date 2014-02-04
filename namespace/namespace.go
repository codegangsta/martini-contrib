// Package namespace allows you to register namespaces for your martini applications
package namespace

import (
	"reflect"

	"github.com/codegangsta/martini"
)

// The MartiniNamespace is a wrapper around ClassicMartini. It implements the Router interface.
type MartiniNamespace struct {
	*martini.ClassicMartini
	Namespace string
	Handlers []martini.Handler
}

// namespace accepts your currenct ClassicMartini instance, the namespace as a string
// as well as any Handlers that should be called for the whole namespace. They will always be
// called before the handlers defined in the routes of the namespace.
func Namespace(m *martini.ClassicMartini, namespace string, handlers ...martini.Handler) {
	namespaceFunc := handlers[len(handlers)-1]
	handlers       = handlers[:len(handlers)-1]
	ns            := newNamespace(m, namespace, handlers)
	arguments     := []reflect.Value{reflect.ValueOf(ns)}
	
	reflect.ValueOf(namespaceFunc).Call(arguments)
}

// The namespace GET handler, call it like you would in martini
func (ns *MartiniNamespace) Get(path string, handlers ...martini.Handler) martini.Route {
	return ns.ClassicMartini.Get(ns.Namespace + path, append(ns.Handlers, handlers...)...)
}

// The namespace PATCH handler, call it like you would in martini
func (ns *MartiniNamespace) Patch(path string, handlers ...martini.Handler) martini.Route {
	return ns.ClassicMartini.Patch(ns.Namespace + path, append(ns.Handlers, handlers...)...)
}

// The namespace POST handler, call it like you would in martini
func (ns *MartiniNamespace) Post(path string, handlers ...martini.Handler) martini.Route {
	return ns.ClassicMartini.Post(ns.Namespace + path, append(ns.Handlers, handlers...)...)
}

// The namespace PUT handler, call it like you would in martini
func (ns *MartiniNamespace) Put(path string, handlers ...martini.Handler) martini.Route {
	return ns.ClassicMartini.Put(ns.Namespace + path, append(ns.Handlers, handlers...)...)
}

// The namespace DELETE handler, call it like you would in martini
func (ns *MartiniNamespace) Delete(path string, handlers ...martini.Handler) martini.Route {
	return ns.ClassicMartini.Delete(ns.Namespace + path, append(ns.Handlers, handlers...)...)
}

// The namespace OPTIONS handler, call it like you would in martini
func (ns *MartiniNamespace) Options(path string, handlers ...martini.Handler) martini.Route {
	return ns.ClassicMartini.Options(ns.Namespace + path, append(ns.Handlers, handlers...)...)
}

// The namespace HEAD handler, call it like you would in martini
func (ns *MartiniNamespace) Head(path string, handlers ...martini.Handler) martini.Route {
	return ns.ClassicMartini.Head(ns.Namespace + path, append(ns.Handlers, handlers...)...)
}

// The namespace handler for any HTTP method, call it like you would in martini
func (ns *MartiniNamespace) Any(path string, handlers ...martini.Handler) martini.Route {
	return ns.ClassicMartini.Any(ns.Namespace + path, append(ns.Handlers, handlers...)...)
}

// NotFound Routing in a namespace works via a globbed route atop of the namespace.
// Therefore, NotFound routing must always be installed as the last function of a namespace
func (ns *MartiniNamespace) NotFound(handlers ...martini.Handler) martini.Route {
	return ns.ClassicMartini.Any(ns.Namespace + "**", handlers...)
}

// Return a new namespace
func newNamespace(m *martini.ClassicMartini, namespace string, handlers []martini.Handler) *MartiniNamespace {
	return &MartiniNamespace{
		m,
		namespace,
		handlers,
	}
}