package plugo

// Generic plugin infrastructure
import (
	"log"
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
type PluginFactory func() (interface{}, error)

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

// Plugin Factories
var PluginFactories = &pluginFactory{
	newPlugin(new(PluginFactory)),
}

type pluginFactory struct {
	*plugin
}

func (p *pluginFactory) Unregister(name string) bool {
	return p.unregister(name)
}

func (p *pluginFactory) Register(component PluginFactory, name string) bool {
	return p.register(component, name)
}

func (p *pluginFactory) Lookup(name string) (PluginFactory, bool) {
	ext, ok := p.lookup(name)
	if !ok {
		return nil, ok
	}
	return ext.(PluginFactory), ok
}

func LoadPluginsWithConfig(cl *ConfigLoader, pluginConfigs []PluginConfig) []interface{} {
	plugins := make([]interface{}, len(pluginConfigs))
	for i, pluginConfig := range pluginConfigs {
		sf, ok := PluginFactories.Lookup(pluginConfig.Name)

		if !ok {
			log.Fatalf("Unable to load plugin with name: %s", pluginConfig.Name)
		}

		s, err := sf()
		if err != nil {
			log.Fatalf("Encountered error loading plugin: %v", err)
		}

		// apply config and validate
		err = cl.ApplyConfig(pluginConfig.Config, s)
		if err != nil {
			log.Fatalf("Encountered error applying configuration to plugin: %v", err)
		}
		err = cl.Validate(s)
		if err != nil {
			log.Fatalf("Encountered error validating plugin configuration: %v", err)
		}

		plugins[i] = s
	}

	return plugins
}
