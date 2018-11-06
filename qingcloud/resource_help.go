package qingcloud

import (
	"encoding/base64"
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

func getNamePointer(d *schema.ResourceData) (value *string, update bool) {
	update = false
	if d.HasChange(resourceName) {
		if d.Get(resourceName).(string) != "" {
			value = qc.String(d.Get(resourceName).(string))
		} else {
			value = qc.String(" ")
		}
		update = !d.IsNewResource()
	}
	return value, update
}
func getDescriptionPointer(d *schema.ResourceData) (*string, bool) {
	var value *string = nil
	if d.HasChange(resourceDescription) {
		if d.Get(resourceDescription).(string) != "" {
			value = qc.String(d.Get(resourceDescription).(string))
		} else {
			value = qc.String(" ")
		}
		return value, true
	}
	return value, false
}
func getSetStringPointer(d *schema.ResourceData, key string) *string {
	if d.Get(key).(string) != "" {
		return qc.String(d.Get(key).(string))
	}
	return nil
}

func getUpdateStringPointer(d *schema.ResourceData, key string) *string {
	if d.Get(key).(string) != "" {
		return qc.String(d.Get(key).(string))
	}
	return qc.String(" ")
}

func isBase64Encoded(data []byte) bool {
	_, err := base64.StdEncoding.DecodeString(string(data))
	return err == nil
}

func getUpdateStringPointerInfo(d *schema.ResourceData, key string) (value *string, update bool) {
	update = false
	if d.HasChange(key) {
		if d.Get(key).(string) != "" {
			value = qc.String(d.Get(key).(string))
		} else {
			value = qc.String(" ")
		}
		update = !d.IsNewResource()
	}
	return
}

func getUpdateIntPointerInfo(d *schema.ResourceData, key string) (value *int, update bool) {
	update = false
	if d.HasChange(key) {
		value = qc.Int(d.Get(key).(int))
		update = !d.IsNewResource()
	}
	return
}
