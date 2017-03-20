package qingcloud

import (
	"fmt"
)

func withinArrayString(limits ...string) func(v interface{}, k string) (ws []string, errors []error) {
	var limitsMap = make(map[string]bool)
	for _, v := range limits {
		limitsMap[v] = true
	}

	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(string)
		if limitsMap[value] {
			return
		}
		errors = append(errors, fmt.Errorf("%q (%q) doesn't match", k, value))
		return
	}
}

func withinArrayInt(limits ...int) func(v interface{}, k string) (ws []string, errors []error) {
	var limitsMap = make(map[int]bool)
	for _, v := range limits {
		limitsMap[v] = true
	}

	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(int)
		if limitsMap[value] {
			return
		}
		errors = append(errors, fmt.Errorf("%q (%q) doesn't match", k, value))
		return
	}
}

func withinArrayIntRange(begin, end int) func(v interface{}, k string) (ws []string, errors []error) {

	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(int)
		if value >= begin && value <= end {
			return
		}
		errors = append(errors, fmt.Errorf("%q (%q) should > %d  && < %d ", k, value, begin, end))
		return
	}
}

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

