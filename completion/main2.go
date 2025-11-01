package completion

import (
	"context"

	"github.com/nyaosorg/go-box/v3"
	rl "github.com/nyaosorg/go-readline-ny"
)

type CmdCompletion2 struct {
	Postfix    string
	Delimiter  string
	Enclosure  string
	Candidates func(fieldsBeforeCursor []string) (completionSet []string, listingSet []string)
}

func (C *CmdCompletion2) String() string {
	return "COMPLETION2"
}

func (C *CmdCompletion2) Call(ctx context.Context, B *rl.Buffer) rl.Result {
	Complete(C.Enclosure, C.Delimiter, B, C.Candidates, C.Postfix)
	return rl.CONTINUE
}

type CmdCompletionOrList2 struct {
	Delimiter  string
	Enclosure  string
	Postfix    string
	Candidates func(fieldsBeforeCursor []string) (completionSet []string, listingSet []string)
}

func (C *CmdCompletionOrList2) String() string {
	return "COMPLETION_OR_LIST2"
}

func (C *CmdCompletionOrList2) Call(ctx context.Context, B *rl.Buffer) rl.Result {
	list := Complete(C.Enclosure, C.Delimiter, B, C.Candidates, C.Postfix)
	if len(list) > 0 {
		B.Out.WriteByte('\n')
		box.Println(list, B.Out)
		B.RepaintAll()
	}
	return rl.CONTINUE
}
