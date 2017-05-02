// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package slice provides a slice sorting function.
package util

import "sort"

// Sort sorts the provided slice using the function less.
// If slice is not a slice, Sort panics.
func Sort(slice interface{}, less func(i, j int) bool) {
	sort.Slice(slice, less)
}
