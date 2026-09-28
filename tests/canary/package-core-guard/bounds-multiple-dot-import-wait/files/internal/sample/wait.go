package sample

import (
	. "github.com/gibbonmi/bench/internal/bounds"
	. "time"
)

func unclassifiedWait() {
	_ = FixedWindow
	Sleep(Duration(17))
}
