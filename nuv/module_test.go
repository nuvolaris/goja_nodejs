package nuv

import (
	_ "embed"
	"os"
	"strings"
	"testing"

	"github.com/nuvolaris/goja"
	"github.com/nuvolaris/goja_nodejs/console"
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

	readDirRes, err := vm.RunString("nuv.readDir('testdata')")
	if err != nil {
		t.Fatal("nuv.readDir() error", err)
	}

	want := []string{"nuv_test.js", "sample.txt", "subfolder"}
	got := readDirRes.Export().([]string)

	if len(got) != len(want) {
		t.Fatal("wrong nuv.readDir() output, want", want, "got", got)
	}

	for i, v := range want {
		if v != got[i] {
			t.Fatal("wrong nuv.readDir() output, want", want, "got", got)
		}
	}
}

//go:embed testdata/nuv_test.js
var nuvTest string

func TestNuvWithScanner(t *testing.T) {
	var stdoutStr, stderrStr string

	printer := console.StdPrinter{
		StdoutPrint: func(s string) { stdoutStr += s },
		StderrPrint: func(s string) { stderrStr += s },
	}

	vm := goja.New()

	registry := new(require.Registry)
	registry.Enable(vm)
	registry.RegisterNativeModule("console", console.RequireWithPrinter(printer))
	registry.RegisterNativeModule(ModuleName, RequireWithScanner(&StdScanner{}))

	Enable(vm)
	console.Enable(vm)

	if c := vm.Get("console"); c == nil {
		t.Fatal("console not found")
	}

	if n := vm.Get("nuv"); n == nil {
		t.Fatal("nuv not found")
	}

	_, err := vm.RunScript("testdata/url_test.js", nuvTest)
	if err != nil {
		if ex, ok := err.(*goja.Exception); ok {
			t.Fatal(ex.String())
		}
		t.Fatal("Failed to process url script.", err)
	}
	os.Remove("testdata/sample2.txt")

	if stdoutStr == "" {
		t.Fatal("stdout is empty")
	}

	out := strings.Split(stdoutStr, "***")

	tt := []struct {
		name     string
		expected string
	}{
		{
			name:     "readFile",
			expected: "a sample text file",
		},
		{
			name:     "writeFile",
			expected: "test write",
		},
		{
			name:     "toYaml",
			expected: "a: 1\nb: 2\n",
		},
		{
			name:     "fromYaml",
			expected: "{\"a\":1,\"b\":2}",
		},
		{
			name:     "scan",
			expected: "testdata testdata/subfolder ",
		},
		{
			name:     "readDir",
			expected: "nuv_test.js,sample.txt,sample2.txt,subfolder",
		},
	}

	for i, test := range tt {
		if out[i] != test.expected {
			t.Fatalf("failed to match %s property. got: %s, expected: %s", test.name, out[i], test.expected)
		}
	}
}
