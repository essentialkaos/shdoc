//go:build gofuzz
// +build gofuzz

package parser

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2025 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bytes"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func Fuzz(data []byte) int {
	_, errs := readData("temp", bytes.NewReader(data))

	if len(errs) != 0 {
		return 0
	}

	return 1
}
