package qingcloud

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/yunify/qingcloud-sdk-go/logger"
	"github.com/yunify/qingcloud-sdk-go/request/errors"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func stringSliceDiff(nl, ol []string) ([]string, []string) {
	var additions []string
	var deletions []string
	for i := 0; i < 2; i++ {
		for _, n := range nl {
			found := false
			for _, o := range ol {
				if n == o {
					found = true
					break
				}
			}
			if !found {
				if i == 0 {
					additions = append(additions, n)
				} else {
					deletions = append(deletions, n)
				}
			}
		}
		if i == 0 {
			nl, ol = ol, nl
		}
	}
	return additions, deletions
}

func retry(attempts int, sleep time.Duration, fn func() error) error {
	if err := fn(); err != nil {
		if s, ok := err.(stop); ok {
			// Return the original error for later checking
			return s.error
		}

		if attempts--; attempts > 0 {
			// Add some randomness to prevent creating a Thundering Herd
			jitter := time.Duration(rand.Int63n(int64(sleep)))
			sleep = sleep + jitter/2

			time.Sleep(sleep)
			logger.Warn("Retry function")
			return retry(attempts, 2*sleep, fn)
		}
		return err
	}

	return nil
}
func retryServerBusy(f func() error) error {
	wraaper := func() error {
		err := f()
		if err != nil {
			if err, ok := err.(errors.QingCloudError); ok {
				if err.RetCode == SERVERBUSY {
					return err
				}
				return stop{err}

			}
			return stop{err}
		}
		return nil
	}
	return simpleRetry(wraaper)
}

type stop struct {
	error
}

//Stop Retry when fn's return value is nil , or fn's return type is stop struct
func simpleRetry(fn func() error) error {
	return retry(100, 10*time.Second, fn)
}

func getQingCloudErr(opration string, retCode *int, message *string, err error) error {
	if err != nil {
		return fmt.Errorf("Error %s : %s ", opration, err)
	}
	if retCode != nil && qc.IntValue(retCode) != 0 {
		return fmt.Errorf("Error %s : %s ", opration, *message)
	}
	return nil
}
