package muxy

// Generic plugin infrastructure
// TODO: Move this into its own project?

import (
	"reflect"
	"runtime"
	"strings"
	"sync"
)

//
// Configurable plugin Factory: given a lookup key (protocol, name etc.)
// return a plugin, give it its' configuration and validate it.
//

// Extension type for adding new plugins
type ProxyFactory func(protocol string) (Middleware, error)
type MiddlewareFactory func() (Middleware, error)
type SymptomFactory func() (Symptom, error)

var registry = struct {
	sync.Mutex
	extpoints map[string]*plugin
}{
	extpoints: make(map[string]*plugin),
}

type plugin struct {
	sync.Mutex
	iface      reflect.Type
	components map[string]interface{}
}

func newPlugin(iface interface{}) *plugin {
	p := &plugin{
		iface:      reflect.TypeOf(iface).Elem(),
		components: make(map[string]interface{}),
	}
	registry.Lock()
	defer registry.Unlock()
	registry.extpoints[p.iface.Name()] = p
	return p
}

func (p *plugin) lookup(name string) (ext interface{}, ok bool) {
	p.Lock()
	defer p.Unlock()
	ext, ok = p.components[name]
	return
}

func (p *plugin) all() map[string]interface{} {
	p.Lock()
	defer p.Unlock()
	all := make(map[string]interface{})
	for k, v := range p.components {
		all[k] = v
	}
	return all
}

func (p *plugin) register(component interface{}, name string) bool {
	p.Lock()
	defer p.Unlock()
	if name == "" {
		comType := reflect.TypeOf(component)
		if comType.Kind() == reflect.Func {
			nameParts := strings.Split(runtime.FuncForPC(
				reflect.ValueOf(component).Pointer()).Name(), ".")
			name = nameParts[len(nameParts)-1]
		} else {
			name = comType.Elem().Name()
		}
	}
	_, exists := p.components[name]
	if exists {
		return false
	}
	p.components[name] = component
	return true
}

func (p *plugin) unregister(name string) bool {
	p.Lock()
	defer p.Unlock()
	_, exists := p.components[name]
	if !exists {
		return false
	}
	delete(p.components, name)
	return true
}

func implements(component interface{}) []string {
	var ifaces []string
	typ := reflect.TypeOf(component)
	for name, p := range registry.extpoints {
		if p.iface.Kind() == reflect.Func && typ.AssignableTo(p.iface) {
			ifaces = append(ifaces, name)
		}
		if p.iface.Kind() != reflect.Func && typ.Implements(p.iface) {
			ifaces = append(ifaces, name)
		}
	}
	return ifaces
}

func Register(component interface{}, name string) []string {
	registry.Lock()
	defer registry.Unlock()
	var ifaces []string
	for _, iface := range implements(component) {
		if ok := registry.extpoints[iface].register(component, name); ok {
			ifaces = append(ifaces, iface)
		}
	}
	return ifaces
}

func Unregister(name string) []string {
	registry.Lock()
	defer registry.Unlock()
	var ifaces []string
	for iface, extpoint := range registry.extpoints {
		if ok := extpoint.unregister(name); ok {
			ifaces = append(ifaces, iface)
		}
	}
	return ifaces
}

// SymptomFactory

var SymptomFactories = &symptomFactory{
	newPlugin(new(SymptomFactory)),
}

type symptomFactory struct {
	*plugin
}

func (p *symptomFactory) Unregister(name string) bool {
	return p.unregister(name)
}

func (p *symptomFactory) Register(component SymptomFactory, name string) bool {
	return p.register(component, name)
}

func (p *symptomFactory) Lookup(name string) (SymptomFactory, bool) {
	ext, ok := p.lookup(name)
	if !ok {
		return nil, ok
	}
	return ext.(SymptomFactory), ok
}

func (p *symptomFactory) All() map[string]SymptomFactory {
	all := make(map[string]SymptomFactory)
	for k, v := range p.all() {
		all[k] = v.(SymptomFactory)
	}
	return all
}

func (p *symptomFactory) Names() []string {
	var names []string
	for k := range p.all() {
		names = append(names, k)
	}
	return names
}
