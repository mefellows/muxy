package plugo

// A configurable is something that can be given an RawConfig
// object, validate itself against the given configuration and return
// an error if that check fails
type Configurable interface {
	Configure(c *RawConfig) error
}
