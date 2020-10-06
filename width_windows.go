package readline

import "os"

// TreatAmbiguousWidthAsNarrow will be set true when the unicode character defined as Ambiguous-with uses 2 cells per character in the current terminal.
var TreatAmbiguousWidthAsNarrow = os.Getenv("WT_SESSION") != "" && os.Getenv("WT_PROFILE_ID") != ""
