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
		errors = append(errors, fmt.Errorf("%q (%q) doesn't match  %q", k, value))
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
		errors = append(errors, fmt.Errorf("%q (%q) doesn't match ", k, value))
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
