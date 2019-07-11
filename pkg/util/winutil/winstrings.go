// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2019 Datadog, Inc.

// +build windows

package winutil

import (
	"bytes"
	"unsafe"
)

// ConvertWindowsStringList Converts a windows-style C list of strings
// (single null terminated elements
// double-null indicates the end of the list) to an array of Go strings
func ConvertWindowsStringList(winput []uint16) []string {
	var retstrings []string
	var buffer bytes.Buffer

	for i := 0; i < (len(winput) - 1); i++ {
		if winput[i] == 0 {
			retstrings = append(retstrings, buffer.String())
			buffer.Reset()

			if winput[i+1] == 0 {
				return retstrings
			}
			continue
		}
		buffer.WriteString(string(rune(winput[i])))
	}
	return retstrings
}

// ConvertWindowsString converts a windows c-string
// into a go string.  Even though the input is array
// of uint8, the underlying data is expected to be
// uint16 (unicode)
func ConvertWindowsString(winput []uint8) string {
	var retstring string
	for i := 0; i < len(winput); i += 2 {
		dbyte := (uint16(winput[i+1]) << 8) + uint16(winput[i])
		if dbyte == 0 {
			break
		}
		retstring += string(rune(dbyte))
	}
	return retstring
}

// ConvertWindowsStringPtr converts a windows c-string
// into a go string.  Even though the input is array
// of uint8, the underlying data is expected to be
// uint16 (unicode)
func ConvertWindowsStringPtr(winput uintptr) string {
	var retstring string
	for i := 0; ; i += 2 {
		lobyte := *(*byte)(unsafe.Pointer(winput))
		hibyte := *(*byte)(unsafe.Pointer(winput + 1))
		dbyte := ((uint16(hibyte) << 8) + uint16(lobyte))
		if dbyte == 0 {
			break
		}
		retstring += string(rune(dbyte))
		winput += 2
	}
	return retstring
}
