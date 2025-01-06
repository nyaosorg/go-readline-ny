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
	"unicode/utf8"

	"github.com/nyaosorg/go-readline-ny/internal/moji"
	"github.com/nyaosorg/go-readline-ny/keys"
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

type ITty interface {
	Open(onSize func(int)) error
	GetKey() (string, error)
	Size() (int, int, error)
	Close() error
}

// Editor is the main class to hold the parameter for ReadLine
type Editor struct {
	KeyMap
	History        IHistory
	Writer         io.Writer
	Out            *bufio.Writer
	PromptWriter   func(io.Writer) (int, error)
	Default        string
	Cursor         int
	LineFeedWriter func(Result, io.Writer) (int, error)
	Tty            ITty
	Coloring       Coloring
	HistoryCycling bool
	mutex          sync.Mutex
	PredictColor   [2]string
	Predictor      func(*Buffer) string
}

const (
	ansiCursorOff = "\x1B[?25l"

	// On Windows 8.1, the cursor is not shown immediately
	// without SetConsoleCursorPosition by `ESC[u`
	ansiCursorOn = "\x1B[?25h\x1B[s\x1B[u"
)

var PredictColorBlueItalic = [...]string{"\x1B[3;22;34m", "\x1B[23;39m"}

// CtrlC is the error when Ctrl-C is pressed.
var CtrlC = (errors.New("^C"))

func (editor *Editor) LookupCommand(key string) Command {
	code := keys.Code(key)
	if f, ok := editor.KeyMap.Lookup(code); ok {
		return f
	}
	if f, ok := GlobalKeyMap.Lookup(code); ok {
		return f
	}
	return SelfInserter(key)
}

func cutEscapeSequenceAndOldLine(s string) string {
	buffer := make([]byte, 0, len(s)*2)
	esc := false           // for ESC[...
	titleSequence := false // for ESC]...\x07
	for i, end := 0, len(s); i < end; i++ {
		r := s[i]
		switch r {
		case '\r', '\n':
			buffer = buffer[:0]
		case '\x1B':
			esc = true
		default:
			if titleSequence {
				if r == '\007' {
					titleSequence = false
					esc = false
				}
			} else if esc {
				if r == ']' {
					titleSequence = true
				} else if ('A' <= r && r <= 'Z') || ('a' <= r && r <= 'z') {
					esc = false
				}
			} else if r == '\b' {
				if lastRune, siz := utf8.DecodeLastRune(buffer); lastRune != utf8.RuneError {
					buffer = buffer[:len(buffer)-siz]
				}
			} else {
				buffer = append(buffer, r)
			}
		}
	}
	return string(buffer)
}

func (editor *Editor) callPromptWriter() (int, error) {
	if editor.PromptWriter == nil {
		_, err := editor.Out.WriteString("\n> ")
		return 2, err
	}
	var buffer strings.Builder
	editor.PromptWriter(&buffer)
	prompt := buffer.String()
	_, err := editor.Out.WriteString(prompt)
	w, _ := moji.MojiWidthAndCountInString(cutEscapeSequenceAndOldLine(prompt))
	return int(w), err
}

// Init replaces nil fields to default values.
// When we refer them before calling Readline,
// We have to call Init explicitly.
func (editor *Editor) Init() {
	if editor.Writer == nil {
		editor.Writer = os.Stdout
	}
	if editor.Out == nil {
		if br, ok := editor.Writer.(*bufio.Writer); ok {
			editor.Out = br
		} else {
			editor.Out = bufio.NewWriter(editor.Writer)
		}
	}
	if editor.History == nil {
		editor.History = _EmptyHistory{}
	}
	if editor.Tty == nil {
		editor.Tty = &_Tty{}
	}
	if editor.Coloring == nil {
		editor.Coloring = _MonoChrome{}
	}
	if editor.Predictor == nil {
		editor.Predictor = predictByHistory
	}
}

// ReadLine calls LineEditor
// - ENTER typed -> returns TEXT and nil
// - CTRL-C typed -> returns "" and readline.CtrlC
// - CTRL-D typed -> returns "" and io.EOF
func (editor *Editor) ReadLine(ctx context.Context) (string, error) {
	editor.Init()
	defer func() {
		editor.Out.WriteString(ansiCursorOn)
		editor.Out.Flush()
	}()

	buffer := Buffer{
		Editor:         editor,
		Buffer:         make([]Cell, 0, 20),
		historyPointer: editor.History.Len(),
		suffix:         nil, // moji.StringToMoji("$"),
	}

	onResize := func(w int) {
		editor.mutex.Lock()
		buffer.termWidth = w
		buffer.ResetViewStart()
		buffer.RepaintLastLine()
		editor.mutex.Unlock()
	}

	if err := editor.Tty.Open(onResize); err != nil {
		return "", err
	}
	defer editor.Tty.Close()

	var err error
	buffer.termWidth, _, err = editor.Tty.Size()
	if err != nil {
		return "", err
	}

	buffer.topColumn, err = editor.callPromptWriter()
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

	io.WriteString(buffer.Out, ansiCursorOn)
	for {
		key, err := buffer.GetKey()
		if err != nil {
			return "", err
		}
		editor.mutex.Lock()

		f := editor.LookupCommand(key)

		io.WriteString(buffer.Out, ansiCursorOff)

		rc := f.Call(ctx, &buffer)

		io.WriteString(buffer.Out, ansiCursorOn)

		if rc != CONTINUE {
			if buffer.suffix != nil {
				buffer.suffix = nil
				buffer.repaintWithoutUpdateSuffix()
			}
			if buffer.LineFeedWriter != nil {
				buffer.LineFeedWriter(rc, buffer.Out)
			} else {
				buffer.Out.WriteByte('\n')
			}
			buffer.Out.Flush()

			result := buffer.String()
			editor.mutex.Unlock()
			if rc == ENTER {
				return result, nil
			} else if rc == INTR {
				return result, CtrlC
			} else {
				return result, io.EOF
			}
		}
		editor.mutex.Unlock()
	}
}
