package completion

import (
	"os"
	"path/filepath"
	"strings"
)

func zeroToOneDot(d string) string {
	if d == "" {
		return "."
	}
	return d
}

func PathComplete(field []string) (completionSet []string, listingSet []string) {
	if len(field) <= 0 {
		return
	}
	target := field[len(field)-1]
	var dir, base string
	if tail := target[len(target)-1]; tail == os.PathSeparator || tail == '/' {
		base = ""
		dir = target
	} else {
		base = filepath.Base(target)
		// Do not use filepath.Dir because it removes "./"
		dir = target[:len(target)-len(base)]
	}
	var sep string
	if strings.ContainsRune(target, '/') {
		sep = "/"
	} else {
		sep = string(os.PathSeparator)
	}

	for {
		dirEntries, err := os.ReadDir(zeroToOneDot(dir))
		if err != nil {
			return
		}
		for _, d := range dirEntries {
			name := d.Name()

			// The case where there is only one candidate directory,
			// the completion is completed and the space will be appeneded.
			// Then, It interferes with the subsequent file name input.
			//
			// To avoid that, we narrow down in advance here,
			// and when it is found  that there is only one directory,
			// we add the files under it to the candidates
			// so that completion does not end

			if name == "." || name == ".." {
				continue
			}
			if len(name) < len(base) {
				continue
			}
			if !strings.EqualFold(base, name[:len(base)]) {
				continue
			}
			// Do not use filepath.Join because it removes "./"
			full := dir + name
			if (d.Type() & (os.ModeDir | os.ModeSymlink)) != 0 {
				name += sep
				full += sep
			}
			completionSet = append(completionSet, full)
			listingSet = append(listingSet, name)
		}
		if len(completionSet) != 1 || !strings.HasSuffix(completionSet[0], sep) {
			return
		}
		base = ""
		dir = completionSet[0]
	}
}
