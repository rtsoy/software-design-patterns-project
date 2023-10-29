package telebot

import "sync"

// lastSentMessage is a struct to store the last sent message for each user.
type lastSentMessage struct {
	store map[int64]string
	mu    *sync.Mutex
}

// newLastSentMessage creates a new instance of lastSentMessage.
func newLastSentMessage() *lastSentMessage {
	return &lastSentMessage{
		store: make(map[int64]string),
		mu:    &sync.Mutex{},
	}
}

// get retrieves the last sent message for a specific user.
func (l *lastSentMessage) get(id int64) string {
	l.mu.Lock()
	defer l.mu.Unlock()

	msg := l.store[id]

	return msg
}

// set stores the last sent message for a specific user.
func (l *lastSentMessage) set(id int64, msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.store[id] = msg
}
