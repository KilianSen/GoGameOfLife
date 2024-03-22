package main

import (
	"GameOfLife/src/GUI"
	"runtime"
)

func main() {
	runtime.LockOSThread()
	GUI.GUI()
}
