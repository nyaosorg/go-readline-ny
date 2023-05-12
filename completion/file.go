package completion

import (
	"os"
	"path/filepath"
	"strings"
)

type File struct{}

func (File) Enclosures() string {
	return `"'`
}

func (File) Delimiters() string {
	return "&|><;"
}

func (File) List(list []*Field) (fullnames []string, basenames []string) {
	if len(list) <= 0 {
		return
	}
	target := list[len(list)-1].Str
	var dir, base string
	if tail := target[len(target)-1]; tail == os.PathSeparator || tail == '/' {
		dir = target[:len(target)-1]
		base = ""
	} else {
		dir = filepath.Dir(target)
		base = filepath.Base(target)
	}
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, d := range dirEntries {
		name := d.Name()
		if name == "." || name == ".." {
			continue
		}
		if len(name) < len(base) {
			continue
		}
		if !strings.EqualFold(base, name[:len(base)]) {
			continue
		}
		full := filepath.Join(dir, name)
		if d.IsDir() {
			name += string(os.PathSeparator)
			full += string(os.PathSeparator)
		}
		fullnames = append(fullnames, full)
		basenames = append(basenames, name)
	}
	return
}
