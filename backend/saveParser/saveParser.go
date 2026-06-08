package saveparser

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L${SRCDIR}/lib -Wl,-rpath,${SRCDIR}/lib -luesave_go_bridge
#include "uesave_bridge.h"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func ConvertUesaveToJSON(fileBytes []byte) (string, error) {
	// Pass Go byte slice pointer to Rust
	cBytes := (*C.uint8_t)(unsafe.Pointer(&fileBytes[0]))
	cLen := C.size_t(len(fileBytes))

	// Call Rust function
	cJsonStr := C.convert_to_json(cBytes, cLen)
	if cJsonStr == nil {
		return "", fmt.Errorf("failed to parse save file or generate JSON")
	}
	// Defer freeing the memory allocated by Rust
	defer C.free_rust_string(cJsonStr)

	// Convert C string back to Go string
	goJsonStr := C.GoString(cJsonStr)
	return goJsonStr, nil
}
