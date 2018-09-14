# lockfile #

Package lockfile implements a simple automatic lockfile (PID) method for golang.

[![Build Status](https://semaphoreci.com/api/v1/fredli74/lockfile/branches/master/badge.svg)](https://semaphoreci.com/fredli74/lockfile)

## Usage

    import "github.com/fredli74/lockfile"

### Example
```go
    if lock, err := lockfile.Lock(filename); err != nil {
        panic(err)
    } else {
        defer lock.Unlock()
    }
```

### Reference

#### func  Lock

```go
func Lock(name string) (*LockFile, error)
```
Lock automatically checks if the file already exists, if so, reads the process
ID from the file and checks if the process is running. If the process is running
a nil lock is returned and an error stating "Locked by other process".

#### func (*LockFile) Unlock

```go
func (l *LockFile) Unlock()
```
Unlock closes and deletes the lock file previously created by Lock()

#### func  ProcessRunning

```go
func ProcessRunning(pid int) bool
```
ProcessRunning is a cross-platform check to work on both Windows and Unix
systems as the os.FindProcess() function works differently.
