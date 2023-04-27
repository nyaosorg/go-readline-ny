package readline

import (
	"context"

	"github.com/nyaosorg/go-readline-ny/keys"
)

const (
	K_BACKSPACE      = "BACKSPACE"
	K_CAPSLOCK       = "CAPSLOCK"
	K_CLEAR          = "CLEAR"
	K_CTRL           = "CTRL"
	K_CTRL_A         = "C_A"
	K_CTRL_B         = "C_B"
	K_CTRL_BREAK     = "C_BREAK"
	K_CTRL_C         = "C_C"
	K_CTRL_D         = "C_D"
	K_CTRL_E         = "C_E"
	K_CTRL_F         = "C_F"
	K_CTRL_G         = "C_G"
	K_CTRL_H         = "C_H"
	K_CTRL_I         = "C_I"
	K_CTRL_J         = "C_J"
	K_CTRL_K         = "C_K"
	K_CTRL_L         = "C_L"
	K_CTRL_M         = "C_M"
	K_CTRL_N         = "C_N"
	K_CTRL_O         = "C_O"
	K_CTRL_P         = "C_P"
	K_CTRL_Q         = "C_Q"
	K_CTRL_R         = "C_R"
	K_CTRL_S         = "C_S"
	K_CTRL_T         = "C_T"
	K_CTRL_U         = "C_U"
	K_CTRL_V         = "C_V"
	K_CTRL_W         = "C_W"
	K_CTRL_X         = "C_X"
	K_CTRL_Y         = "C_Y"
	K_CTRL_Z         = "C_Z"
	K_CTRL_UNDERBAR  = "C_UNDERBAR"
	K_CTRL_LBRACKET  = "C_["
	K_CTRL_RBRACKET  = "C_]"
	K_CTRL_BACKSLASH = "C_\\"
	K_CTRL_CARET     = "C_^"
	K_DELETE         = "DEL"
	K_DOWN           = "DOWN"
	K_CTRL_DOWN      = "C_DOWN"
	K_END            = "END"
	K_ENTER          = "ENTER"
	K_ESCAPE         = "ESCAPE"
	K_F1             = "F1"
	K_F10            = "F10"
	K_F11            = "F11"
	K_F12            = "F12"
	K_F13            = "F13"
	K_F14            = "F14"
	K_F15            = "F15"
	K_F16            = "F16"
	K_F17            = "F17"
	K_F18            = "F18"
	K_F19            = "F19"
	K_F2             = "F2"
	K_F20            = "F20"
	K_F21            = "F21"
	K_F22            = "F22"
	K_F23            = "F23"
	K_F24            = "F24"
	K_F3             = "F3"
	K_F4             = "F4"
	K_F5             = "F5"
	K_F6             = "F6"
	K_F7             = "F7"
	K_F8             = "F8"
	K_F9             = "F9"
	K_HOME           = "HOME"
	K_LEFT           = "LEFT"
	K_CTRL_LEFT      = "C_LEFT"
	K_PAGEDOWN       = "PAGEDOWN"
	K_PAGEUP         = "PAGEUP"
	K_PAUSE          = "PAUSE"
	K_RIGHT          = "RIGHT"
	K_CTRL_RIGHT     = "C_RIGHT"
	K_SHIFT          = "SHIFT"
	K_UP             = "UP"
	K_CTRL_UP        = "C_UP"
	K_ALT_A          = "M_A"
	K_ALT_B          = "M_B"
	K_ALT_BACKSPACE  = "M_BACKSPACE"
	K_ALT_BREAK      = "M_BREAK"
	K_ALT_C          = "M_C"
	K_ALT_D          = "M_D"
	K_ALT_E          = "M_E"
	K_ALT_F          = "M_F"
	K_ALT_G          = "M_G"
	K_ALT_H          = "M_H"
	K_ALT_I          = "M_I"
	K_ALT_J          = "M_J"
	K_ALT_K          = "M_K"
	K_ALT_L          = "M_L"
	K_ALT_M          = "M_M"
	K_ALT_N          = "M_N"
	K_ALT_O          = "M_O"
	K_ALT_P          = "M_P"
	K_ALT_Q          = "M_Q"
	K_ALT_R          = "M_R"
	K_ALT_S          = "M_S"
	K_ALT_T          = "M_T"
	K_ALT_U          = "M_U"
	K_ALT_V          = "M_V"
	K_ALT_W          = "M_W"
	K_ALT_X          = "M_X"
	K_ALT_Y          = "M_Y"
	K_ALT_Z          = "M_Z"
	K_ALT_OEM_2      = "M_OEM_2"
)

