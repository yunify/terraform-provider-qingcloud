package qingcloud

import (
	"math/rand"
	"time"

	"github.com/yunify/qingcloud-sdk-go/logger"
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
func IsServerBusy(RetCode int) bool {
	return RetCode == SERVERBUSY
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

type stop struct {
	error
}

func simpleRetry(fn func() error) error {
	return retry(100, 10*time.Second, fn)
}
