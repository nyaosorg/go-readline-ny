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
	KeyMap: map[keys.Code]Command{
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
	KeyMap map[keys.Code]Command
}

func (km *KeyMap) BindKey(key keys.Code, f Command) {
	if km.KeyMap == nil {
		km.KeyMap = map[keys.Code]Command{}
	}
	km.KeyMap[key] = f
}

type AnonymousCommand func(context.Context, *Buffer) Result

func (f AnonymousCommand) String() string {
	return "anonymous"
}

func (f AnonymousCommand) Call(ctx context.Context, B *Buffer) Result {
	return f(ctx, B)
}
