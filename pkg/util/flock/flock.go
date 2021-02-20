package flock

import (
	"fmt"
	"syscall"
	"time"
)

const (
	perm = 0666
)

type Flocker struct {
	fd int
}

func WithLock(filename string, closure func() error) error {
	lock, err := New(filename)
	if err != nil {
		return err
	}

	defer lock.Close()

	if err := lock.Lock(); err != nil {
		return err
	}

	return closure()
}

func New(filename string) (*Flocker, error) {
	fd, err := syscall.Open(filename, syscall.O_CREAT|syscall.O_RDONLY, perm)
	if err != nil {
		return nil, err
	}
	return &Flocker{fd: fd}, nil
}

func (m *Flocker) Lock() error {
	return syscall.Flock(m.fd, syscall.LOCK_EX)
}

func (m *Flocker) tryLock() (bool, error) {
	err := syscall.Flock(m.fd, syscall.LOCK_EX|syscall.LOCK_NB)
	if err == nil {
		return true, nil
	}
	switch err {
	case syscall.EWOULDBLOCK:
		return false, nil
	default:
		return false, err
	}
}

func (m *Flocker) LockWithTimeout(totalWait time.Duration, waitIteration time.Duration) error {
	iterations := int(totalWait/waitIteration) + 1
	for i := 0; i < iterations; i++ {
		if i > 0 {
			time.Sleep(waitIteration)
		}
		if haveLock, err := m.tryLock(); haveLock {
			return nil // success
		} else if err != nil {
			return err
		}
	}
	return fmt.Errorf("could not get lock within %+v", totalWait)
}

func (m *Flocker) Unlock() error {
	return syscall.Flock(m.fd, syscall.LOCK_UN)
}

// Close unlocks the lock and closes the underlying file descriptor.
func (m *Flocker) Close() error {
	m.Unlock() // ignore err here
	return syscall.Close(m.fd)
}
