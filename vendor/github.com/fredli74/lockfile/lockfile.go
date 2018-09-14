// Copyright 2015 Fredrik Lidström. All rights reserved.
// Use of this source code is governed by the standard MIT License (MIT)
// that can be found in the LICENSE file.

// Package lockfile implements a simple automatic lockfile (PID) method for golang.
//		if lock, err := lockfile.Lock(filename); err != nil {
//			panic(err)
//		} else {
//			defer lock.Unlock()
//		}
package lockfile

import (
	"errors"
	"fmt"
	"os"
)

var AlreadyLocked = errors.New("Locked by other process")

type LockFile struct {
	name string
	file *os.File
}

// Lock automatically checks if the file already exists, if so, reads the process ID
// from the file and checks if the process is running. If the process is running a nil
// lock is returned and an error stating "Locked by other process".
func Lock(name string) (*LockFile, error) {
	var err error

	lock := LockFile{name: name}

	if lock.file, err = os.OpenFile(lock.name, os.O_CREATE|os.O_RDWR, os.ModeTemporary|0640); err == nil {
		var pid int
		if _, err = fmt.Fscanf(lock.file, "%d\n", &pid); err == nil {
			if pid != os.Getpid() {
				if ProcessRunning(pid) {
					return nil, AlreadyLocked
				}
			}
		}

		lock.file.Seek(0, 0)
		if n, err := fmt.Fprintf(lock.file, "%d\n", os.Getpid()); err == nil {
			lock.file.Truncate(int64(n))
			return &lock, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

// Unlock closes and deletes the lock file previously created by Lock()
func (l *LockFile) Unlock() {
	l.file.Close()
	os.Remove(l.name)
}
