package utils

import "fmt"

func RenderClear() {
	fmt.Print("\x1b[2J\x1b[H")
}