const (
	F_ACCEPT_LINE          = "ACCEPT_LINE"
	F_BACKWARD_CHAR        = "BACKWARD_CHAR"
	F_BACKWARD_WORD        = "BACKWARD_WORD"
	F_BACKWARD_DELETE_CHAR = "BACKWARD_DELETE_CHAR"
	F_BEGINNING_OF_LINE    = "BEGINNING_OF_LINE"
	F_CLEAR_SCREEN         = "CLEAR_SCREEN"
	F_DELETE_CHAR          = "DELETE_CHAR"
	F_DELETE_OR_ABORT      = "DELETE_OR_ABORT"
	F_END_OF_LINE          = "END_OF_LINE"
	F_FORWARD_CHAR         = "FORWARD_CHAR"
	F_FORWARD_WORD         = "FORWARD_WORD"
	F_HISTORY_DOWN         = "HISTORY_DOWN" // for compatible
	F_HISTORY_UP           = "HISTORY_UP"   // for compatible
	F_NEXT_HISTORY         = "NEXT_HISTORY"
	F_PREVIOUS_HISTORY     = "PREVIOUS_HISTORY"
	F_INTR                 = "INTR"
	F_ISEARCH_BACKWARD     = "ISEARCH_BACKWARD"
	F_KILL_LINE            = "KILL_LINE"
	F_KILL_WHOLE_LINE      = "KILL_WHOLE_LINE"
	F_PASS                 = "PASS"
	F_QUOTED_INSERT        = "QUOTED_INSERT"
	F_REPAINT_ON_NEWLINE   = "REPAINT_ON_NEWLINE"
	F_SWAPCHAR             = "SWAPCHAR"
	F_UNIX_LINE_DISCARD    = "UNIX_LINE_DISCARD"
	F_UNIX_WORD_RUBOUT     = "UNIX_WORD_RUBOUT"
	F_YANK                 = "YANK"
	F_YANK_WITH_QUOTE      = "YANK_WITH_QUOTE"
	F_UNDO                 = "UNDO"
)

var name2code = map[string]keys.Code{
	K_BACKSPACE:      keys.Backspace,
	K_CTRL_A:         keys.CtrlA,
	K_CTRL_B:         keys.CtrlB,
	K_CTRL_C:         keys.CtrlC,
	K_CTRL_D:         keys.CtrlD,
	K_CTRL_E:         keys.CtrlE,
	K_CTRL_F:         keys.CtrlF,
	K_CTRL_G:         keys.CtrlG,
	K_CTRL_H:         keys.CtrlH,
	K_CTRL_I:         keys.CtrlI,
	K_CTRL_J:         keys.CtrlJ,
	K_CTRL_K:         keys.CtrlK,
	K_CTRL_L:         keys.CtrlL,
	K_CTRL_M:         keys.CtrlM,
	K_CTRL_N:         keys.CtrlN,
	K_CTRL_O:         keys.CtrlO,
	K_CTRL_P:         keys.CtrlP,
	K_CTRL_Q:         keys.CtrlQ,
	K_CTRL_R:         keys.CtrlR,
	K_CTRL_S:         keys.CtrlS,
	K_CTRL_T:         keys.CtrlT,
	K_CTRL_U:         keys.CtrlU,
	K_CTRL_V:         keys.CtrlV,
	K_CTRL_W:         keys.CtrlW,
	K_CTRL_X:         keys.CtrlX,
	K_CTRL_Y:         keys.CtrlY,
	K_CTRL_Z:         keys.CtrlZ,
	K_CTRL_LBRACKET:  keys.CtrlLBracket, // C-[
	K_CTRL_BACKSLASH: keys.CtrlBackslash,
	K_CTRL_RBRACKET:  keys.CtrlRBracket, // C-]
	K_CTRL_CARET:     keys.CtrlCaret,    // C-^
	K_CTRL_UNDERBAR:  keys.CtrlUnderbar,
	// K_DELETE:        "\x7F",
	K_ENTER:         keys.Enter,
	K_ESCAPE:        keys.Escape,
	K_ALT_A:         keys.AltA,
	K_ALT_B:         keys.AltB,
	K_ALT_BACKSPACE: keys.Backspace,
	K_ALT_C:         keys.AltC,
	K_ALT_D:         keys.AltD,
	K_ALT_E:         keys.AltE,
	K_ALT_F:         keys.AltF,
	K_ALT_G:         keys.AltG,
	K_ALT_H:         keys.AltH,
	K_ALT_I:         keys.AltI,
	K_ALT_J:         keys.AltJ,
	K_ALT_K:         keys.AltK,
	K_ALT_L:         keys.AltL,
	K_ALT_M:         keys.AltM,
	K_ALT_N:         keys.AltN,
	K_ALT_O:         keys.AltO,
	K_ALT_P:         keys.AltP,
	K_ALT_Q:         keys.AltQ,
	K_ALT_R:         keys.AltR,
	K_ALT_S:         keys.AltS,
	K_ALT_T:         keys.AltT,
	K_ALT_U:         keys.AltU,
	K_ALT_V:         keys.AltV,
	K_ALT_W:         keys.AltW,
	K_ALT_X:         keys.AltX,
	K_ALT_Y:         keys.AltY,
	K_ALT_Z:         keys.AltZ,
	K_ALT_OEM_2:     keys.AltOEM2,
	K_CLEAR:         keys.Clear,
	K_CTRL:          keys.Ctrl,
	K_CTRL_BREAK:    keys.CtrlBreak,
	K_DELETE:        keys.Delete,
	K_DOWN:          keys.Down,
	K_CTRL_DOWN:     keys.Down,
	K_END:           keys.End,
	K_F10:           keys.F10,
	K_F11:           keys.F11,
	K_F12:           keys.F12,
	K_F13:           keys.F13,
	K_F14:           keys.F14,
	K_F15:           keys.F15,
	K_F16:           keys.F16,
	K_F17:           keys.F17,
	K_F18:           keys.F18,
	K_F19:           keys.F19,
	K_F1:            keys.F1,
	K_F20:           keys.F20,
	K_F21:           keys.F21,
	K_F22:           keys.F22,
	K_F23:           keys.F23,
	K_F24:           keys.F24,
	K_F2:            keys.F2,
	K_F3:            keys.F3,
	K_F4:            keys.F4,
	K_F5:            keys.F5,
	K_F6:            keys.F6,
	K_F7:            keys.F7,
	K_F8:            keys.F8,
	K_F9:            keys.F9,
	K_HOME:          keys.Home,
	K_LEFT:          keys.Left,
	K_CTRL_LEFT:     keys.CtrlLeft,
	K_PAGEDOWN:      keys.PageDown,
	K_PAGEUP:        keys.PageUp,
	K_PAUSE:         keys.Pause,
	K_RIGHT:         keys.Right,
	K_CTRL_RIGHT:    keys.CtrlRight,
	K_UP:            keys.Up,
	K_CTRL_UP:       keys.CtrlUp,
}

