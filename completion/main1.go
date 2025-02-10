package completion

import (
	"context"

	"github.com/nyaosorg/go-box/v2"
	rl "github.com/nyaosorg/go-readline-ny"
)

// Deprecated.
type Completion interface {
	Delimiters() string
	Enclosures() string
	List(fields []string) (completionSet, listingSet []string)
}

// Deprecated. Use CmdCompletion2
type CmdCompletion struct {
	Completion
	Postfix string
}

func (C CmdCompletion) String() string {
	return "COMPLETION"
}

func (C CmdCompletion) Call(ctx context.Context, B *rl.Buffer) rl.Result {
	Complete(C.Enclosures(), C.Delimiters(), B, C.List, C.Postfix)
	return rl.CONTINUE
}

// Deprecated. Use CmdCompletionOrList2
type CmdCompletionOrList struct {
	Completion
	Postfix string
}

func (C CmdCompletionOrList) String() string {
	return "COMPLETION_OR_LIST"
}

func (C CmdCompletionOrList) Call(ctx context.Context, B *rl.Buffer) rl.Result {
	list := Complete(C.Enclosures(), C.Delimiters(), B, C.List, C.Postfix)
	if len(list) > 0 {
		B.Out.WriteByte('\n')
		box.Print(ctx, list, B.Out)
		B.RepaintAll()
	}
	return rl.CONTINUE
}
