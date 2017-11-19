package qingcloud

import (
	"math/rand"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
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

func isServerBusy(err error) error {
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

type stop struct {
	error
}

//Stop Retry when fn's return value is nil , or fn's return type is stop struct
func simpleRetry(fn func() error) error {
	return retry(100, 10*time.Second, fn)
}

func getNamePointer(d *schema.ResourceData) (*string, bool) {
	var value *string = nil
	if d.HasChange("name") {
		if d.Get("name").(string) != "" {
			value = qc.String(d.Get("name").(string))
		} else {
			value = qc.String(" ")
		}
		return value, true
	}
	return value, false
}
func getDescriptionPointer(d *schema.ResourceData) (*string, bool) {
	var value *string = nil
	if d.HasChange("description") {
		if d.Get("description").(string) != "" {
			value = qc.String(d.Get("description").(string))
		} else {
			value = qc.String(" ")
		}
		return value, true
	}
	return value, false
}
