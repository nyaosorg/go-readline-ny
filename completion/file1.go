package completion

// Deprecated: Use PathComplete, and CmdCompletion2 or CmdCompletionList2
type File struct{}

func (File) Enclosures() string {
	return `"'`
}

func (File) Delimiters() string {
	return "&|><;"
}

func (File) List(field []string) (completionSet []string, listingSet []string) {
	return PathComplete(field)
}
