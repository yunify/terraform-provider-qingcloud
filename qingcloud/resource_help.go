package qingcloud

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
