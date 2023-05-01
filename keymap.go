package readline

import (
	"context"
	"fmt"
	"strings"

	"github.com/nyaosorg/go-readline-ny/keys"
)

// Command is the interface for object bound to key-mapping
type Command interface {
	String() string
	Call(ctx context.Context, buffer *Buffer) Result
}

var name2code = map[string]keys.Code{
	"BACKSPACE":   keys.Backspace,
	"CLEAR":       keys.Clear,
	"CTRL":        keys.Ctrl,
	"CTRL_BREAK":  keys.CtrlBreak,
	"CTRL_DOWN":   keys.Down,
	"C_A":         keys.CtrlA,
	"C_B":         keys.CtrlB,
	"C_C":         keys.CtrlC,
	"C_D":         keys.CtrlD,
	"C_E":         keys.CtrlE,
	"C_F":         keys.CtrlF,
	"C_G":         keys.CtrlG,
	"C_H":         keys.CtrlH,
	"C_I":         keys.CtrlI,
	"C_J":         keys.CtrlJ,
	"C_K":         keys.CtrlK,
	"C_L":         keys.CtrlL,
	"C_LEFT":      keys.CtrlLeft,
	"C_M":         keys.CtrlM,
	"C_N":         keys.CtrlN,
	"C_O":         keys.CtrlO,
	"C_P":         keys.CtrlP,
	"C_Q":         keys.CtrlQ,
	"C_R":         keys.CtrlR,
	"C_RIGHT":     keys.CtrlRight,
	"C_S":         keys.CtrlS,
	"C_T":         keys.CtrlT,
	"C_U":         keys.CtrlU,
	"C_UNDERBAR":  keys.CtrlUnderbar,
	"C_UP":        keys.CtrlUp,
	"C_V":         keys.CtrlV,
	"C_W":         keys.CtrlW,
	"C_X":         keys.CtrlX,
	"C_Y":         keys.CtrlY,
	"C_Z":         keys.CtrlZ,
	"C_[":         keys.CtrlLBracket,
	"C_\\":        keys.CtrlBackslash,
	"C_]":         keys.CtrlRBracket,
	"C_^":         keys.CtrlCaret,
	"DELETE":      keys.Delete,
	"DOWN":        keys.Down,
	"END":         keys.End,
	"ENTER":       keys.Enter,
	"ESCAPE":      keys.Escape,
	"F1":          keys.F1,
	"F10":         keys.F10,
	"F11":         keys.F11,
	"F12":         keys.F12,
	"F13":         keys.F13,
	"F14":         keys.F14,
	"F15":         keys.F15,
	"F16":         keys.F16,
	"F17":         keys.F17,
	"F18":         keys.F18,
	"F19":         keys.F19,
	"F2":          keys.F2,
	"F20":         keys.F20,
	"F21":         keys.F21,
	"F22":         keys.F22,
	"F23":         keys.F23,
	"F24":         keys.F24,
	"F3":          keys.F3,
	"F4":          keys.F4,
	"F5":          keys.F5,
	"F6":          keys.F6,
	"F7":          keys.F7,
	"F8":          keys.F8,
	"F9":          keys.F9,
	"HOME":        keys.Home,
	"LEFT":        keys.Left,
	"M_A":         keys.AltA,
	"M_B":         keys.AltB,
	"M_BACKSPACE": keys.Backspace,
	"M_C":         keys.AltC,
	"M_D":         keys.AltD,
	"M_E":         keys.AltE,
	"M_F":         keys.AltF,
	"M_G":         keys.AltG,
	"M_H":         keys.AltH,
	"M_I":         keys.AltI,
	"M_J":         keys.AltJ,
	"M_K":         keys.AltK,
	"M_L":         keys.AltL,
	"M_M":         keys.AltM,
	"M_N":         keys.AltN,
	"M_O":         keys.AltO,
	"M_OEM_2":     keys.AltOEM2,
	"M_P":         keys.AltP,
	"M_Q":         keys.AltQ,
	"M_R":         keys.AltR,
	"M_S":         keys.AltS,
	"M_T":         keys.AltT,
	"M_U":         keys.AltU,
	"M_V":         keys.AltV,
	"M_W":         keys.AltW,
	"M_X":         keys.AltX,
	"M_Y":         keys.AltY,
	"M_Z":         keys.AltZ,
	"PAGEDOWN":    keys.PageDown,
	"PAGEUP":      keys.PageUp,
	"PAUSE":       keys.Pause,
	"RIGHT":       keys.Right,
	"UP":          keys.Up,
}

