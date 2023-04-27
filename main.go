package readline

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	//tty "github.com/nyaosorg/go-readline-ny/tty10"
	"github.com/nyaosorg/go-readline-ny/keys"
	tty "github.com/nyaosorg/go-readline-ny/tty8"
)

// Result is the type for readline's result.
type Result int

const (
	// CONTINUE is returned by key-functions to continue the line editor
	CONTINUE Result = iota
	// ENTER is returned by key-functions when Enter key is pressed
	ENTER Result = iota
	// ABORT is returned by key-functions when Ctrl-D is pressed with no command-line
	ABORT Result = iota
	// INTR is returned by key-functions when Ctrl-C is pressed
	INTR Result = iota
)

// String makes Result to fmt.Stringer
func (R Result) String() string {
	switch R {
	case CONTINUE:
		return "CONTINUE"
	case ENTER:
		return "ENTER"
	case ABORT:
		return "ABORT"
	case INTR:
		return "INTR"
	default:
		return "ERROR"
	}
}

// KeyFuncT is the interface for object bound to key-mapping
type KeyFuncT interface {
	Call(ctx context.Context, buffer *Buffer) Result
}

// KeyGoFuncT is the implement of KeyFuncT which has a name and a function
type KeyGoFuncT struct {
	Func func(ctx context.Context, buffer *Buffer) Result
	Name string
}

// Call calls the function the receiver contains
func (K *KeyGoFuncT) Call(ctx context.Context, buffer *Buffer) Result {
	if K.Func == nil {
		return CONTINUE
	}
	return K.Func(ctx, buffer)
}

// String returns KeyGoFuncT's name
func (K KeyGoFuncT) String() string {
	return K.Name
}

var defaultKeyMap = map[keys.Code]KeyFuncT{
	keys.AltB:         name2func(F_BACKWARD_WORD),
	keys.AltF:         name2func(F_FORWARD_WORD),
	keys.AltV:         name2func(F_YANK),
	keys.AltY:         name2func(F_YANK_WITH_QUOTE),
	keys.Backspace:    name2func(F_BACKWARD_DELETE_CHAR),
	keys.Ctrl:         name2func(F_PASS),
	keys.CtrlA:        name2func(F_BEGINNING_OF_LINE),
	keys.CtrlB:        name2func(F_BACKWARD_CHAR),
	keys.CtrlC:        name2func(F_INTR),
	keys.CtrlD:        name2func(F_DELETE_OR_ABORT),
	keys.CtrlE:        name2func(F_END_OF_LINE),
	keys.CtrlF:        name2func(F_FORWARD_CHAR),
	keys.CtrlH:        name2func(F_BACKWARD_DELETE_CHAR),
	keys.CtrlK:        name2func(F_KILL_LINE),
	keys.CtrlL:        name2func(F_CLEAR_SCREEN),
	keys.CtrlLeft:     name2func(F_BACKWARD_WORD),
	keys.CtrlM:        name2func(F_ACCEPT_LINE),
	keys.CtrlN:        name2func(F_HISTORY_DOWN),
	keys.CtrlP:        name2func(F_HISTORY_UP),
	keys.CtrlQ:        name2func(F_QUOTED_INSERT),
	keys.CtrlR:        name2func(F_ISEARCH_BACKWARD),
	keys.CtrlRight:    name2func(F_FORWARD_WORD),
	keys.CtrlT:        name2func(F_SWAPCHAR),
	keys.CtrlU:        name2func(F_UNIX_LINE_DISCARD),
	keys.CtrlUnderbar: name2func(F_UNDO),
	keys.CtrlV:        name2func(F_QUOTED_INSERT),
	keys.CtrlW:        name2func(F_UNIX_WORD_RUBOUT),
	keys.CtrlY:        name2func(F_YANK),
	keys.CtrlZ:        name2func(F_UNDO),
	keys.Delete:       name2func(F_DELETE_CHAR),
	keys.Down:         name2func(F_HISTORY_DOWN),
	keys.End:          name2func(F_END_OF_LINE),
	keys.Escape:       name2func(F_KILL_WHOLE_LINE),
	keys.Home:         name2func(F_BEGINNING_OF_LINE),
	keys.Left:         name2func(F_BACKWARD_CHAR),
	keys.Right:        name2func(F_FORWARD_CHAR),
	keys.Up:           name2func(F_HISTORY_UP),
}

