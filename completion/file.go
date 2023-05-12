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

func (File) List(field []string) (fullnames []string, basenames []string) {
	if len(field) <= 0 {
		return
	}
	target := field[len(field)-1]
	var dir, base string
	if tail := target[len(target)-1]; tail == os.PathSeparator || tail == '/' {
		dir = target[:len(target)-1]
		base = ""
	} else {
		dir = filepath.Dir(target)
		base = filepath.Base(target)
	}
	for {
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
			if (d.Type() & (os.ModeDir | os.ModeSymlink)) != 0 {
				name += string(os.PathSeparator)
				full += string(os.PathSeparator)
			}
			fullnames = append(fullnames, full)
			basenames = append(basenames, name)
		}
		if len(fullnames) != 1 || fullnames[0][len(fullnames[0])-1] != os.PathSeparator {
			return
		}
		dir = fullnames[0]
		base = ""
	}
}
