module github.com/zetamatta/go-readline-ny

go 1.14

require (
	github.com/atotto/clipboard v0.1.2
	github.com/mattn/go-colorable v0.1.7
	github.com/mattn/go-runewidth v0.0.9
	github.com/mattn/go-tty v0.0.3
	golang.org/x/sys v0.0.0-20200922070232-aee5d888a860 // indirect
)

replace github.com/mattn/go-tty v0.0.3 => ./go-tty-fork20200922