func normWord(src string) string {
	return strings.Replace(strings.ToUpper(src), "-", "_", -1)
}

// KeyMap is the class for key-bindings
type KeyMap struct {
	KeyMap map[keys.Code]KeyFuncT
}

// BindKeyFunc binds function to key
func (km *KeyMap) BindKeyFunc(key string, f KeyFuncT) error {
	key = normWord(key)
	if char, ok := name2code[key]; ok {
		if km.KeyMap == nil {
			km.KeyMap = map[keys.Code]KeyFuncT{}
		}
		km.KeyMap[char] = f
		return nil
	}
	return fmt.Errorf("%s: no such keyname", key)
}

// BindKeyClosure binds closure to key by name
func (km *KeyMap) BindKeyClosure(name string, f func(context.Context, *Buffer) Result) error {
	return km.BindKeyFunc(name, &KeyGoFuncT{Func: f, Name: "annonymous"})
}

// GetBindKey returns the function assigned to given key
func (km *KeyMap) GetBindKey(key string) KeyFuncT {
	key = normWord(key)
	if ch, ok := name2code[key]; ok {
		if km.KeyMap != nil {
			if f, ok := km.KeyMap[ch]; ok {
				return f
			}
		}
	}
	return nil
}

// GlobalKeyMap is the global keymap for users' customizing
var GlobalKeyMap KeyMap

type ITty interface {
	Raw() (func() error, error)
	ReadRune() (rune, error)
	Buffered() bool
	Open() error
	Close() error
	Size() (int, int, error)
	GetResizeNotifier() func() (int, int, bool)
}

// Editor is the main class to hold the parameter for ReadLine
type Editor struct {
	KeyMap
	History        IHistory
	Writer         io.Writer
	Out            *bufio.Writer
	Prompt         func() (int, error)
	Default        string
	Cursor         int
	LineFeed       func(Result)
	Tty            ITty
	Coloring       Coloring
	HistoryCycling bool
}

// GetBindKey returns the function assigned to given key
func (editor *Editor) GetBindKey(key string) KeyFuncT {
	key = normWord(key)
	if code, ok := name2code[key]; ok {
		return editor.getKeyFunction(code.String())
	}
	return nil
}

// GetFunc returns KeyFuncT-object by name
func GetFunc(funcName string) (KeyFuncT, error) {
	rc := name2func(normWord(funcName))
	if rc != nil {
		return rc, nil
	}
	return nil, fmt.Errorf("%s: not found in the function-list", funcName)
}

// BindKeySymbol assigns function to key by names.
func (editor *KeyMap) BindKeySymbol(keyName, funcName string) error {
	funcValue := name2func(normWord(funcName))
	if funcValue == nil {
		return fmt.Errorf("%s: no such function", funcName)
	}
	return editor.BindKeyFunc(keyName, funcValue)
}

const (
	ansiCursorOff = "\x1B[?25l"

	// On Windows 8.1, the cursor is not shown immediately
	// without SetConsoleCursorPosition by `ESC[u`
	ansiCursorOn = "\x1B[?25h\x1B[s\x1B[u"
)

// CtrlC is the error when Ctrl-C is pressed.
var CtrlC = (errors.New("^C"))

var mu sync.Mutex

