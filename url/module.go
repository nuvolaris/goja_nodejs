package url

import (
	"math"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/net/idna"

	"github.com/nuvolaris/goja"
	"github.com/nuvolaris/goja_nodejs/require"
)

const ModuleName = "url"

type urlModule struct {
	r *goja.Runtime

	URLSearchParamsPrototype         *goja.Object
	URLSearchParamsIteratorPrototype *goja.Object
}

func Require(runtime *goja.Runtime, module *goja.Object) {
	exports := module.Get("exports").(*goja.Object)
	m := &urlModule{
		r: runtime,
	}
	exports.Set("URL", m.createURLConstructor())
	exports.Set("URLSearchParams", m.createURLSearchParamsConstructor())
	exports.Set("domainToASCII", m.domainToASCII)
	exports.Set("domainToUnicode", m.domainToUnicode)
}

func Enable(runtime *goja.Runtime) {
	m := require.Require(runtime, ModuleName).ToObject(runtime)
	runtime.Set("URL", m.Get("URL"))
	runtime.Set("URLSearchParams", m.Get("URLSearchParams"))
}

func init() {
	require.RegisterCoreModule(ModuleName, Require)
}
