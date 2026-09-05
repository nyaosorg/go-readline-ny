package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/nyaosorg/go-readline-ny/moji"
)

func main() {
	fmt.Println("SurrogatePairOk:", moji.SurrogatePairOk)
	fmt.Println("ZeroWidthJoinSequenceOk:", moji.ZeroWidthJoinSequenceOk)
	fmt.Println("VariationSequenceOk:", moji.VariationSequenceOk)
	fmt.Println("ModifierSequenceOk:", moji.ModifierSequenceOk)
	fmt.Println("AmbiguousIsWide:", moji.AmbiguousIsWide)

	var buffer strings.Builder
	for _, s := range os.Args[1:] {
		value, err := strconv.ParseInt(s, 16, 32)
		if err != nil {
			fmt.Fprintln(os.Stderr, s, err.Error())
			os.Exit(1)
		}
		fmt.Printf("runewidth(U+%X): %d\n", value, runewidth.RuneWidth(rune(value)))
		buffer.WriteRune(rune(value))
	}
	text := buffer.String()
	w, c := moji.MojiWidthAndCountInString(text)
	fmt.Printf("moji.MojiWidthAndCountInString(%s)=%d:width,%d:count\n",
		text, w, c)
}
