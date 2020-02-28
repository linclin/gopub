package grpool

import "time"

const (
	timeoutErrorString = "withtimeout: Operation timed out"
)

type timeoutError struct{}

func (timeoutError) Error() string {
	return timeoutErrorString
}
func Do(timeout time.Duration, fn func() (interface{}, error)) (result interface{}, timedOut bool, err error) {
	resultCh := make(chan *resultWithError, 1)

	go func() {
		result, err := fn()
		resultCh <- &resultWithError{result, err}
	}()

	select {
	case <-time.After(timeout):
		return nil, true, timeoutError{}
	case rwe := <-resultCh:
		return rwe.result, false, rwe.err
	}
}

type resultWithError struct {
	result interface{}
	err    error
}
