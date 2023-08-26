package nuv

import (
	"os"
	"testing"

	"github.com/nuvolaris/goja"
	"github.com/nuvolaris/goja_nodejs/require"
)

func TestNuv(t *testing.T) {
	vm := goja.New()

	new(require.Registry).Enable(vm)
	Enable(vm)

	if n := vm.Get("nuv"); n == nil {
		t.Fatal("nuv not found")
	}

	content, err := vm.RunString("nuv.readFile('testdata/sample.txt')")
	if err != nil {
		t.Fatal("nuv.readFile error", err)
	}
	if content.Export().(string) != "a sample text file" {
		t.Fatal("wrong nuv.readFile output, want 'a sample text file', got", content)
	}

	if _, err := vm.RunString("nuv.writeFile('testdata/written.txt', 'sample from js')"); err != nil {
		t.Fatal("nuv.writeFile error", err)
	}

	// check that testdata/written.txt exists
	writtenContent, err := vm.RunString("nuv.readFile('testdata/written.txt')")
	if err != nil {
		t.Fatal("nuv.readFile error after writeFile", err)
	}
	if writtenContent.Export().(string) != "sample from js" {
		t.Fatal("wrong nuv.readFile output after writeFile, want 'sample from js', got", writtenContent)
	}
	// remove testdata/written.txt with go
	_ = os.Remove("testdata/written.txt")

	result, err := vm.RunString("nuv.toYaml({ version: 3 })")
	if err != nil {
		t.Fatal("nuv.toYaml() error", err)
	}
	if result.Export().(string) != "version: 3\n" {
		t.Fatal("wrong nuv.toYaml() output, want 'version: 3\n', got", result)
	}

	objRes, err := vm.RunString("nuv.fromYaml('version: 3')")
	if err != nil {
		t.Fatal("nuv.fromYaml() error", err)
	}

	if objRes.Export().(map[string]interface{})["version"].(int) != 3 {
		t.Fatal("wrong nuv.fromYaml() output, want '3', got", objRes)
	}

	scanRes, err := vm.RunString("nuv.scan('testdata', (folder) => folder + ' ')")
	if err != nil {
		t.Fatal("nuv.scan() error", err)
	}

	if scanRes.Export().(string) != "testdata testdata/subfolder " {
		t.Fatal("wrong nuv.scan() output, want 'testdata testdata/subfolder ', got", scanRes)
	}
}
