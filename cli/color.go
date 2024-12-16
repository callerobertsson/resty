// Package cli color functions.
package cli

// Terminal color constants.
var (
	NORM     = "\033[0m"
	TITLE    = "\033[1;33m" // title - yellow bold
	SUBTITLE = "\033[37m"   // greyish
	SELECTED = "\033[1;32m" // selected - dark green
	NOTICE   = "\033[1;31m" // red, bold
)

// Turn colors on.
func colorOn() {
	NORM = "\033[0m"
	TITLE = "\033[1;33m"    // title - yellow bold
	SUBTITLE = "\033[37m"   // greyish
	SELECTED = "\033[1;32m" // selected - dark green
	NOTICE = "\033[1;31m"   // red, bold
}

// Turn color off.
func colorOff() {
	NORM = ""
	TITLE = ""    // title - yellow bold
	SUBTITLE = "" // greyish
	SELECTED = "" // selected - dark green
	NOTICE = ""   // red, bold
}
