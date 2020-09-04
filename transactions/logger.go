// Package transactions adds persistence to key/value store by logging the
// trasactions.
package transactions

// TransactionLogger is the interface that groups methods for logging
// transactions.
type TransactionLogger interface {
	WriteDelete(key string)
	WritePut(key, value string)
	Err() <-chan error
	Run()
	ReadEvents() (<-chan Event, <-chan error)
}

// EventType represents a transaction event type.
type EventType byte

// Types of transaction events.
const (
	_                     = iota // iota == 0; ignore the zero valus
	EventDelete EventType = iota // iota == 1
	EventPut                     // iota == 2; implicitly repeat
)

// Event represents a transaction event.
type Event struct {
	Sequence  uint64    // a unique record ID
	EventType EventType // action taken
	Key       string    // the key affected by this transaction
	Value     string    // the value of a PUT transaction
}
