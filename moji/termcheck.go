package moji

import (
	"os"
)

var (
	isVsCodeTerminal = os.Getenv("VSCODE_PID") != ""

	isWindowsTerminal = os.Getenv("WT_SESSION") != "" && os.Getenv("WT_PROFILE_ID") != "" && !isVsCodeTerminal

	isWezTerm = os.Getenv("TERM_PROGRAM") == "WezTerm"

	isContour = os.Getenv("TERMINAL_NAME") == "contour"

	// SurrogatePairOk is true when the surrogated pair unicode is supported
	// If it is false, <NNNN> is displayed instead.
	SurrogatePairOk = isWindowsTerminal || isWezTerm || isContour

	// ZeroWidthJoinSequenceOk is true when ZWJ(U+200D) is supported.
	// If it is false, <NNNN> is displayed instead.
	ZeroWidthJoinSequenceOk = isWindowsTerminal || isWezTerm || isContour

	// VariationSequenceOk is true when Variation Sequences are supported.
	// If it is false, <NNNN> is displayed instead.
	VariationSequenceOk = isWindowsTerminal || isWezTerm || isContour

	// ModifierSequenceOk is false, SkinTone sequence are treated as two
	// character
	ModifierSequenceOk = isWindowsTerminal || isWezTerm

	// AmbiguousIsWide is true, EastAsianAmbiguous are treated as two width
	AmbiguousIsWide = !isWindowsTerminal
)
