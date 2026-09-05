package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

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
		if value < 0 || value > utf8.MaxRune || (value >= 0xD800 && value <= 0xDFFF) {
			fmt.Fprintln(os.Stderr, s, "is not a valid Unicode scalar value")
			os.Exit(1)
		}
		r := rune(value)
		fmt.Printf("runewidth(U+%X): %d\n", value, runewidth.RuneWidth(r))
		buffer.WriteRune(r)
	}
	text := buffer.String()
	w, c := moji.MojiWidthAndCountInString(text)
	fmt.Printf("moji.MojiWidthAndCountInString(%s)=%d:width,%d:count\n",
		text, w, c)
}
