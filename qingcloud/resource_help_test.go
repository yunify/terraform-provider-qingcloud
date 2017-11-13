package qingcloud

import (
	"testing"
)

func TestStringSliceDiff(t *testing.T) {
	nls := [][]string{
		[]string{"hello", "world", "golang"},
		[]string{"hello", "world", "golang"},
		[]string{"hello", "world"},
		[]string{"hello"},
	}
	ols := [][]string{
		[]string{"hello", "world", "golang"},
		[]string{"hello", "world", "java"},
		[]string{"hello"},
		[]string{"hello", "world"},
	}
	additionsList := [][]string{
		[]string{},
		[]string{"golang"},
		[]string{"world"},
		[]string{},
	}
	deletionsList := [][]string{
		[]string{},
		[]string{"java"},
		[]string{},
		[]string{"world"},
	}
	for k, nl := range nls {
		additions, deletions := stringSliceDiff(nl, ols[k])
		a, b := stringSliceDiff(additions, additionsList[k])
		c, d := stringSliceDiff(deletions, deletionsList[k])
		if len(a) != 0 || len(b) != 0 || len(c) != 0 || len(d) != 0 {
			t.Errorf("test case %d want additions: %+v deletions: %+v, got additions: %+v deletions: %+v", k, additionsList[k], deletionsList[k], additions, deletions)
		}
	}
}

func TestIsServerBusy(t *testing.T) {
	exceptValues := []int{5200, 132, 356}
	validValues := []int{5100}
	for _, v := range validValues {
		result := IsServerBusy(v)
		if !result {
			t.Fatalf("%q should be a vaild value in %#v", v, validValues)
		}
	}
	for _, v := range exceptValues {
		result := IsServerBusy(v)
		if result {
			t.Fatalf("%q should be a except value in %#v", v, exceptValues)
		}
	}

}
