package readline

import (
	"github.com/nyaosorg/go-readline-ny/internal/ttysub"
)

type XTty = ttysub.XTty

// Deprecated: it has a problem where only the first line is taken when multiple lines are pasted from the clipboard, etc
func GetKey(tty XTty) (string, error) {
	clean, err := tty.Raw()
	if err != nil {
		return "", err
	}
	defer clean()

	return ttysub.GetOneKey(tty)
}
