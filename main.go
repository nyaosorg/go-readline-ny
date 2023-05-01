package readline

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
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
func (editor *Editor) GetBindKey(key string) Command {
	key = normWord(key)
	if code, ok := name2code[key]; ok {
		return editor.loolupCommand(code.String())
	}
	return nil
}

const (
	ansiCursorOff = "\x1B[?25l"

	// On Windows 8.1, the cursor is not shown immediately
	// without SetConsoleCursorPosition by `ESC[u`
	ansiCursorOn = "\x1B[?25h\x1B[s\x1B[u"
)

// CtrlC is the error when Ctrl-C is pressed.
var CtrlC = errors.New("^C")

var mu sync.Mutex

func (editor *Editor) loolupCommand(key string) Command {
	code := keys.Code(key)
	if f, ok := editor.KeyMap[code]; ok {
		return f
	}
	if f, ok := GlobalKeyMap[code]; ok {
		return f
	}
	return SelfInserter(key)
}

// ReadLine calls LineEditor
// - ENTER typed -> returns TEXT and nil
// - CTRL-C typed -> returns "" and readline.CtrlC
// - CTRL-D typed -> returns "" and io.EOF
func (editor *Editor) ReadLine(ctx context.Context) (string, error) {
	if editor.KeyMap == nil {
		editor.KeyMap = KeyMap{}
	}
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
		editor.History = _EmptyHistory{}
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

	buffer.topColumn, err = editor.Prompt()
	if err != nil {
		// unable to get prompt-string.
		fmt.Fprintf(buffer.Out, "%s\n$ ", err.Error())
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
	buffer.startChangeWidthEventLoop(buffer.termWidth, editor.Tty.GetResizeNotifier())

	for {
		buffer.Out.Flush()

		key, err := buffer.GetKey()
		if err != nil {
			return "", err
		}
		mu.Lock()

		f := editor.loolupCommand(key)

		io.WriteString(buffer.Out, ansiCursorOff)

		rc := f.Call(ctx, &buffer)

		io.WriteString(buffer.Out, ansiCursorOn)

		if rc != CONTINUE {
			buffer.LineFeed(rc)

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
