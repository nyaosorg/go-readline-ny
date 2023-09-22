package nameutils

import (
	"fmt"

	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/keys"
)

func BindKeyFunc(km *readline.KeyMap, key string, f readline.Command) error {
	key = keys.NormalizeName(key)
	if code, ok := keys.NameToCode[key]; ok {
		km.BindKey(code, f)
		return nil
	}
	return fmt.Errorf("%s: no such keyname", key)
}

func GetFunc(name string) (readline.Command, error) {
	f, ok := readline.NameToFunc[keys.NormalizeName(name)]
	if ok {
		return f, nil
	}
	return nil, fmt.Errorf("%s: not found in the function-list", name)
}

func BindKeySymbol(km *readline.KeyMap, key, funcName string) error {
	f, err := GetFunc(key)
	if err != nil {
		return err
	}
	return BindKeyFunc(km, key, f)
}
