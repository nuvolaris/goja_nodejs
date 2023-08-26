package nuv

import (
	"github.com/nuvolaris/goja"
	"github.com/nuvolaris/goja_nodejs/require"
)

const ModuleName = "nuv"

type Nuv struct {
	runtime *goja.Runtime
	scanner Scanner
}

type Scanner interface {
	readFile(string) (string, error)                 //read an entire file
	writeFile(string, string) error                  // write an entire file
	readDir(string) []string                         // read a folder and return an array of filenames
	toYaml(map[string]interface{}) (string, error)   // encode js object into a yaml string
	fromYaml(string) (map[string]interface{}, error) // decode a string assuming it is yaml in a js object
	scan(string, func(string) string) string         // walks the substree starting in root, execute a function for each folder
}

func Require(runtime *goja.Runtime, module *goja.Object) {
	requireWithScanner(defaultNuvScanner)(runtime, module)
}

func RequireWithScanner(scanner Scanner) require.ModuleLoader {
	return requireWithScanner(scanner)
}

func requireWithScanner(scanner Scanner) require.ModuleLoader {
	return func(runtime *goja.Runtime, module *goja.Object) {
		nuv := &Nuv{
			runtime: runtime,
			scanner: scanner,
		}

		o := module.Get("exports").(*goja.Object)
		o.Set("readFile", nuv.readJSFunc())
		o.Set("writeFile", nuv.writeJSFunc())
		o.Set("readDir", nuv.readDirJSFunc())
		o.Set("scan", nuv.scanJSFunc())
		o.Set("toYaml", nuv.toYamlJSFunc(nuv.runtime))
		o.Set("fromYaml", nuv.fromYamlJSFunc())
	}
}

func Enable(runtime *goja.Runtime) {
	runtime.Set("nuv", require.Require(runtime, ModuleName))
}

func init() {
	require.RegisterCoreModule(ModuleName, Require)
}

func (nuv *Nuv) readJSFunc() func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			panic(nuv.runtime.NewTypeError("readFile() requires one argument"))
		}
		arg := call.Argument(0).String()
		output, err := nuv.scanner.readFile(arg)
		if err != nil {
			panic(err)
		}
		return nuv.runtime.ToValue(output)
	}
}

func (nuv *Nuv) writeJSFunc() func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			panic(nuv.runtime.NewTypeError("writeFile() requires two arguments"))
		}
		arg1 := call.Argument(0).String()
		arg2 := call.Argument(1).String()
		err := nuv.scanner.writeFile(arg1, arg2)
		if err != nil {
			panic(err)
		}
		return nil
	}
}

func (nuv *Nuv) readDirJSFunc() func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			panic(nuv.runtime.NewTypeError("readDir() requires one argument"))
		}
		arg := call.Argument(0).String()
		output := nuv.scanner.readDir(arg)
		return nuv.runtime.ToValue(output)
	}
}

func (nuv *Nuv) toYamlJSFunc(rt *goja.Runtime) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			panic(nuv.runtime.NewTypeError("toYaml() requires one argument"))
		}

		inputObj, ok := call.Argument(0).Export().(map[string]interface{})
		if !ok {
			panic(nuv.runtime.NewTypeError("toYaml() requires an object as argument"))
		}

		output, err := nuv.scanner.toYaml(inputObj)
		if err != nil {
			panic(err)
		}
		return nuv.runtime.ToValue(output)
	}
}

func (nuv *Nuv) fromYamlJSFunc() func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			panic(nuv.runtime.NewTypeError("fromYaml() requires one argument"))
		}
		arg := call.Argument(0).String()
		outputObj, err := nuv.scanner.fromYaml(arg)
		if err != nil {
			panic(err)
		}

		return nuv.runtime.ToValue(outputObj)
	}
}

func (nuv *Nuv) scanJSFunc() func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			panic(nuv.runtime.NewTypeError("scan() requires two arguments"))
		}
		arg1 := call.Argument(0).String()
		arg2, ok := call.Argument(1).Export().(func(goja.FunctionCall) goja.Value)
		if !ok {
			panic(nuv.runtime.NewTypeError("scan() requires a function as second argument"))
		}

		f := func(path string) string {
			return arg2(goja.FunctionCall{
				This:      nuv.runtime.ToValue(nil),
				Arguments: []goja.Value{nuv.runtime.ToValue(path)},
			}).String()
		}

		output := nuv.scanner.scan(arg1, f)
		return nuv.runtime.ToValue(output)
	}
}
