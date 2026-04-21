### Q. Why isn’t my input automatically added to history?

A. This is by design.
`go-readline-ny` does not automatically add input to history.

In many real-world use cases (such as shells), the final string to be recorded in history is not always the raw user input. For example:

* History expansion (e.g. `!` substitution)
* Post-processing or normalization
* Skipping registration on errors

If readline automatically added the raw input, applications would need to:

1. Retrieve the just-added entry
2. Modify or replace it
3. Handle rollback on failure

This would introduce unnecessary complexity.

To avoid this, history management is entirely delegated to the application.
You are expected to explicitly call your history container’s `Add` method at the appropriate time.

### Q. Why doesn’t the history interface include an Add method?

A. The history interface is intentionally read-only:

```go
type IHistory interface {
    Len() int
    At(int) string
}
```

Different applications may have different policies for:

* When to record history
* What to record
* How to store it

By not enforcing a write API, the library allows full flexibility in history management.

### Example

If you are using the provided `simplehistory` package, you can update the history like this:

```../examples/example-history.go
package main

import (
    "context"
    "fmt"

    "github.com/nyaosorg/go-readline-ny"
    "github.com/nyaosorg/go-readline-ny/simplehistory"
)

func main() {
    var editor readline.Editor
    h := simplehistory.New()
    editor.History = h
    for {
        text, err := editor.ReadLine(context.Background())
        if err != nil {
            fmt.Printf("ERR=%s\n", err.Error())
            return
        }
        fmt.Printf("TEXT=%s\n", text)
        h.Add(text)

        // editor.History.Add(text)
        // -> compile error: History interface does not have Add method
        //    But the value returned by `simplehistory.New()` does.
    }
}
```

The `Container` type returned by `simplehistory.New()` already implements an `Add` method, so in typical cases you just need to call it after `ReadLine`.
