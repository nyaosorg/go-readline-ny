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

- (#2,#3) Add flag: Editor.HistoryCycling
- Color: support Multi SelectGraphcReition Parameters (on v0.6.2)
- Implement short function SGR1,SGR2,SGR3,SGR4 for color setting

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

temporaly replace go-tty with forked version to support emoji in WindowsTerminal 1.5

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
