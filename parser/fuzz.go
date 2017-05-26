// +build gofuzz

package parser

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2017 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
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
