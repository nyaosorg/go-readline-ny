package readline_test

import (
	"context"
	"errors"
	"io"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/nyaosorg/go-ttyadapter/auto"

	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/keys"
)

const f = "\U0001F468\u200D\U0001F33E"

type snapShot struct {
	View string
}

func (s *snapShot) Call(_ context.Context, B *readline.Buffer) readline.Result {
	var shot strings.Builder
	start := B.ViewStart
	shot.WriteString(B.SubString(0, start))
	shot.WriteByte('[')
	width := readline.WidthT(0)
	viewWidth := B.ViewWidth()
	isCursorBarDisplayed := false
	isViewEndDisplayed := false
	for i := start; i < len(B.Buffer); i++ {
		if i == B.Cursor {
			shot.WriteByte('|')
			isCursorBarDisplayed = true
		}
		w1 := B.Buffer[i].Moji.Width()
		if width <= viewWidth && width+w1 > viewWidth {
			shot.WriteByte(']')
			isViewEndDisplayed = true
		}
		shot.WriteString(B.Buffer[i].String())
		width += w1
	}
	if !isCursorBarDisplayed {
		shot.WriteByte('|')
	}
	if !isViewEndDisplayed {
		shot.WriteByte(']')
	}
	s.View = shot.String()
	return readline.CONTINUE
}

func (s *snapShot) String() string {
	return "snap-shot"
}

func (s *snapShot) Register(editor *readline.Editor) {
	editor.KeyMap.BindKey(keys.Code(s.String()), s)
}

func tryAll(t *testing.T, texts ...string) (string, []string) {
	var buffer strings.Builder
	editor := readline.Editor{
		Tty:          &auto.Pilot{Text: texts},
		Writer:       &buffer,
		PromptWriter: func(w io.Writer) (int, error) { return 0, nil },
	}
	result, err := editor.ReadLine(context.Background())
	if err != nil {
		t.Fatalf("ERR=%s", err.Error())
		return "", nil
	}

	outputPieces := strings.Split(buffer.String(), "\x1B[?25h\x1B[?25l")
	for i, s := range outputPieces {
		s = strings.ReplaceAll(s, "\x1B[?25h", "")
		s = strings.ReplaceAll(s, "\x1B[?25l", "")
		s = strings.ReplaceAll(s, "\u200D", "<ZWJ>")
		s = strings.ReplaceAll(s, "\x1B", "<ESC>")
		outputPieces[i] = s
	}
	return result, outputPieces
}

func TestKeyFuncBackSpace(t *testing.T) {
	result, outputs := tryAll(t, f, "\b", "x", "\r")
	expect := "x"
	if runtime.GOOS != "windows" || os.Getenv("WT_SESSION") == "" {
		expect = "\U0001F468\u200Dx"
	}
	if result != expect {
		t.Fatalf("TEXT=%s", result)
		return
	}
	if false {
		for _, o := range outputs {
			println(o)
		}
	}
}

func keyTest(t *testing.T, expView string, width int, typed ...string) error {
	t.Helper()

	expText := strings.ReplaceAll(expView, "[", "")
	expText = strings.ReplaceAll(expText, "|", "")
	expText = strings.ReplaceAll(expText, "]", "")

	ss := &snapShot{}
	prtSc := ss.String()
	typed = append(typed, prtSc)
	typed = append(typed, keys.Enter)

	editor := &readline.Editor{
		Tty: &auto.Pilot{
			Width: width + int(readline.ScrollMargin),
			Text:  typed},
		Writer:       io.Discard,
		PromptWriter: func(w io.Writer) (int, error) { return 0, nil },
	}
	ss.Register(editor)
	result, err := editor.ReadLine(context.Background())
	if err != nil {
		if errors.Is(err, readline.CtrlC) || errors.Is(err, io.EOF) {
			return err
		}
		t.Fatalf("ERR=%s", err.Error())
	}
	if result != expText {
		t.Fatalf("text: expect %#v, but %#v", expText, result)
	}
	if ss.View != expView {
		t.Fatalf("view: expect %#v, but %#v", expView, ss.View)
	}
	return err
}

