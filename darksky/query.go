package darksky

// Where ranges a slice of DataPoint and filters the results to the items that pass
// the provided compare function
func Where(dps []DataPoint, compare func(DataPoint) bool) []DataPoint {
	res := make([]DataPoint, 0, len(dps))
	for _, d := range dps {
		if compare(d) {
			res = append(res, d)
		}
	}
	return res
}
