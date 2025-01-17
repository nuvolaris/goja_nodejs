package process

import (
	"os"
	"strings"

	"github.com/nuvolaris/goja"
	"github.com/nuvolaris/goja_nodejs/require"
)

const ModuleName = "process"

type Process struct {
	env  map[string]string
	argv []string
}

func Require(runtime *goja.Runtime, module *goja.Object) {
	p := &Process{
		env:  make(map[string]string),
		argv: os.Args,
	}

	for _, e := range os.Environ() {
		envKeyValue := strings.SplitN(e, "=", 2)
		p.env[envKeyValue[0]] = envKeyValue[1]
	}

	o := module.Get("exports").(*goja.Object)
	o.Set("env", p.env)
	o.Set("argv", p.argv)
}

func Enable(runtime *goja.Runtime) {
	runtime.Set("process", require.Require(runtime, ModuleName))
}

func init() {
	require.RegisterCoreModule(ModuleName, Require)
}
