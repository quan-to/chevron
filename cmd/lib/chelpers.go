package main

import "C"
import (
	"github.com/quan-to/slog"
	"unsafe"
)

const TRUE = C.int(1)
const FALSE = C.int(0)
const ERROR = C.int(-1)
const OK = TRUE

func copyStringToC(dst *C.char, src []byte, n int) {
	dstPtr := unsafe.Pointer(dst)
	rBuf := (*[1 << 30]byte)(dstPtr)
	srcLen := len(src)

	copyLen := srcLen
	if srcLen > n-1 {
		copyLen = n - 1
	}

	copy(rBuf[:], src[:copyLen])
	if srcLen < n-1 {
		for i := srcLen; i < n-1; i++ {
			rBuf[i] = 0x00
		}
	}
}

func copyFromCToGo(dst []byte, src *C.char, n int) {
	dstPtr := unsafe.Pointer(src)
	rBuf := (*[1 << 30]byte)(dstPtr)
	copy(dst, rBuf[:n])
}

func init() {
	// Disable logging for library mode
	slog.SetDebug(false)
	slog.SetInfo(false)
	slog.SetError(false)
	slog.SetWarning(false)
}

func main() {}
