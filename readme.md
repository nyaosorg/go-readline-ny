go-readline-ny
==============

go-readline-ny is the readline library used in the command line shell [NYAGOS](https://github.com/zetamatta/nyagos).

sample1.go
----------

The most simple sample

```go
package main

import (
    "context"
    "fmt"

    "github.com/zetamatta/go-readline-ny"
)

func main() {
    editor := readline.Editor{}
    text, err := editor.ReadLine(context.Background())
    if err != nil {
        fmt.Printf("ERR=%s\n", err.Error())
    } else {
        fmt.Printf("TEXT=%s\n", text)
    }
}
```

example2.go
-----------

Tiny Shell

```go
package main

import (
    "context"
    "fmt"
    "os"
    "os/exec"
    "strings"

    "github.com/mattn/go-colorable"

    "github.com/zetamatta/go-readline-ny"
    "github.com/zetamatta/go-readline-ny/simplehistory"
)

func main() {
    history := simplehistory.New()

    editor := readline.Editor{
        Prompt:  func() (int, error) { return fmt.Print("$ ") },
        Writer:  colorable.NewColorableStdout(),
        History: history,
    }
    fmt.Println("Tiny Shell. Type Ctrl-D to quit.")
    for {
        text, err := editor.ReadLine(context.Background())

        if err != nil {
            fmt.Printf("ERR=%s\n", err.Error())
            return
        }

        fields := strings.Fields(text)
        if len(fields) <= 0 {
            continue
        }
        cmd := exec.Command(fields[0], fields[1:]...)
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        cmd.Stdin = os.Stdin

        cmd.Run()

        history.Add(text)
    }
}
```
