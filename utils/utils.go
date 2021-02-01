package utils

import "sort"

func SetToSlice(set map[string]struct{}) []string {
	res := make([]string, 0, len(set))
	for k := range set {
		res = append(res, k)
	}
	sort.Strings(res)
	return res
}
