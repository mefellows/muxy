package muxy

import (
  "github.com/mefellows/muxy/config"
)

// Middleware's are executed in stacked order before or after a Symptom,
// and are perfect for jobs like instrumentation. They are given a read-only copy
// of the runtime context and are executed asynchronously so cannot block a request/response
//

type Middleware interface {
	Configure(c *config.RawConfig)
}
