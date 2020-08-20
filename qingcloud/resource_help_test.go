package qingcloud

import (
	"testing"
)

func TestStringSliceDiff(t *testing.T) {
	nls := [][]string{
		{"hello", "world", "golang"},
		{"hello", "world", "golang"},
		{"hello", "world"},
		{"hello"},
	}
	ols := [][]string{
		{"hello", "world", "golang"},
		{"hello", "world", "java"},
		{"hello"},
		{"hello", "world"},
	}
	additionsList := [][]string{
		{},
		{"golang"},
		{"world"},
		{},
	}
	deletionsList := [][]string{
		{},
		{"java"},
		{},
		{"world"},
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
