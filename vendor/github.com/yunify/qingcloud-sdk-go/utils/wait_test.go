package utils

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestWaitForSpecificOrError(t *testing.T) {

	waitInterval := 100 * time.Millisecond
	timeout := 10*waitInterval + waitInterval/2
	times := 0
	err := WaitForSpecificOrError(func() (bool, error) {
		times++
		println("times", times)
		if times == 3 {
			return true, nil
		}
		return false, nil
	}, timeout, waitInterval)
	assert.NoError(t, err)
	assert.Equal(t, 3, times)

	times = 0
	err = WaitForSpecificOrError(func() (bool, error) {
		times++
		println("times", times)
		if times == 3 {
			return false, errors.New("error")
		}
		return false, nil
	}, timeout, waitInterval)
	assert.Error(t, err)
	assert.Equal(t, 3, times)

	times = 0
	err = WaitForSpecificOrError(func() (bool, error) {
		times++
		println("times", times)
		return false, nil
	}, timeout, waitInterval)
	assert.Error(t, err)
	tErr, ok := err.(*TimeoutError)
	assert.True(t, ok)
	assert.Equal(t, timeout, tErr.timeout)
	assert.Equal(t, 10, times)
}
