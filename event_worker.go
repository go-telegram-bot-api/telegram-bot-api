package tgbotapi

// EventWorker provides interface to keep event ID up to date
type EventWorker interface {
	UpdateEventID(eventID int) error
	GetLastEventID() (int, error)
}

type fakeEventWorker struct {
	lastEventID int
}

func (w *fakeEventWorker) UpdateEventID(eventID int) error {
	w.lastEventID = eventID
	return nil
}

func (w *fakeEventWorker) GetLastEventID() (int, error) {
	return w.lastEventID, nil
}
