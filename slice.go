package readline

func sliceInsert[T any](slice []T, at int, elements ...T) []T {
	// expand buffer (elements... is dummy)
	slice = append(slice, elements...)

	// shift
	copy(slice[at+len(elements):], slice[at:])

	// copy elements
	copy(slice[at:at+len(elements)], elements)

	return slice
}

func sliceDelete[T any](slice []T, at int, n int) []T {
	copy(slice[at:], slice[at+n:])
	return slice[:len(slice)-n]
}
