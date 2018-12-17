package remote_signer

func stringIndexOf(v string, a []string) int {
	for i, vo := range a {
		if vo == v {
			return i
		}
	}

	return -1
}
