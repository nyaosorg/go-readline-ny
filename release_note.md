- Fix: a trash text 'm' was printed when Coloring.Init/Next returns 0

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

- Support Emoji Moifier Sequence (skin tone) : something with &x1F3FB;(U+1F3FB)～ &x1F3FF;(U+1F3FF)
