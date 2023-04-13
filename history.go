package readline

import (
	"context"
)

// IHistory is the interface ReadLine can use as container for history.
// It can be set to Editor.History field
type IHistory interface {
	Len() int
	At(int) string
}

type _EmptyHistory struct{}

// Len always returns zero because the receiver is dummy.
func (*_EmptyHistory) Len() int { return 0 }

// At always returns empty-string because the receiver is dummy.
func (*_EmptyHistory) At(int) string { return "" }

func keyFuncHistoryUp(ctx context.Context, this *Buffer) Result {
	if this.History.Len() <= 0 {
		return CONTINUE
	}
	if this.historyPointer <= 0 {
		if !this.HistoryCycling {
			return CONTINUE
		}
		this.historyPointer = this.History.Len()
	}
	this.historyPointer--
	keyFuncClear(ctx, this)
	if this.historyPointer >= 0 {
		this.InsertString(0, this.History.At(this.historyPointer))
		this.ViewStart = 0
		this.Cursor = 0
		keyFuncTail(ctx, this)
	}
	return CONTINUE
}

func keyFuncHistoryDown(ctx context.Context, this *Buffer) Result {
	if this.History.Len() <= 0 {
		return CONTINUE
	}
	if this.historyPointer+1 > this.History.Len() {
		return CONTINUE
	}
	this.historyPointer++
	keyFuncClear(ctx, this)
	if this.historyPointer < this.History.Len() {
		this.InsertString(0, this.History.At(this.historyPointer))
		this.ViewStart = 0
		this.Cursor = 0
		keyFuncTail(ctx, this)
	}
	return CONTINUE
}
