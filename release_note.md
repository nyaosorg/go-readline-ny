v0.8.2
======
12,Aug. 2022

- Reset color before printing the first character

v0.8.1
======
25,Jun. 2022

- Fix: On Ctrl-E typed, sometimes non-space character remains on the cursor.

v0.8.0
======
29 Apr. 2022

- Enable surrogate-pair on WezTerm
- Enable surrogate-pair, ZWJ,and VS on Contour Terminal
- Do not use BACKSPACE(\b) as output

v0.7.0
======
26 Feb. 2022

- Coloring.Init() has to return default colors (This interface's compatibility is broken)
- Use ESC[49m (default bgcolor) instead of ESC[40m

v0.6.3
======
29 Dec. 2021

- (#2,#3) Add flag: Editor.HistoryCycling
- Color: support Multi SelectGraphcReition Parameters (on v0.6.2)
- Implement short function SGR1,SGR2,SGR3,SGR4 for color setting

v0.6.1
======
10 Dec. 2021

- Support color.

v0.5.0
======
12 Sep. 2021

- Change owner: zetamatta to nyaosorg

v0.4.14
=======
27 Aug. 2021

- (nyagos-412) The widths of Box Drawing (U+2500-257F) were incorrect on the legacy terminals (on not Winows Terminal in Windows10)  
    for [East Asian Ambiguous Character ・ Issue #412 ・ zetamatta/nyagos](https://github.com/zetamatta/nyagos/issues/412)

v0.4.13
=======
05 Jul. 2021

- Support Mathematical Bold Capital (U+1D400 - U+1D7FF) on the Windows Terminal

v0.4.12
=======
03 May. 2021

- Disable the surrogate-pair in the terminal of VisualStudioCode because it is not supported.
This problem surfaces when we start VSCode from the WindowsTerminal.

v0.4.11
=======

- Support Emoji Moifier Sequence (skin tone) : something with &x1F3FB;(U+1F3FB)～ &x1F3FF;(U+1F3FF)
