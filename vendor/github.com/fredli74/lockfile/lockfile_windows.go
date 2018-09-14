// Copyright 2015 Fredrik Lidström. All rights reserved.
// Use of this source code is governed by the standard MIT License (MIT)
// that can be found in the LICENSE file.

// +build windows

package lockfile

import (
	"os"
)

// ProcessRunning is a cross-platform check to work on both Windows and Unix systems as the os.FindProcess() function works differently.
func ProcessRunning(pid int) bool {
	_, err := os.FindProcess(pid)
	return err == nil
}
