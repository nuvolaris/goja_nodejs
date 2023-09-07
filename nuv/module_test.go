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

	t.Run("readFile", func(t *testing.T) {
		content, err := vm.RunString("nuv.readFile('testdata/sample.txt')")
		if err != nil {
			t.Fatal("nuv.readFile error", err)
		}
		if content.Export().(string) != "a sample text file" {
			t.Fatal("wrong nuv.readFile output, want 'a sample text file', got", content)
		}
	})

	t.Run("writeFile", func(t *testing.T) {
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
	})

	t.Run("toYaml", func(t *testing.T) {
		result, err := vm.RunString("nuv.toYaml({ version: 3 })")
		if err != nil {
			t.Fatal("nuv.toYaml() error", err)
		}
		if result.Export().(string) != "version: 3\n" {
			t.Fatal("wrong nuv.toYaml() output, want 'version: 3\n', got", result)
		}
	})

	t.Run("fromYaml", func(t *testing.T) {
		objRes, err := vm.RunString("nuv.fromYaml('version: 3')")
		if err != nil {
			t.Fatal("nuv.fromYaml() error", err)
		}

		if objRes.Export().(map[string]interface{})["version"].(int) != 3 {
			t.Fatal("wrong nuv.fromYaml() output, want '3', got", objRes)
		}
	})

	t.Run("scan", func(t *testing.T) {
		scanRes, err := vm.RunString("nuv.scan('testdata', (folder) => folder + ' ')")
		if err != nil {
			t.Fatal("nuv.scan() error", err)
		}

		if scanRes.Export().(string) != "testdata testdata/subfolder " {
			t.Fatal("wrong nuv.scan() output, want 'testdata testdata/subfolder ', got", scanRes)
		}
	})

	t.Run("readDir", func(t *testing.T) {
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
	})

	t.Run("basePath", func(t *testing.T) {
		basePatRes, err := vm.RunString("nuv.basePath('testdata/test/sample')")
		if err != nil {
			t.Fatal("nuv.basePath() error", err)
		}

		if basePatRes.Export().(string) != "sample" {
			t.Fatal("wrong nuv.basePath() output, want 'sample', got", basePatRes)
		}
	})

	t.Run("fileExt", func(t *testing.T) {
		fileExtRes, err := vm.RunString("nuv.fileExt('testdata/test/sample.txt')")
		if err != nil {
			t.Fatal("nuv.fileExt() error", err)
		}

		if fileExtRes.Export().(string) != ".txt" {
			t.Fatal("wrong nuv.fileExt() output, want '.txt', got", fileExtRes)
		}
	})

	t.Run("isDir", func(t *testing.T) {
		isDirRes, err := vm.RunString("nuv.isDir('testdata/test/sample.txt')")
		if err != nil {
			t.Fatal("nuv.isDir() error", err)
		}

		if isDirRes.Export().(bool) != false {
			t.Fatal("wrong nuv.isDir() output, want 'false', got", isDirRes)
		}

		isDirRes, err = vm.RunString("nuv.isDir('testdata')")
		if err != nil {
			t.Fatal("nuv.isDir() error", err)
		}

		if isDirRes.Export().(bool) != true {
			t.Fatal("wrong nuv.isDir() output, want 'true', got", isDirRes)
		}
	})

	t.Run("joinPath", func(t *testing.T) {
		joinPathRes, err := vm.RunString("nuv.joinPath('testdata', 'test/sample.txt')")
		if err != nil {
			t.Fatal("nuv.joinPath() error", err)
		}

		if joinPathRes.Export().(string) != "testdata/test/sample.txt" {
			t.Fatal("wrong nuv.joinPath() output, want 'testdata/test/sample.txt', got", joinPathRes)
		}
	})

	t.Run("exists", func(t *testing.T) {
		existsRes, err := vm.RunString("nuv.exists('testdata/sample.txt')")
		if err != nil {
			t.Fatal("nuv.exists() error", err)
		}

		if existsRes.Export().(bool) != true {
			t.Fatal("wrong nuv.exists() output, want 'true', got", existsRes)
		}

		existsRes, err = vm.RunString("nuv.exists('testdata/sample2.txt')")
		if err != nil {
			t.Fatal("nuv.exists() error", err)
		}

		if existsRes.Export().(bool) != false {
			t.Fatal("wrong nuv.exists() output, want 'false', got", existsRes)
		}
	})

	t.Run("nuvExec", func(t *testing.T) {
		oldNuvVersion := os.Getenv("NUV_VERSION")
		os.Setenv("NUV_VERSION", "3")
		_, err := vm.RunString("nuv.nuvExec('nuv -v')")
		os.Setenv("NUV_VERSION", oldNuvVersion)
		if err != nil {
			t.Fatal("nuv.nuvExec() error", err)
		}
	})
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