func (editor *Editor) getKeyFunction(key string) KeyFuncT {
	code := keys.Code(key)
	if editor.KeyMap.KeyMap != nil {
		if f, ok := editor.KeyMap.KeyMap[code]; ok {
			return f
		}
	}
	if f, ok := GlobalKeyMap.KeyMap[code]; ok {
		return f
	}
	if f, ok := defaultKeyMap[code]; ok {
		return f
	}
	return &KeyGoFuncT{
		Func: func(ctx context.Context, this *Buffer) Result {
			return keyFuncInsertSelf(ctx, this, key)
		},
		Name: key,
	}
}

// ReadLine calls LineEditor
// - ENTER typed -> returns TEXT and nil
// - CTRL-C typed -> returns "" and readline.CtrlC
// - CTRL-D typed -> returns "" and io.EOF
func (editor *Editor) ReadLine(ctx context.Context) (string, error) {
	if editor.Writer == nil {
		editor.Writer = os.Stdout
	}
	if editor.Out == nil {
		editor.Out = bufio.NewWriter(editor.Writer)
	}
	defer func() {
		editor.Out.WriteString(ansiCursorOn)
		editor.Out.Flush()
	}()

	if editor.Prompt == nil {
		editor.Prompt = func() (int, error) {
			editor.Out.WriteString("\n> ")
			return 2, nil
		}
	}
	if editor.History == nil {
		editor.History = new(_EmptyHistory)
	}
	if editor.LineFeed == nil {
		editor.LineFeed = func(Result) {
			editor.Out.WriteByte('\n')
		}
	}
	if editor.Tty == nil {
		editor.Tty = &tty.Tty{}
	}
	buffer := Buffer{
		Editor:         editor,
		Buffer:         make([]Cell, 0, 20),
		historyPointer: editor.History.Len(),
	}

	if err := editor.Tty.Open(); err != nil {
		return "", fmt.Errorf("go-tty.Open: %s", err.Error())
	}
	defer editor.Tty.Close()

	var err error
	buffer.termWidth, _, err = editor.Tty.Size()
	if err != nil {
		return "", fmt.Errorf("go-tty.Size: %s", err.Error())
	}

	var err1 error
	buffer.topColumn, err1 = editor.Prompt()
	if err1 != nil {
		// unable to get prompt-string.
		fmt.Fprintf(buffer.Out, "%s\n$ ", err1.Error())
		buffer.topColumn = 2
	} else if buffer.topColumn >= buffer.termWidth-3 {
		// ViewWidth is too narrow to edit.
		io.WriteString(buffer.Out, "\n")
		buffer.topColumn = 0
	}
	buffer.InsertString(0, editor.Default)
	if buffer.Cursor > len(buffer.Buffer) {
		buffer.Cursor = len(buffer.Buffer)
	}
	buffer.RepaintAfterPrompt()

	cursorOnSwitch := false

	buffer.startChangeWidthEventLoop(buffer.termWidth, editor.Tty.GetResizeNotifier())

	for {
		mu.Lock()
		if !cursorOnSwitch {
			io.WriteString(buffer.Out, ansiCursorOn)
			cursorOnSwitch = true
		}
		buffer.Out.Flush()

		mu.Unlock()
		key1, err := buffer.GetKey()
		if err != nil {
			return "", err
		}
		mu.Lock()

		f := editor.getKeyFunction(key1)

		if fg, ok := f.(*KeyGoFuncT); !ok || fg.Func != nil {
			io.WriteString(buffer.Out, ansiCursorOff)
			cursorOnSwitch = false
			buffer.Out.Flush()
		}
		rc := f.Call(ctx, &buffer)
		if rc != CONTINUE {
			buffer.LineFeed(rc)

			if !cursorOnSwitch {
				io.WriteString(buffer.Out, ansiCursorOn)
			}
			buffer.Out.Flush()
			result := buffer.String()
			mu.Unlock()
			if rc == ENTER {
				return result, nil
			} else if rc == INTR {
				return result, CtrlC
			} else {
				return result, io.EOF
			}
		}
		mu.Unlock()
	}
}