// KeyCode from
// http://msdn.microsoft.com/ja-jp/library/windows/desktop/dd375731(v=vs.85).aspx

var _name2func = map[string]func(context.Context, *Buffer) Result{
	F_ACCEPT_LINE:          keyFuncEnter,
	F_BACKWARD_CHAR:        keyFuncBackward,
	F_BACKWARD_WORD:        keyFuncBackwardWord,
	F_BACKWARD_DELETE_CHAR: keyFuncBackSpace,
	F_BEGINNING_OF_LINE:    keyFuncHead,
	F_CLEAR_SCREEN:         keyFuncCLS,
	F_DELETE_CHAR:          keyFuncDelete,
	F_DELETE_OR_ABORT:      keyFuncDeleteOrAbort,
	F_END_OF_LINE:          keyFuncTail,
	F_FORWARD_CHAR:         keyFuncForward,
	F_FORWARD_WORD:         keyFuncForwardWord,
	F_HISTORY_DOWN:         keyFuncHistoryDown, // for compatible
	F_HISTORY_UP:           keyFuncHistoryUp,   // for compatible
	F_NEXT_HISTORY:         keyFuncHistoryDown,
	F_PREVIOUS_HISTORY:     keyFuncHistoryUp,
	F_INTR:                 keyFuncIntr,
	F_ISEARCH_BACKWARD:     keyFuncIncSearch,
	F_KILL_LINE:            keyFuncClearAfter,
	F_KILL_WHOLE_LINE:      keyFuncClear,
	F_PASS:                 nil,
	F_QUOTED_INSERT:        keyFuncQuotedInsert,
	F_UNIX_LINE_DISCARD:    keyFuncClearBefore,
	F_UNIX_WORD_RUBOUT:     keyFuncWordRubout,
	F_YANK:                 keyFuncPaste,
	F_YANK_WITH_QUOTE:      keyFuncPasteQuote,
	F_SWAPCHAR:             keyFuncSwapChar,
	F_REPAINT_ON_NEWLINE:   keyFuncRepaintOnNewline,
	F_UNDO:                 keyFuncUndo,
}

func name2func(keyName string) KeyFuncT {
	if p, ok := _name2func[keyName]; ok {
		return &KeyGoFuncT{
			Func: p,
			Name: keyName,
		}
	}
	return nil
}
