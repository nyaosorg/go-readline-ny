package readline

import "os"

var TreatAmbiguousWidthAsNarrow = os.Getenv("WT_SESSION") != "" && os.Getenv("WT_PROFILE_ID") != ""
