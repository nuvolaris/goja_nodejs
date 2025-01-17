package nuv

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	defaultNuvScanner Scanner = &StdScanner{}
)

// StdScanner implements the nuv.Scanner interface
// with the utilities for nuv -scan.
type StdScanner struct {
}

// readDir implements Scanner.
func (*StdScanner) readDir(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var filenames []string
	for _, entry := range entries {
		filenames = append(filenames, entry.Name())
	}

	return filenames, nil
}

// readFile implements Scanner.
func (*StdScanner) readFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// writeFile implements Scanner.
func (*StdScanner) writeFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

// toYaml implements Scanner.
func (*StdScanner) toYaml(data map[string]interface{}) (string, error) {
	// Convert the map to a YAML string
	yamlBytes, err := yaml.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(yamlBytes), nil
}

// fromYaml implements Scanner.
func (*StdScanner) fromYaml(yamlStr string) (map[string]interface{}, error) {
	// Convert the YAML string to a Go map
	var data map[string]interface{}
	err := yaml.Unmarshal([]byte(yamlStr), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// scan implements Scanner.
func (*StdScanner) scan(root string, f func(string) string) string {
	strBuilder := strings.Builder{}
	filepath.WalkDir(
		root,
		func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				strBuilder.WriteString(f(path))
			}
			return nil
		},
	)

	return strBuilder.String()
}

func (*StdScanner) basePath(path string) string {
	return filepath.Base(path)
}

func (*StdScanner) fileExt(path string) string {
	return filepath.Ext(path)
}

func (*StdScanner) isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func (*StdScanner) joinPath(path1 string, path2 string) string {
	return filepath.Join(path1, path2)
}

func (*StdScanner) exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (*StdScanner) nuvExec(cmd string, args ...string) string {
	shCmd := exec.Command(cmd, args...)
	out, _ := shCmd.CombinedOutput()
	return string(out)
}
