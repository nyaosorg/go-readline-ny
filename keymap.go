package readline

import (
	"context"

	"github.com/nyaosorg/go-readline-ny/keys"
)

// Command is the interface for object bound to key-mapping
type Command interface {
	String() string
	Call(ctx context.Context, buffer *Buffer) Result
}

// GlobalKeyMap is the global keymap for users' customizing
var GlobalKeyMap = &KeyMap{
	table: map[keys.Code]Command{
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
	},
}

// KeyMap is the class for key-bindings
type KeyMap struct {
	table map[keys.Code]Command
}

func (km *KeyMap) BindKey(key keys.Code, f Command) {
	if km.table == nil {
		km.table = map[keys.Code]Command{}
	}
	if f == nil {
		delete(km.table, key)
	} else {
		km.table[key] = f
	}
}

func (km *KeyMap) Lookup(key keys.Code) (Command, bool) {
	if km.table == nil {
		return nil, false
	}
	f, ok := km.table[key]
	return f, ok
}

// AnonymousCommand is a type that defines an unnamed command.
type AnonymousCommand func(context.Context, *Buffer) Result

func (f AnonymousCommand) String() string {
	return "anonymous"
}

func (f AnonymousCommand) Call(ctx context.Context, B *Buffer) Result {
	return f(ctx, B)
}
