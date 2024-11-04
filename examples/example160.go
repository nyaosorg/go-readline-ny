package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/simplehistory"
)

func predictByHistory(B *readline.Buffer) string {
	current := strings.TrimSpace(B.String())
	for i := B.History.Len() - 1; i >= 0; i-- {
		h := strings.TrimSpace(B.History.At(i))
		if strings.HasPrefix(h, current) {
			return h[len(current):]
		}
	}
	return ""
}

func main() {
	history := simplehistory.New()

	editor := &readline.Editor{
		PredictColor: [...]string{"\x1B[3;22;34m", "\x1B[23;39m"},
		Predictor:    predictByHistory,
		History:      history,
	}
	for {
		text, err := editor.ReadLine(context.Background())
		if err != nil {
			fmt.Printf("ERR=%s\n", err.Error())
			return
		}
		fmt.Printf("TEXT=%s\n", text)
		history.Add(text)
	}
}
