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