func TestCmdBackspace(t *testing.T) {
	keyTest(t, "[a|]", 80, "a", "b", "\b")
	keyTest(t, "[|]", 80, "a", "b", "\b", "\b")
	keyTest(t, "[|]", 80, "a", "b", "\b", "\b", "\b")
	keyTest(t, "[a|c]", 80, "a", "b", "c", keys.Left, "\b")
	keyTest(t, "[１|３]", 80, "１", "２", "３", keys.Left, "\b")
	keyTest(t, "12[345|]", 4, "1", "2", "3", "4", "5", "6", "\b")
	keyTest(t, "12[34|6]", 4, "1", "2", "3", "4", "5", "6", keys.Left, "\b")
	keyTest(t, "12[|4567]", 4, "1", "2", "3", "4", "5", "6", "7", keys.Left, keys.Left, keys.Left, keys.Left, "\b")
	keyTest(t, "1[|3456]7", 4, "1", "2", "3", "4", "5", "6", "7", keys.Left, keys.Left, keys.Left, keys.Left, keys.Left, "\b")
	keyTest(t, "１２[３４５|]", 8, "１", "２", "３", "４", "５", "６", "\b")
	keyTest(t, "１２[３４|６]", 8, "１", "２", "３", "４", "５", "６", keys.Left, "\b")
	keyTest(t, "１２[|４５６７]", 8, "１", "２", "３", "４", "５", "６", "７", keys.Left, keys.Left, keys.Left, keys.Left, "\b")
	keyTest(t, "１[|３４５６]７", 8, "１", "２", "３", "４", "５", "６", "７", keys.Left, keys.Left, keys.Left, keys.Left, keys.Left, "\b")
}

func TestCmdBackwardChar(t *testing.T) {
	keyTest(t, "[a|b]", 80, "a", "b", keys.CtrlB)
	keyTest(t, "[a|b]", 80, "a", "b", keys.Left)
	keyTest(t, "ab[|cdef]g", 4, "a", "b", "c", "d", "e", "f", "g", keys.Left, keys.Left, keys.Left, keys.Left, keys.Left)
	keyTest(t, "[|abcd]efg", 4, "a", "b", "c", "d", "e", "f", "g", keys.Left, keys.Left, keys.Left, keys.Left, keys.Left, keys.Left, keys.Left)
	keyTest(t, "１２[３|４]", 4, "１", "２", "３", "４", keys.Left)
	keyTest(t, "１２[|３４]", 4, "１", "２", "３", "４", keys.Left, keys.Left)
	keyTest(t, "１[|２３]４", 4, "１", "２", "３", "４", keys.Left, keys.Left, keys.Left)
	keyTest(t, "[|１２]３４", 4, "１", "２", "３", "４", keys.Left, keys.Left, keys.Left, keys.Left)
	keyTest(t, "[|１２]３４", 4, "１", "２", "３", "４", keys.Left, keys.Left, keys.Left, keys.Left, keys.Left)
}

func TestCmdInterrupt(t *testing.T) {
	result := keyTest(t, "[|]", 80, "a", "b", "c", keys.CtrlC)
	if result != readline.CtrlC {
		t.Fatalf("expect CtrlC, but %v", result)
	}
}

func TestCmdBackwardKillWord(t *testing.T) {
	keyTest(t, "[foo   |]", 80, "foo   bar   baz", keys.AltBackspace, keys.AltBackspace)
	keyTest(t, "[foo   |baz]", 80, "foo   bar   baz", keys.AltB, keys.AltBackspace)
	keyTest(t, "[foo   |  baz]", 80, "foo   bar   baz", keys.AltB, keys.Left, keys.Left, keys.AltBackspace)
}


func TestCmdKillWord(t *testing.T) {
	keyTest(t, "[foo  |  baz]", 80, "foo  bar  baz", keys.AltB, keys.AltB, keys.AltD)
	keyTest(t, "[fo|]", 80, "foo", keys.Left, keys.AltD)
}
