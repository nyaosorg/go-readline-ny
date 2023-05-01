package readline

import (
	"context"
	"fmt"

	"github.com/nyaosorg/go-readline-ny/keys"
)

// Deprecated: BindKeyFunc binds function to key
func (km *KeyMap) BindKeyFunc(key string, f Command) error {
	key = keys.NormalizeName(key)
	if code, ok := keys.NameToCode[key]; ok {
		km.BindKey(code, f)
		return nil
	}
	return fmt.Errorf("%s: no such keyname", key)
}

// Deprecated: BindKeyClosure binds closure to key by name
func (km *KeyMap) BindKeyClosure(name string, f func(context.Context, *Buffer) Result) error {
	return km.BindKeyFunc(name, AnonymousCommand(f))
}

// Deprecated: GetBindKey returns the function assigned to given key
func (km *KeyMap) GetBindKey(key string) Command {
	key = keys.NormalizeName(key)
	if ch, ok := keys.NameToCode[key]; ok {
		if km.KeyMap != nil {
			if f, ok := km.KeyMap[ch]; ok {
				return f
			}
		}
	}
	return nil
}

// Deprecated: GetFunc returns Command-object by name
func GetFunc(name string) (Command, error) {
	f, ok := name2func[keys.NormalizeName(name)]
	if ok {
		return f, nil
	}
	return nil, fmt.Errorf("%s: not found in the function-list", name)
}

// Deprecated: BindKeySymbol assigns function to key by names.
func (km *KeyMap) BindKeySymbol(key, funcName string) error {
	f, err := GetFunc(key)
	if err != nil {
		return err
	}
	return km.BindKeyFunc(key, f)
}

// Deprecated: GetBindKey returns the function assigned to given key
func (editor *Editor) GetBindKey(key string) Command {
	key = keys.NormalizeName(key)
	if code, ok := keys.NameToCode[key]; ok {
		return editor.loolupCommand(code.String())
	}
	return nil
}

// Deprecated: use GoCommand instead
type KeyGoFuncT = GoCommand
