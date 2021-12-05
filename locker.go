package tgbotapi

// Locker provides interface for lock/unlock realization
type Locker interface {
	Lock() bool
	Unlock()
}

type fakeLocker struct {}

func (l *fakeLocker) Lock() bool {
	return true
}

func (l *fakeLocker) Unlock() {}
