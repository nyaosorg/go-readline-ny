package moji

import (
	"testing"

	"github.com/mattn/go-runewidth"
)

const benchRunes = `
func BenchmarkDirectRuneWidth(b *testing.B){
	b.ResetTimer()
	for i := 0 ; i < b.N ; i++ {
		for _,c := range benchRunes {
			_ = runewidth.RuneWidth(c)
		}
	}
}

func BenchmarkCachedRuneWidth(b *testing.B){
	b.ResetTimer()
	for i := 0 ; i < b.N ; i++ {
		for _,c := range benchRunes {
			_ = getWidth(c)
		}
	}
}
`

func BenchmarkDirectRuneWidth(b *testing.B){
	b.ResetTimer()
	for i := 0 ; i < b.N ; i++ {
		for _,c := range benchRunes {
			_ = runewidth.RuneWidth(c)
		}
	}
}

func BenchmarkCachedRuneWidth(b *testing.B){
	b.ResetTimer()
	for i := 0 ; i < b.N ; i++ {
		for _,c := range benchRunes {
			_ = getWidth(c)
		}
	}
}
