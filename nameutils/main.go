package nameutils

import (
	"context"
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

func BindKeyClosure(km *readline.KeyMap, name string, f func(context.Context, *readline.Buffer) readline.Result) error {
	return BindKeyFunc(km, name, readline.AnonymousCommand(f))
}

func GetBindKeyMap(km *readline.KeyMap, key string) readline.Command {
	key = keys.NormalizeName(key)
	if ch, ok := keys.NameToCode[key]; ok {
		f, _ := km.Lookup(ch)
		return f
	}
	return nil
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

func GetBindKeyEditor(editor *readline.Editor, key string) readline.Command {
	key = keys.NormalizeName(key)
	if code, ok := keys.NameToCode[key]; ok {
		return editor.LookupCommand(code.String())
	}
	return nil
}
