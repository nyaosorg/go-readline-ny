( **English** / [Japanese](release_note_ja.md) )

- Updated [go-ttyadapter] to [v0.1.0] and adjusted to its API changes [#12]
- Addressed warnings reported by Go Report Card [#13]

[#12]: https://github.com/nyaosorg/go-readline-ny/pull/12
[#13]: https://github.com/nyaosorg/go-readline-ny/pull/13
[go-ttyadapter]: https://github.com/nyaosorg/go-ttyadapter
[v0.1.0]: https://github.com/nyaosorg/go-ttyadapter/releases/tag/v0.1.0

v1.12.1
=======
Nov 2, 2025

- Fixed a panic in completion.PathComplete that occurred when the last token was an empty string. It now correctly returns all files in the current directory. (#10)
- Adapt to the API changes in go-box v3 (#11)

v1.12.0
=======
Nov 1, 2025

- The meaning of `Buffer.ViewWidth()` has been adjusted: it now returns the maximum number of columns that can be displayed at once without scrolling, rather than the text length at which scrolling starts.
- Migrated the terminal input subpackages as follows:
    - `tty8` → "[github.com/nyaosorg/go-ttyadapter]/tty8" for [github.com/mattn/go-tty]
    - `tty10` → "[github.com/nyaosorg/go-ttyadapter]/tty10" for [golang.org/x/term]
    - `auto` → "[github.com/nyaosorg/go-ttyadapter]/auto"

[github.com/mattn/go-tty]: https://github.com/mattn/go-tty
[golang.org/x/term]: https://pkg.go.dev/golang.org/x/term
[github.com/nyaosorg/go-ttyadapter]: https://github.com/nyaosorg/go-ttyadapter

v1.11.0
=======
Oct 24, 2025

- Exported the wrapper type of `"github.com/mattn/go-tty".TTY` as `"github.com/nyaosorg/go-readline-ny/tty8".Tty`, allowing key input (including control sequences) to be handled per key rather than per rune.

v1.10.0
=======
Oct 18, 2025

- Implemented the following changes in response to Issue [#452] in nyaosorg/nyagos:
  - Added symbolic identifiers and string names for keys  
    | Identifier\*1        | Symbol\*2       | Key Combination   |
    | ------------------- | -------------- | ----------------- |
    | `keys.CtrlPageDown` | `"C_PAGEDOWN"` | `Ctrl`+`PageDown` |
    | `keys.CtrlPageUp`   | `"C_PAGEUP"`   | `Ctrl`+`PageUp`   |
    | `keys.CtrlHome`     | `"C_HOME"`     | `Ctrl`+`Home`     |
    | `keys.CtrlEnd`      | `"C_END"`      | `Ctrl`+`End`      |
  - Added default key bindings  
    |Key Combination| Function |
    |---------------|----------|
    |`Ctrl`+`Home` | Delete from the beginning of the line to the cursor (same as `Ctrl`+`U`) |
    |`Ctrl`+`End` | Delete from the cursor to the end of the line (same as `Ctrl`+`K`) |
- Updated [mattn/go-runewidth] from v0.0.16 to v0.0.19

\*1 Constants defined in `"github.com/nyaosoorg/go-readline-ny/keys"`  
\*2 Keys of the mapping `"keys".NameToCode` that associates symbolic strings with key sequences

[#452]: https://github.com/nyaosorg/nyagos/issues/452
[mattn/go-runewidth]: https://github.com/mattn/go-runewidth

v1.9.1
======
Jun 25, 2025

- Fixed a build error caused by multiple `main` functions under the `examples` directory when running `go test ./...`.
- Completion candidates are now listed even when the word to complete is empty.

v1.9.0
======
Feb 12, 2025

- Fix: Completion failed to insert a space after the word when there was only one candi
date.
- Fix: The second and subsequent characters in .Enclosure(s) were not used.
- Add a new completion feature: complete.CmdCompletionList2 and CmdCompletion2

v1.8.0
======
Feb 9, 2025

- By default, avoid using the operating system’s clipboard as the cut buffer. It is now possible to replace the cut buffer with the `Clipboard` field of the `Editor` type. [#9]
- Export the package `Moji` and its types that were previously internal.

[#9]: https://github.com/nyaosorg/go-readline-ny/issues/9

v1.7.4
======
Jan 28, 2025

- Experimental change: The value `-1 - (the byte position of the cursor)` is now given as the second parameter instead of `-1` when calling the `FindAllStringIndex` method of the `Pattern` field in the syntax highlighting structure.

v1.7.3
======
Jan 23, 2025

- Improved processing time when the syntax highlighting evaluation function is slow by reducing the number of its calls

v1.7.2
======
Jan 22, 2025

- Made the following changes for go-multiline-ny:
    - Added a hook: `Editor.AfterCommand`, which is called after executing commands.
    - Exported previously unexported types and functions: NewEscapeSequenceId, EscapeSequenceId, HighlightToColoring, HighlightColorSequence, and ColorInterface.

v1.7.1
======
Jan 18, 2025

- Fixed a panic in the syntax highlighting code when Japanese characters were entered using go-readline-skk.
- Deprecated the `GetKey` function.
- Unexported the `GetKeys` and `GetRawKey` functions.

v1.7.0
======
Jan 12, 2025

- Introduced a new interface for syntax highlighting.
    - Pattern specifies the range to which the highlight is applied. A `regexp.Regexp` is acceptable, but any type is sufficient as long as it has `FindAllStringIndex(string, int) [][]int`.
    - `Sequence` specifies the escape sequence to apply to the `Pattern` parts.
    - `DefaultColor` specifies the sequence used for parts where no highlight is applied.
    - `ResetColor` specifies the sequence used when completing the output of the input text.
    - The existing syntax highlighting interface, `Coloring`, is planned to be deprecated.

```go
editor := &readline.Editor{
    // :
    Highlight: []readline.Highlight{
        {Pattern: regexp.MustCompile("&"), Sequence: "\x1B[33;49;22m"},
        {Pattern: regexp.MustCompile(`"[^"]*"`), Sequence: "\x1B[35;49;22m"},
        {Pattern: regexp.MustCompile(`%[^%]*%`), Sequence: "\x1B[36;49;1m"},
        {Pattern: regexp.MustCompile("\u3000"), Sequence: "\x1B[37;41;22m"},
    },
    ResetColor:     "\x1B[0m",
    DefaultColor:   "\x1B[33;49;1m",
}
```

v1.6.3
======
Jan 7, 2025

- Adjusted: Avoided double nesting of `bufio.Writer`.
- Fixed: default color setting was incorrectly set to bold. [#8]

Thanks to [@brammeleman]

[@brammeleman]: https://github.com/brammeleman
[#8]: https://github.com/nyaosorg/go-readline-ny/releases/tag/v1.6.3

v1.6.2
======
Nov 20, 2024

- Fix: in incremental search on UNIX-like platforms, Backspace-key `\x7F` did not work to remove the previous character (though Ctrl-H worked)

v1.6.1
======
Nov 8, 2024

- Fix: some text was missing when pasting multi-lines using the terminal feature of Linux Desktop (for [hymkor/go-multiline-ny v0.16.2](https://github.com/hymkor/go-multiline-ny/releases/tag/v0.16.2))

v1.6.0
======
Nov 4, 2024

- Enable to replace the function to predict (`(*Editor) Predictor = ...`)

v1.5.0
======
Oct 6, 2024

- Implement the prediction like PowerShell 7

It is enabled with setting the escape sequences at starting and ending for drawing predicted text

```
editor := &readline.Editor{
    PredictColor: [...]string{"\x1B[3;22;34m", "\x1B[23;39m"},
}
```

v1.4.1
======
Oct 2, 2024

- Enable to replace the function to calcurate the width of Zero-Width-Join-Sequence for WindowsTerminal 1.22 (SetZWJSWidthGetter)

v1.4.0
======
Jun 16, 2024

- Ctrl-P/N: save the modified entry when switching history, and restore when switching again, until Enter is pressed

v1.3.1
======
Apr 21, 2024

- [#6] Fix: Sub package completion fails on empty field  
  (Thanks [@glejeune])

[#6]: https://github.com/nyaosorg/go-readline-ny/pull/6
[@glejeune]: https://github.com/glejeune

v1.3.0
======
Apr 17, 2024

- Add constant: keys.ShiftTAB = "\x1B[Z"
- Simplify the terminal interface and implements (Compatibility around `ITty` is broken)

v1.2.0
======
Feb 29, 2024

- "keys": Key name constants are now untyped (originally `keys.Code`)

v1.1.0
======
Feb 27, 2024

- "completion": Append the value of CmdCompletion.Postfix or CmdCompletionOrList.Postfix instead of one space when there is only one candidate. (The default value is empty string)

v1.0.1
======
Oct 08, 2023

- Fix: the color can not be changed where the charactor is not simple codepoint such as ZERO WIDTH JOIN SEQUENCE, VARIATION SELECTOR SEQUENCE...

v1.0.0
======
Oct 06, 2023

- Just changed the version to v1.0.0

v0.15.2
=======
Oct 02, 2023

- Fix: Coloring.Next(CursorPositionDummyRune) was not called when the cursor is at the end of the string
- Add `(ColorSequence) Chain` that joins two instances of `ColorSequence`

v0.15.1
=======
Oct 01, 2023

- Implement `(ColorSequence) Add`
- `Coloring.Next` recieves CursorPositionDummyRune(U+E000) on the cursor position now

v0.15.0
=======
Sep 29, 2023

- Remove the deprecated fields, methods and functions for v1.0.0
    - `KeyGoFuncT`. Use `GoCommand` instead
    - `moji.GetCharWidth`
    - `GetFunc`. Use `nameutils.GetFunc` instead
    - `(*Editor) LineFeed`. use `(*Editor) LineFeedWriter` instead
    - `(*Editor) Prompt`. use `(*Editor) PromptWriter` instead
    - `(*Editor) GetBindKey`
    - `(*KeyMap) BindKeyFunc`. Use `nameutils.BindKeyFunc()` instead
    - `(*KeyMap) BindKeyClosure`
    - `(*KeyMap) GetBindKey`
    - `(*KeyMap) BindKeySymbol`. Use `nameutils.BindKeySymbol` instead

v0.14.1
=======
Sep 13, 2023

- Publish the function GetKey(*tty*)
    - The parameter *tty* is expected to be set to a type of '[go-tty.TTY]' or a compatible type. It must have methods: Raw(), ReadRune(), and Buffered().
- Set `Deprecated` comment on the field `Editor.Prompt`

[go-tty.TTY]: https://pkg.go.dev/github.com/mattn/go-tty#TTY

v0.14.0
=======
Oct 28, 2023

- Even if `(*Editor) PromptWriter` outputs `Ctrl-H` or `ESC]...\007`, count the width of the prompt correctly now
- Implement `(*Buffer) RepaintLastLine()` that outputs the last line of the prompt and user input-text. It outputs prompt in which `\n` are replaced to `\r`.

v0.13.2
=======
Jul 29, 2023

- Fix the literals that should be written as `\x` were `0x` on keys/code.go
- Fix cursor wouldn't appear on startup when called with cursor off

v0.13.1
=======
May 28, 2023

- Add new method: `(*KeyMap) Lookup`

v0.13.0
=======
May 19, 2023

- Tab characters can now be represented by a few spaces up to every fourth position instead of ^I

v0.12.3
=======
May 16, 2023

- Add new method: `(*Editor) Init`
    - It replaces nil fields to default values.
      When we refer Editor's fields before calling `(*Editor) Readline`,
      we have to call `(*Editor) Init` explicitly.

v0.12.2
========
May 15, 2023

- Fix: completion.File failed when the path did not contain a directory

v0.12.1
========
May 15, 2023

- CmdCompletion and CmdCompletionOrList narrows down candidates now. So Completion interface side does not have to do
- Completion.List can now be omitted by setting basenames to nil
- completion.File: Fixed: filename completion did not match anyone when ./ is included in the path because the filepath package removes ./ in the path.

v0.12.0
=======
May 13, 2023

- Reimported an improved subset of nyagos completion as a subpackage `completion`
- Rewrite `examples/example2.go` to use the sub-package `completion`
- Add a field LineFeedWriter.

v0.11.7
=======
May 12, 2023

- Change the global sync.Mutex variable to a field in an `Editor` instance. In a command bound to a key, it was imposible to create a new `Editor` instance and call `(*Editor) ReadLine`.
- Rename `(*Editor) loolup(KEY)` to look up a command mapped to a KEY from both instance's table and global table to `LookUp` (exposed)
- `(KeyMap) BindKey(KEY,nil)` now removes the function assigned to KEY
- Implement `(c Cell) String` that behaves equivalently to `b:=&strings.Builder{};c.Moji.WriteTo(b);b.String()`
- `(*Buffer) GetKey` now calls `(*Buffer).Out.Flush`, so the user no longer needs to call `flush` explicitly.

v0.11.6
=======
May 8, 2023

- Reduced memory allocation counts for functions `StringToMoji` and `GetStringWidth`
- Implement a method `MojiCountInString` that counts the number of `Moji` in a string

v0.11.5
=======
May 7, 2023

- (#5) Fix Coloring is wrong on teraterm connecting to Ubuntu at executing `make demo`

v0.11.4
=======
May 6, 2023

- Update go-tty to v0.0.5 for https://github.com/hymkor/go-multiline-ny/issues/1

v0.11.3
=======
May 5, 2023

- `(*Editor)`: Add an new field `PromptWriter func(io.Writer)(int,error)`

v0.11.2
=======
May 1, 2023

- Remove `(*Buffer) Write`
- Hide `(*Buffer) RefreshColor()`
- Sub package: `tty10`: fix goroutine leak
- Add type: `AnonymousCommand` and `SelfInserter`

v0.11.1
=======
Apr 28, 2023

- Create sub-package: `keys` that defines key codes
- Rename `KeyGoFuncT` to `GoCommand`
- Rename `KeyFuncT` to `Command`
- Hide `GetKey(ITty)`
- Remove Key name constants: `K_mmmm`. Use `keys.*`
- Remove Command name constants: `F_mmmm`. Use `(Command) String()` method
- Remove `(Result) String()`

v0.11.0
=======
Apr 26, 2023

- Remove the fork version of [go-tty] and use the original one v0.0.4.  
  Because Windows Terminal's bug was fixed that is the reason to fork.  
  Test:
    - &#x2460; OK: CIRCLE DIGIT ONE: U+2460
    - &#x1F468;&#x200D;&#x1F33E; OK: FARMER: MAN(U+1F468)+ZERO WIDTH JOINER(U+200D)+EAR OF RICE(U+1F33E)
    - &#x908A;&#xE0104; OK: KANJI with VARIATION SELECTOR(U+908A U+E0104)
- Add internal switch to use "golang.org/[x/term]" instead of "[go-tty]". Currently [go-tty] is used.
- Add internal switch to use "golang.org/[x/text/width]" instead of "[go-runewidth]". Currently [go-runewidth] is used.
- Remove the variable SurrogatePairOk. Use functions EnableSurrogatePair() and IsSurrogatePairEnabled()
- Remove the function NewDefaultTty(). Use golang.org/[x/term] or [go-tty]

[go-tty]: https://github.com/mattn/go-tty
[go-runewidth]: https://github.com/mattn/go-runewidth
[x/term]: https://pkg.go.dev/golang.org/x/term
[x/text/width]: https://pkg.go.dev/golang.org/x/text/width

v0.10.1
=======
Apr 14, 2023 ( Used in nyagos-4.4.13\_2 )

- Fix: some constants for color were broken at v0.10.0

v0.10.0
=======
Apr 13, 2023

- Change type `Coloring` interface
    - Init and Next() returns ColorSequence(int64) instead of int32
    - It can output `ESC[0m` now

v0.9.1
======
Apr 10, 2023

- Fix: a trash text 'm' was printed when `Coloring.Init()` / `Next()` returns 0

v0.9.0
=======
Apr 9, 2023

- Rename: `_Cell` (unexported type) to `Cell` (exported). A instance of Cell contains a set of code points (=`Moji`) and color information (unexported) for one gryph.
- Rename: `string2moji`(unexported function) to `StringToMoji` (exported) that converts string to an array of `Moji`. A instance of `Moji` contains code points for one gryph.

v0.8.5
======
Apr 2, 2023

- Fix: imcompatibility on v0.8.4  
  Restore the removed method `(*KeyMap)GetBindKey(string)`  
  because v0.8.4 failed to link on nyagos

v0.8.4
======
Mar 25, 2023

- Fix: GetBindkey returned nil when key is in default state

v0.8.3
======
Sep 24, 2022 ( Used in nyagos-4.4.13\_0 )

- Sample color: vimbatch: change foreground color `ESC[37m` to `ESC[39m`

v0.8.2
======
Aug 12, 2022

- Reset color before printing the first character

v0.8.1
======
Jun 25, 2022

- Fix: On Ctrl-E typed, sometimes non-space character remains on the cursor.

v0.8.0
======
Apr 29, 2022

- Enable surrogate-pair on WezTerm
- Enable surrogate-pair, ZWJ,and VS on Contour Terminal
- Do not use BACKSPACE(`\b`) as output

v0.7.0
======
Feb 26, 2022

- Coloring.Init() has to return default colors (This interface's compatibility is broken)
- Use `ESC[49m` (default bgcolor) instead of `ESC[40m`

v0.6.3
======
Dec 29, 2021

- ([#2],[#3]) Add flag: Editor.HistoryCycling
- Color: support Multi SelectGraphcReition Parameters (on v0.6.2)
- Implement short function SGR1,SGR2,SGR3,SGR4 for color setting

[#2]: https://github.com/nyaosorg/go-readline-ny/issues/2
[#3]: https://github.com/nyaosorg/go-readline-ny/issues/3

Thanks to [@ram-on]

[@ram-on]: https://github.com/ram-on

v0.6.1
======
Dec 10, 2021

- Support color.

v0.5.0
======
Sep 12, 2021

- Change owner: zetamatta to nyaosorg

v0.4.14
=======
Aug 27, 2021

- (nyagos-412) The widths of Box Drawing (U+2500-257F) were incorrect on the legacy terminals (on not Winows Terminal in Windows10)  
    for [East Asian Ambiguous Character ・ Issue #412 ・ zetamatta/nyagos](https://github.com/zetamatta/nyagos/issues/412)

v0.4.13
=======
Jul 5, 2021

- Support Mathematical Bold Capital (U+1D400 - U+1D7FF) on the Windows Terminal

v0.4.12
=======
May 3, 2021

- Disable the surrogate-pair in the terminal of VisualStudioCode because it is not supported.
This problem surfaces when we start VSCode from the WindowsTerminal.

v0.4.11
=======
Apr 14, 2021

- Support Emoji Moifier Sequence (skin tone) : something with &#x1F3FB;(U+1F3FB)～ &#x1F3FF;(U+1F3FF)

v0.4.10
=======
Apr 14, 2021

- Fix the problem the keyup code is entered which was pressed before calling .ReadLine method

v0.4.9
=======
Apr 14, 2021

- Support RAINBOW FLAG (U+1F3F3 U+200D U+1F308 &#x1F3F3;&#x200D;&#x1F308;)

v0.4.8
=======
Apr 14, 2021

- WAVING WHITE FLAG and its variations (U+1F3F3 &amp; U+1F3F3 U+FE0F / &#x1F3F3; &amp; &#x1F3F3;&#xFE0F;)

v0.4.7
=======
Apr 14, 2021

- Support REGIONAL INDICATOR (U+1F1E6 "&#x1F1E6;"..U+1F1FF "&#x1F1FF;" )

v0.4.6
=======
Feb 27, 2021

- Support editing COMBINING ENCLOSING KEYCAP after Variation Selector (&#x0023;&#xFE0F;&#x20E3;) in WindowsTerminal

v0.4.5
=======
Feb 27, 2021

Variation Selector Sequence can include ZeroWidthJoinerSequence for Emoji:WOMAN FACEPALMING

v0.4.4
=======
Feb 27, 2021

Fix: the view was broken when ANYONE + C-b + MANFARMER(or any ZeroWidthJoin Sequence) not via clipboard

v0.4.3
=======
Feb 17, 2021

Fix: CIRCLED DIGITS (e.g. ①) could not be input in WindowsTerminal 1.5

v0.4.2
=======
Feb 14, 2021

- include forked go-tty into internal directory

v0.4.1
=======
Feb 14, 2021

- Use the forked version of go-tty permanetly

v0.4.0
=======
Feb 11, 2021

- temporaly replace go-tty with forked version to support emoji in WindowsTerminal 1.5
- [#1] Support Alt-/ key bind

Thanks to [@masamitsu-murase]

[@masamitsu-murase]: https://github.com/masamitsu-murase
[#1]: https://github.com/nyaosorg/go-readline-ny/pull/1

v0.3.0
=======
Jan 11, 2021

Support Variation Selectors 

See also
-  [異体字セレクタ - Wikipedia](https://ja.wikipedia.org/wiki/%E7%95%B0%E4%BD%93%E5%AD%97%E3%82%BB%E3%83%AC%E3%82%AF%E3%82%BF)
- [UTS #37: Unicode Ideographic Variation Database](https://www.unicode.org/reports/tr37/)

v0.2.8
=======
Dec 13, 2020

Fix for [入力した文字がプロンプトにめり込む状態で CTRL+W で文字を消していくと Panic をおこして nyagos が落ちます ・ Issue #396 ・ zetamatta/nyagos](https://github.com/zetamatta/nyagos/issues/396)

v0.2.7
=======
Nov 20, 2020

Use go-tty includingthe patch on [Fix: the first key after terminal-window activated was input twice. by zetamatta ・ Pull Request #40 ・ mattn/go-tty](https://github.com/mattn/go-tty/pull/40)

v0.2.6
=======
Nov 15, 2020

- Temporarily replace mattn/go-tty by zetamatta/go-tty for https://github.com/zetamatta/nyagos/issues/393

v0.2.4
=======
Nov 14, 2020

- Ctrl-Y: trim the last CRLF on pasting

v0.2.3
======
Oct 23, 2020

- Incremental Search: compare case-insensitively

v0.2.2
=======
Oct 9, 2020

- Fix the cursor position broken on the unicode VARIATION SELECTOR-1..16

v0.2.1
=======
Oct 6, 2020

Support Zero-Width-Join Sequence.

v0.1.10
=======
Sep 29, 2020

Refer to https://github.com/mattn/go-tty/commit/b0f19ff1ae2faa49f3be9c0f304009b2cde03b97

v0.1.9
=======
Sep 27, 2020

Refer zetamatta/go-tty (fork of mattn/go-tty on Sep.22,2020 )
( Revert to v0.1.5 )

v0.1.8
=======
Sep 27, 2020

Remove go-tty-fork2290922's go.mod / go.sum

v0.1.7
=======
Sep 27, 2020

Add TreatAmbiguousWidthAsNarrow And IsSurrogatePairOk again.

v0.1.6
=======
Sep 27, 2020

Embed fork of go-tty

v0.1.5
=======
Sep 23, 2020

Regard Ambiguous width character as 1 cell on Windows Terminal

v0.1.4
=======
Sep 23, 2020

Support surrogat pair

v0.1.3
======
Sep 22, 2020

Use zetamatta/go-tty instead of mattn/go-tty temporary

v0.1.2
======
May 24, 2020

Fix width of '<7F>' (when Ctrl-V and Ctrl-H are typed)

v0.1.1
======
May 24, 2020

Support SurrogatePaired Letter on WindowsTerminal

v0.1.0
======
Mar 21, 2020


Copy from github.com/zetamatta/nyagos/readline
