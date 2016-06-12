// Package executors provides a registry for all accepted
// executors.
package executors

import (
	"runtime"

	"github.com/xsb/dog/dog"
	"github.com/xsb/dog/executors/def"
)

func init() {
	dog.RegisterExecutor("sh", def.NewDefaultExecutor("sh"))
	dog.RegisterExecutor("bash", def.NewDefaultExecutor("bash"))
	dog.RegisterExecutor("python", def.NewDefaultExecutor("python"))
	dog.RegisterExecutor("ruby", def.NewDefaultExecutor("ruby"))

	switch runtime.GOOS {
	case "windows":
		dog.RegisterExecutor("system", def.NewDefaultExecutor("cmd"))
	default:
		dog.RegisterExecutor("system", def.NewDefaultExecutor("sh"))
	}
}
