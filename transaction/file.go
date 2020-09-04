package transaction

import (
	"bufio"
	"fmt"
	"os"
)

// NewFileLogger creates a new TransactionLogger that writes
// transaction log into a file.
func NewFileLogger(filename string) (Logger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("can't open transaction log file: %w", err)
	}
	// NOTE: there's no pointer type in the signature but a pointer is returned.
	// This is because Go doesn't allow pointers to interface types.
	return &FileTransactionLogger{file: file}, nil
}

// FileTransactionLogger implements TransactionLogger inferface by logging into
// a filesystem file.
type FileTransactionLogger struct {
	events       chan<- Event // write-only channel for sending events
	errors       <-chan error // read-only channel for receiving errors
	lastSequence uint64       // the last used event sequence number
	file         *os.File     // the location of the transaction log file
}

// WritePut sends an event down the events channel.
func (l *FileTransactionLogger) WritePut(key, value string) {
	l.events <- Event{EventType: EventPut, Key: key, Value: value}
}

// WriteDelete sends an event down the events channel.
func (l *FileTransactionLogger) WriteDelete(key string) {
	l.events <- Event{EventType: EventDelete, Key: key}
}

// Err returns an error channel.
func (l *FileTransactionLogger) Err() <-chan error {
	return l.errors
}

// Run writes transaction events.
func (l *FileTransactionLogger) Run() {
	// Buffered channel handles bursts of events w/o being slowed by disk IO.
	events := make(chan Event, 16)
	l.events = events

	// Buffer value of 1 allows us to send an error in a non-blocking manner.
	errors := make(chan error, 1)
	l.errors = errors

	go func() {
		for e := range events { // retrieve the next event
			l.lastSequence++

			_, err := fmt.Fprintf( // write the event to the log
				l.file,
				"%d\t%d\t%s\t%s\n",
				l.lastSequence, e.EventType, e.Key, e.Value)
			if err != nil {
				errors <- err
				return
			}
		}
	}()
}

// ReadEvents reads transaction events.
func (l *FileTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	scanner := bufio.NewScanner(l.file)
	outEvent := make(chan Event)
	outError := make(chan error, 1)

	go func() {
		var e Event

		defer close(outEvent)
		defer close(outError)

		for scanner.Scan() {
			line := scanner.Text()

			fmt.Sscanf(
				line, "%d\t%d\t%s\t%s",
				&e.Sequence, &e.EventType, &e.Key, &e.Value)

			// Sanity check! Are the sequence numbers in increasing order?
			if l.lastSequence >= e.Sequence {
				outError <- fmt.Errorf(
					"transaction numbers out of sequence (last %v, next %v)",
					l.lastSequence, e.Sequence)
				return
			}

			l.lastSequence = e.Sequence // update last used sequence #

			outEvent <- e // send the event along
		}

		if err := scanner.Err(); err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
		}
	}()

	return outEvent, outError
}
