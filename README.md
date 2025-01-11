[![Go Reference](https://pkg.go.dev/badge/github.com/nyaosorg/go-readline-ny.svg)](https://pkg.go.dev/github.com/nyaosorg/go-readline-ny)
[![Go Report Card](https://goreportcard.com/badge/github.com/nyaosorg/go-readline-ny)](https://goreportcard.com/report/github.com/nyaosorg/go-readline-ny)
[![Wiki](https://img.shields.io/badge/Wiki-orange)](https://github.com/nyaosorg/go-readline-ny/wiki)
[![License](https://img.shields.io/badge/License-MIT-red)](https://github.com/nyaosorg/go-readline-ny/blob/master/LICENSE)

The New Yet another Readline for Go (go-readline-ny)
====================================================

**go-readline-ny** is a one-line input library for CUI applications written in Go. It is designed to be extensible to meet diverse needs and has been used in the command-line shell [NYAGOS] for a long time.

[NYAGOS]: https://github.com/nyaosorg/nyagos

- Emacs-like key-bindings
- Input history
- Word completion (file names, command names, or any names in given array)
- Syntax highlighting
- Supported platforms: Windows and Linux
- Full Unicode (UTF8) support, including:
    - Surrogate-pair
    - Emojis (via clipboard)
    - Zero-width joiner (via clipboard)
    - Variation selectors (via clipboard pasted with Ctrl-Y)
- Add-Ons:
    - [SKK] (Japanese input method editor)
    - [Multi-lines Editing][go-multiline-ny]

[SKK]: https://github.com/nyaosorg/go-readline-skk
[go-multiline-ny]: https://github.com/hymkor/go-multiline-ny

![Zero-Width-Joiner sample on Windows-Terminal](./emoji.png)

![](./colorcmdline.png)

examples
--------

### [example1.go](./examples/example1.go)

The most simple sample.

```examples/example1.go
package main

import (
    "context"
    "fmt"

    "github.com/nyaosorg/go-readline-ny"
)

func main() {
    var editor readline.Editor
    text, err := editor.ReadLine(context.Background())
    if err != nil {
        fmt.Printf("ERR=%s\n", err.Error())
    } else {
        fmt.Printf("TEXT=%s\n", text)
    }
}
```

If the target platform includes Windows, you have to import and use [go-colorable](https://github.com/mattn/go-colorable) like example2.go .

### [example2.go](./examples/example2.go)

Tiny Shell. This is a sample of prompt change, colorization, filename completion and history browsing.

```examples/example2.go
package main

import (
    "context"
    "fmt"
    "io"
    "os"
    "os/exec"
    "regexp"
    "strings"

    "github.com/mattn/go-colorable"

    "github.com/nyaosorg/go-readline-ny"
    "github.com/nyaosorg/go-readline-ny/completion"
    "github.com/nyaosorg/go-readline-ny/keys"
    "github.com/nyaosorg/go-readline-ny/simplehistory"
)

func main() {
    history := simplehistory.New()

    editor := &readline.Editor{
        PromptWriter: func(w io.Writer) (int, error) {
            return io.WriteString(w, "\x1B[36;22m$ ") // print `$ ` with cyan
        },
        Writer:  colorable.NewColorableStdout(),
        History: history,
        Highlight: []readline.Highlight{
            {Pattern: regexp.MustCompile("&"), Sequence: "\x1B[33;49;22m"},
            {Pattern: regexp.MustCompile(`"[^"]*"`), Sequence: "\x1B[35;49;22m"},
            {Pattern: regexp.MustCompile(`%[^%]*%`), Sequence: "\x1B[36;49;1m"},
            {Pattern: regexp.MustCompile("\u3000"), Sequence: "\x1B[37;41;22m"},
        },
        HistoryCycling: true,
        PredictColor:   [...]string{"\x1B[3;22;34m", "\x1B[23;39m"},
        ResetColor:     "\x1B[0m",
        DefaultColor:   "\x1B[33;49;1m",
    }

    editor.BindKey(keys.CtrlI, completion.CmdCompletionOrList{
        Completion: completion.File{},
        Postfix:    " ",
    })
    // If you do not want to list files with double-tab-key,
    // use `CmdCompletion` instead of `CmdCompletionOrList`

    fmt.Println("Tiny Shell. Type Ctrl-D to quit.")
    for {
        text, err := editor.ReadLine(context.Background())

        if err != nil {
            fmt.Printf("ERR=%s\n", err.Error())
            return
        }

        fields := strings.Fields(text)
        if len(fields) <= 0 {
            continue
        }
        cmd := exec.Command(fields[0], fields[1:]...)
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        cmd.Stdin = os.Stdin

        cmd.Run()

        history.Add(text)
    }
}
```

### example3.go

- [example3.go](./examples/example3.go)

This is a sample to change key-bindings to diamond cursor.

### example4.go

- [example4.go](./examples/example4.go)

This is a sample that implements the function to start the text editor defined by the environment variable EDITOR and import the edited contents when the ESCAPE key is pressed.

Acknowledgements
----------------

- [@masamitsu-murase](https://github.com/masamitsu-murase) [#1]
- [@ram-on](https://github.com/ram-on) [#2]
- [@glejeune](https://github.com/glejeune) [#6]
- [@brammeleman](https://github.com/brammeleman) [#8]

[#1]: https://github.com/nyaosorg/go-readline-ny/pull/1
[#2]: https://github.com/nyaosorg/go-readline-ny/issues/2
[#6]: https://github.com/nyaosorg/go-readline-ny/pull/6
[#8]: https://github.com/nyaosorg/go-readline-ny/issues/8

Release notes
-------------

- [Release notes (English)](release_note_en.md)
- [Release notes (Japanese)](release_note_ja.md)

License
-------

MIT License
