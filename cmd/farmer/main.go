package main

import (
	"fmt"
	"time"
)

func main() {
	farmer := "\U0001F468\u200D\U0001F33E"
	fmt.Printf("(a) %s%s%s:OK\n", farmer, farmer, farmer)
	fmt.Printf("(b) %s\b\b%s\b\b%s:NG\n", farmer, farmer, farmer)
	fmt.Printf("(c) %s\b%s\b%s:NG\n", farmer, farmer, farmer)
	fmt.Printf("(d) %s\b\b\b\b\b%s\b\b\b\b\b%s:OK\n", farmer, farmer, farmer)
	fmt.Println("------ When sleep 1 second after three farmer are drawn ------")
	fmt.Printf("(x) %s%s%s", farmer, farmer, farmer)
	time.Sleep(time.Second)
	fmt.Printf("\b\b\b\b\b: NG ???\n")
}
