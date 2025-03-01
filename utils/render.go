package utils

import "fmt"

func RenderClear() {
	// TODO: Test if this works on both Linux and Win
	fmt.Printf(string("\x1b[2J\x1b[H"))
	// if runtime.GOOS == "windows" {
	// 	cmd := exec.Command("cmd", "/c", "cls")
	// 	cmd.Stdout = os.Stdout
	// 	_ = cmd.Run()
	// 	return
	// }
	//
	// cmd := exec.Command("clear")
	// cmd.Stdout = os.Stdout
	// _ = cmd.Run()
}
