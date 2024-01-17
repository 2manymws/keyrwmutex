package keyrwmutex

import (
	"testing"
	"time"
)

const (
	callbackTimeout = 1 * time.Second
)

func newKeyMutexes() []*KeyRWMutex {
	return []*KeyRWMutex{
		New(0),
		New(1),
		New(2),
		New(4),
	}
}

// Copy from https://github.com/kubernetes/utils/blob/master/keymutex/keymutex_test.go
func Test_SingleLock_NoUnlock(t *testing.T) {
	for _, km := range newKeyMutexes() {
		// Arrange
		key := "fakeid"
		callbackCh := make(chan interface{})

		// Act
		go lockAndCallback(km, key, callbackCh)

		// Assert
		verifyCallbackHappens(t, callbackCh)
	}
}

// Copy from https://github.com/kubernetes/utils/blob/master/keymutex/keymutex_test.go
func Test_SingleLock_SingleUnlock(t *testing.T) {
	for _, km := range newKeyMutexes() {
		// Arrange
		key := "fakeid"
		callbackCh := make(chan interface{})

		// Act & Assert
		go lockAndCallback(km, key, callbackCh)
		verifyCallbackHappens(t, callbackCh)
		km.UnlockKey(key)
	}
}

// Copy from https://github.com/kubernetes/utils/blob/master/keymutex/keymutex_test.go
func Test_DoubleLock_DoubleUnlock(t *testing.T) {
	for _, km := range newKeyMutexes() {
		// Arrange
		key := "fakeid"
		callbackCh1stLock := make(chan interface{})
		callbackCh2ndLock := make(chan interface{})

		// Act & Assert
		go lockAndCallback(km, key, callbackCh1stLock)
		verifyCallbackHappens(t, callbackCh1stLock)
		go lockAndCallback(km, key, callbackCh2ndLock)
		verifyCallbackDoesntHappens(t, callbackCh2ndLock)
		km.UnlockKey(key)
		verifyCallbackHappens(t, callbackCh2ndLock)
		km.UnlockKey(key)
	}
}

func Test_SingleRLock_NoRUnlock(t *testing.T) {
	for _, km := range newKeyMutexes() {
		// Arrange
		key := "fakeid"
		callbackCh := make(chan interface{})

		// Act
		go rLockAndCallback(km, key, callbackCh)

		// Assert
		verifyCallbackHappens(t, callbackCh)
	}
}

func Test_DoubleRLock_DoubleRUnlock(t *testing.T) {
	for _, km := range newKeyMutexes() {
		// Arrange
		key := "fakeid"
		callbackCh1stLock := make(chan interface{})
		callbackCh2ndLock := make(chan interface{})

		// Act & Assert
		go rLockAndCallback(km, key, callbackCh1stLock)
		verifyCallbackHappens(t, callbackCh1stLock)
		go rLockAndCallback(km, key, callbackCh2ndLock)
		verifyCallbackHappens(t, callbackCh2ndLock)
	}
}

func Test_RLock_Lock(t *testing.T) {
	for _, km := range newKeyMutexes() {
		// Arrange
		key := "fakeid"
		callbackCh1stLock := make(chan interface{})
		callbackCh2ndLock := make(chan interface{})

		// Act & Assert
		go rLockAndCallback(km, key, callbackCh1stLock)
		verifyCallbackHappens(t, callbackCh1stLock)
		go lockAndCallback(km, key, callbackCh2ndLock)
		verifyCallbackDoesntHappens(t, callbackCh2ndLock)
		km.RUnlockKey(key)
		verifyCallbackHappens(t, callbackCh2ndLock)
	}
}

func Test_Lock_RLock(t *testing.T) {
	for _, km := range newKeyMutexes() {
		// Arrange
		key := "fakeid"
		callbackCh1stLock := make(chan interface{})
		callbackCh2ndLock := make(chan interface{})

		// Act & Assert
		go lockAndCallback(km, key, callbackCh1stLock)
		verifyCallbackHappens(t, callbackCh1stLock)
		go rLockAndCallback(km, key, callbackCh2ndLock)
		verifyCallbackDoesntHappens(t, callbackCh2ndLock)
		km.UnlockKey(key)
		verifyCallbackHappens(t, callbackCh2ndLock)
	}
}

func lockAndCallback(km *KeyRWMutex, id string, callbackCh chan<- interface{}) {
	km.LockKey(id)
	callbackCh <- true
}

func rLockAndCallback(km *KeyRWMutex, id string, callbackCh chan<- interface{}) {
	km.RLockKey(id)
	callbackCh <- true
}

func verifyCallbackHappens(t *testing.T, callbackCh <-chan interface{}) bool {
	t.Helper()
	select {
	case <-callbackCh:
		return true
	case <-time.After(callbackTimeout):
		t.Fatalf("Timed out waiting for callback.")
		return false
	}
}

func verifyCallbackDoesntHappens(t *testing.T, callbackCh <-chan interface{}) bool {
	t.Helper()
	select {
	case <-callbackCh:
		t.Fatalf("Unexpected callback.")
		return false
	case <-time.After(callbackTimeout):
		return true
	}
}
