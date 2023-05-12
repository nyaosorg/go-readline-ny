package completion

import (
	"context"
	"strings"
	"unicode"

	"github.com/nyaosorg/go-box/v2"
	rl "github.com/nyaosorg/go-readline-ny"
)

type Completion interface {
	Delimiters() string
	Enclosures() string
	List(fields []string) (fullnames, basenames []string)
}

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
		if strings.IndexRune(q, c) < 0 {
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
			if strings.Index(spaces, c) < 0 {
				break
			}
			i++
			if i >= B.Cursor {
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
				if strings.Index(spaces, c) >= 0 {
					fields = append(fields, removeQuotes(B.SubString(start, i), quotes))
					lastWordStart = start
					break
				}
				if strings.Index(del, c) >= 0 {
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

func complete(quotes, del string, B *rl.Buffer, C Completion) []string {
	fields, lastWordStart := split(quotes, del, B)
	list, baselist := C.List(fields)

	if len(list) <= 0 {
		return nil
	}
	if len(list) == 1 {
		str := list[0]
		if strings.EqualFold(fields[len(fields)-1], str) {
			B.Out.WriteByte('\a')
		} else {
			if len(quotes) > 0 && len(del) > 0 && strings.ContainsAny(str, " \t\r\n\v\f"+del) {
				str = string(quotes[0]) + str + string(quotes[0])
			}
			B.ReplaceAndRepaint(lastWordStart, str+" ")
		}
		return nil
	}
	prefix := commonPrefix(list)
	if strings.EqualFold(fields[len(fields)-1], prefix) {
		B.Out.WriteByte('\a')
		return baselist
	} else {
		if len(quotes) > 0 && hasToInsertQuotation(list, " \t\r\n\v\f"+del) {
			prefix = string(quotes[0]) + prefix
		}
		B.ReplaceAndRepaint(lastWordStart, prefix)
		return nil
	}
}

type CmdCompletion struct {
	Completion
}

func (C CmdCompletion) String() string {
	return "COMPLETION"
}

func (C CmdCompletion) Call(ctx context.Context, B *rl.Buffer) rl.Result {
	complete(C.Enclosures(), C.Delimiters(), B, C)
	return rl.CONTINUE
}

type CmdCompletionOrList struct {
	Completion
}

func (C CmdCompletionOrList) String() string {
	return "COMPLETION_OR_LIST"
}

func (C CmdCompletionOrList) Call(ctx context.Context, B *rl.Buffer) rl.Result {
	list := complete(C.Enclosures(), C.Delimiters(), B, C)
	if list != nil && len(list) > 0 {
		B.Out.WriteByte('\n')
		box.Print(ctx, list, B.Out)
		B.RepaintAll()
	}
	return rl.CONTINUE
}
