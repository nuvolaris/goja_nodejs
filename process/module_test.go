package process

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/nuvolaris/goja"
	"github.com/nuvolaris/goja_nodejs/require"
)

func TestProcessArgvStructure(t *testing.T) {
	vm := goja.New()

	new(require.Registry).Enable(vm)
	Enable(vm)

	if c := vm.Get("process"); c == nil {
		t.Fatal("process not found")
	}

	if c, err := vm.RunString("process.argv"); c == nil || err != nil {
		t.Fatal("error accessing process.argv")
	}

	if c, err := vm.RunString("process.argv.length"); c == nil || err != nil {
		t.Fatal("error accessing process.argv.length")
	}
}

func TestProcessArgvValues(t *testing.T) {
	vm := goja.New()

	new(require.Registry).Enable(vm)
	Enable(vm)

	if c, err := vm.RunString("process.argv[0]"); c == nil || err != nil {
		t.Fatal("error accessing process.argv[0]")
	}

	if c, err := vm.RunString("process.argv[1]"); c == nil || err != nil {
		t.Fatal("error accessing process.argv[1]")
	}

	if c, err := vm.RunString("process.argv[2]"); c == nil || err != nil {
		t.Fatal("error accessing process.argv[2]")
	}
}

func TestProcessArgvValuesArtificial(t *testing.T) {
	vm := goja.New()

	new(require.Registry).Enable(vm)
	Enable(vm)

	jsRes, err := vm.RunString("process.argv[0]")
	if err != nil {
		t.Fatal(fmt.Sprintf("Error executing: %s", err))
	}

	if strings.Contains(jsRes.String(), "process.test") == false {
		t.Fatal(fmt.Sprintf("Error executing: got %s but expected %s", jsRes, "goja"))
	}

	jsRes, err = vm.RunString("process.argv[1]")
	if err != nil {
		t.Fatal(fmt.Sprintf("Error executing: %s", err))
	}

	if strings.Contains(jsRes.String(), "-test") == false {
		t.Fatal(fmt.Sprintf("Error executing: got %s but expected %s", jsRes, "test"))
	}

	jsRes, err = vm.RunString("process.argv[2]")
	if err != nil {
		t.Fatal(fmt.Sprintf("Error executing: %s", err))
	}
}

func TestProcessEnvStructure(t *testing.T) {
	vm := goja.New()

	new(require.Registry).Enable(vm)
	Enable(vm)

	if c := vm.Get("process"); c == nil {
		t.Fatal("process not found")
	}

	if c, err := vm.RunString("process.env"); c == nil || err != nil {
		t.Fatal("error accessing process.env")
	}
}

func TestProcessEnvValuesArtificial(t *testing.T) {
	os.Setenv("GOJA_IS_AWESOME", "true")
	defer os.Unsetenv("GOJA_IS_AWESOME")

	vm := goja.New()

	new(require.Registry).Enable(vm)
	Enable(vm)

	jsRes, err := vm.RunString("process.env['GOJA_IS_AWESOME']")

	if err != nil {
		t.Fatal(fmt.Sprintf("Error executing: %s", err))
	}

	if jsRes.String() != "true" {
		t.Fatal(fmt.Sprintf("Error executing: got %s but expected %s", jsRes, "true"))
	}
}

func TestProcessEnvValuesBrackets(t *testing.T) {
	vm := goja.New()

	new(require.Registry).Enable(vm)
	Enable(vm)

	for _, e := range os.Environ() {
		envKeyValue := strings.SplitN(e, "=", 2)
		jsExpr := fmt.Sprintf("process.env['%s']", envKeyValue[0])

		jsRes, err := vm.RunString(jsExpr)

		if err != nil {
			t.Fatal(fmt.Sprintf("Error executing %s: %s", jsExpr, err))
		}

		if jsRes.String() != envKeyValue[1] {
			t.Fatal(fmt.Sprintf("Error executing %s: got %s but expected %s", jsExpr, jsRes, envKeyValue[1]))
		}
	}
}
