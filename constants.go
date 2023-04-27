package readline

import (
	"github.com/nyaosorg/go-readline-ny/keys"
)

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
