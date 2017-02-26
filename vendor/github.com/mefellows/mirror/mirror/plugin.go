package mirror

import (
	"github.com/mefellows/mirror/filesystem"
	"reflect"
	"runtime"
	"strings"
	"sync"
)

// Extension type for adding new FileSystem adapters
type FileSystemFactory func(protocol string) (filesystem.FileSystem, error)

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

// FileSystemFactory

var FileSystemFactories = &fileSystemFactory{
	newPlugin(new(FileSystemFactory)),
}

type fileSystemFactory struct {
	*plugin
}

func (p *fileSystemFactory) Unregister(name string) bool {
	return p.unregister(name)
}

func (p *fileSystemFactory) Register(component FileSystemFactory, name string) bool {
	return p.register(component, name)
}

func (p *fileSystemFactory) Lookup(name string) (FileSystemFactory, bool) {
	ext, ok := p.lookup(name)
	if !ok {
		return nil, ok
	}
	return ext.(FileSystemFactory), ok
}

func (p *fileSystemFactory) All() map[string]FileSystemFactory {
	all := make(map[string]FileSystemFactory)
	for k, v := range p.all() {
		all[k] = v.(FileSystemFactory)
	}
	return all
}

func (p *fileSystemFactory) Names() []string {
	var names []string
	for k := range p.all() {
		names = append(names, k)
	}
	return names
}
