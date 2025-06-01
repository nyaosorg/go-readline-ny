package completion

import (
	"strings"
	"unicode"

	rl "github.com/nyaosorg/go-readline-ny"
)

func commonPrefix(list []string) string {
	if len(list) < 1 {
		return ""
	}
	common := list[0]
	var cr, fr *strings.Reader
	// to make English case of completed word to the shortest candidate
	minimumLength := len(list[0])
	minimumIndex := 0
	for index, f := range list[1:] {
		if cr != nil {
			cr.Reset(common)
		} else {
			cr = strings.NewReader(common)
		}
		if fr != nil {
			fr.Reset(f)
		} else {
			fr = strings.NewReader(f)
		}
		i := 0
		var buffer strings.Builder
		for {
			ch, _, err := cr.ReadRune()
			if err != nil {
				break
			}
			fh, _, err := fr.ReadRune()
			if err != nil || unicode.ToUpper(ch) != unicode.ToUpper(fh) {
				break
			}
			buffer.WriteRune(ch)
			i++
		}
		common = buffer.String()
		if len(f) < minimumLength {
			minimumLength = len(f)
			minimumIndex = index + 1
		}
	}
	return list[minimumIndex][:len(common)]
}

func removeQuotes(s, q string) string {
	var buffer strings.Builder
	for _, c := range s {
		if !strings.ContainsRune(q, c) {
			buffer.WriteRune(c)
		}
	}
	return buffer.String()
}

func split(quotes, del string, B *rl.Buffer) (fields []string, lastWordStart int) {
	i := 0
	const spaces = " \t\r\n\v\f"

	for i < B.Cursor {
		// skip space
		for {
			c := B.Buffer[i].String()
			if !strings.Contains(spaces, c) {
				break
			}
			i++
			if i >= B.Cursor {
				fields = append(fields, "")
				return
			}
		}
		start := i
		bits := 0
		for {
			c := B.Buffer[i].String()
			if j := strings.Index(quotes, c); j >= 0 {
				bits ^= (1 << j)
			} else if bits == 0 {
				if strings.Contains(spaces, c) {
					fields = append(fields, removeQuotes(B.SubString(start, i), quotes))
					lastWordStart = start
					break
				}
				if strings.Contains(del, c) {
					fields = append(fields, removeQuotes(B.SubString(start, i), quotes))
					fields = append(fields, c)
					lastWordStart = i
					i++
					break
				}
			}
			i++
			if i >= B.Cursor {
				fields = append(fields, removeQuotes(B.SubString(start, i), quotes))
				lastWordStart = start
				return
			}
		}
	}
	return
}

func hasToInsertQuotation(list []string, spaceAndSoOn string) bool {
	for _, s := range list {
		if strings.ContainsAny(s, spaceAndSoOn) {
			return true
		}
	}
	return false
}

func removeUnmatches(full, base []string, source string) (newFull, newBase []string) {
	for i, name := range full {
		if len(name) >= len(source) && strings.EqualFold(source, name[:len(source)]) {
			newFull = append(newFull, name)
			newBase = append(newBase, base[i])
		}
	}
	return
}

func Complete(quotes, del string, B *rl.Buffer, getCandidates func([]string) ([]string, []string), postfix string) []string {
	fields, lastWordStart := split(quotes, del, B)
	if len(fields) == 0 {
		return nil
	}
	q := ""
	if len(quotes) > 0 {
		q = quotes[:1]
	}
	if qq := B.Buffer[lastWordStart].String(); strings.Contains(quotes, qq) {
		q = qq
	}
	list, baselist := getCandidates(fields)
	if baselist == nil || len(baselist) <= 0 {
		baselist = list
	}
	list, baselist = removeUnmatches(list, baselist, fields[len(fields)-1])

	if len(list) <= 0 {
		return nil
	}
	if len(list) == 1 {
		str := list[0]
		if len(quotes) > 0 && len(del) > 0 && strings.ContainsAny(str, " \t\r\n\v\f"+del) {
			str = q + str + q
		}
		B.ReplaceAndRepaint(lastWordStart, str+postfix)
		return nil
	}
	prefix := commonPrefix(list)
	if strings.EqualFold(fields[len(fields)-1], prefix) {
		B.Out.WriteByte('\a')
		return baselist
	} else {
		if len(quotes) > 0 && hasToInsertQuotation(list, " \t\r\n\v\f"+del) {
			prefix = q + prefix
		}
		B.ReplaceAndRepaint(lastWordStart, prefix)
		return nil
	}
}
