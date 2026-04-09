package readline

import (
	"context"
	"io"
	"strings"

	"github.com/nyaosorg/go-readline-ny/keys"
	"github.com/nyaosorg/go-readline-ny/moji"
)

var PredictColorBlueItalic = [...]string{"\x1B[3;22;34m", "\x1B[23;39m"}

type predict struct {
	suffix       []Moji
	PredictColor [2]string
	Predictor    func(*Buffer) string
}

func predictByHistory(B *Buffer) string {
	current := B.String()
	for i := B.History.Len() - 1; i >= 0; i-- {
		h := B.History.At(i)
		if len(h) >= len(current) && strings.EqualFold(h[:len(current)], current) {
			return h
		}
	}
	return ""
}

func (P *predict) onAfterRenderPredict(B *Buffer, availWidth int) {
	if P.PredictColor[0] == "" || B.Cursor != len(B.Buffer) {
		return
	}
	if len(B.Buffer) <= 0 {
		P.suffix = nil
	} else if P.Predictor != nil {
		P.suffix = moji.StringToMoji(P.Predictor(B))
	} else {
		P.suffix = moji.StringToMoji(predictByHistory(B))
	}
	if len(P.suffix) <= len(B.Buffer) {
		return
	}
	io.WriteString(B.Out, B.PredictColor[0]) // "\x1B[3;22;39m"
	sfx := P.suffix
	for i := len(B.Buffer); i < len(sfx); i++ {
		c := sfx[i]
		availWidth -= int(c.Width())
		if availWidth < 0 {
			break
		}
		c.PrintTo(B.Out)
	}
	io.WriteString(B.Out, P.PredictColor[1]) // "\x1B[23m"
}

func (P *predict) accept(b *Buffer) {
	if len(b.suffix) <= 0 {
		return
	}
	var s strings.Builder
	for _, c := range b.suffix {
		c.WriteTo(&s)
	}
	b.ReplaceAndRepaint(0, s.String())
}

func (p *predict) install(editor *Editor) {
	if editor.OnAfterRender != nil {
		return
	}
	editor.OnAfterRender = editor.predict.onAfterRenderPredict
	editor.KeyMap.BindKey(keys.Right, CmdForwardCharOrAcceptPredict)
	editor.KeyMap.BindKey(keys.CtrlF, CmdForwardCharOrAcceptPredict)
}

var CmdAcceptPredict = NewGoCommand("ACCEPT_PREDICT", cmdAcceptPredict)

func cmdAcceptPredict(ctx context.Context, b *Buffer) Result {
	b.accept(b)
	return CONTINUE
}

var CmdForwardCharOrAcceptPredict = NewGoCommand("FORWARD_CHAR_OR_ACCEPT_PREDICT", cmdForwardCharOrAcceptPredict)

func cmdForwardCharOrAcceptPredict(ctx context.Context, b *Buffer) Result {
	if b.Cursor < len(b.Buffer) {
		return cmdForwardChar(ctx, b)
	}
	b.accept(b)
	return CONTINUE
}