// KeyMap is the class for key-bindings
type KeyMap map[keys.Code]Command

// GlobalKeyMap is the global keymap for users' customizing
var GlobalKeyMap = KeyMap{
	keys.AltB:         CmdBackwardWord,
	keys.AltF:         CmdForwardWord,
	keys.AltV:         CmdYank,
	keys.AltY:         CmdYankWithQuote,
	keys.Backspace:    CmdBackwardDeleteChar,
	keys.CtrlA:        CmdBeginningOfLine,
	keys.CtrlB:        CmdBackwardChar,
	keys.CtrlC:        CmdInterrupt,
	keys.CtrlD:        CmdDeleteOrAbort,
	keys.CtrlE:        CmdEndOfLine,
	keys.CtrlF:        CmdForwardChar,
	keys.CtrlH:        CmdBackwardDeleteChar,
	keys.CtrlK:        CmdKillLine,
	keys.CtrlL:        CmdClearScreen,
	keys.CtrlLeft:     CmdBackwardWord,
	keys.CtrlM:        CmdAcceptLine,
	keys.CtrlN:        CmdNextHistory,
	keys.CtrlP:        CmdPreviousHistory,
	keys.CtrlQ:        CmdQuotedInsert,
	keys.CtrlR:        CmdISearchBackward,
	keys.CtrlRight:    CmdForwardWord,
	keys.CtrlT:        CmdSwapChar,
	keys.CtrlU:        CmdUnixLineDiscard,
	keys.CtrlUnderbar: CmdUndo,
	keys.CtrlV:        CmdQuotedInsert,
	keys.CtrlW:        CmdUnixWordRubout,
	keys.CtrlY:        CmdYank,
	keys.CtrlZ:        CmdUndo,
	keys.Delete:       CmdDeleteChar,
	keys.Down:         CmdNextHistory,
	keys.End:          CmdEndOfLine,
	keys.Escape:       CmdKillWholeLine,
	keys.Home:         CmdBeginningOfLine,
	keys.Left:         CmdBackwardChar,
	keys.Right:        CmdForwardChar,
	keys.Up:           CmdPreviousHistory,
}

func normWord(src string) string {
	return strings.Replace(strings.ToUpper(src), "-", "_", -1)
}

func (km KeyMap) BindKey(key keys.Code, f Command) {
	km[key] = f
}

// BindKeyFunc binds function to key
func (km KeyMap) BindKeyFunc(key string, f Command) error {
	key = normWord(key)
	if code, ok := name2code[key]; ok {
		km.BindKey(code, f)
		return nil
	}
	return fmt.Errorf("%s: no such keyname", key)
}

// BindKeyClosure binds closure to key by name
func (km KeyMap) BindKeyClosure(name string, f func(context.Context, *Buffer) Result) error {
	return km.BindKeyFunc(name, &GoCommand{Func: f, Name: "annonymous"})
}

// GetBindKey returns the function assigned to given key
func (km KeyMap) GetBindKey(key string) Command {
	key = normWord(key)
	if ch, ok := name2code[key]; ok {
		if f, ok := km[ch]; ok {
			return f
		}
	}
	return nil
}

// GetFunc returns Command-object by name
func GetFunc(name string) (Command, error) {
	f, ok := name2func[normWord(name)]
	if ok {
		return f, nil
	}
	return nil, fmt.Errorf("%s: not found in the function-list", name)
}

// BindKeySymbol assigns function to key by names.
func (km KeyMap) BindKeySymbol(key, funcName string) error {
	f, err := GetFunc(key)
	if err != nil {
		return err
	}
	return km.BindKeyFunc(key, f)
}
