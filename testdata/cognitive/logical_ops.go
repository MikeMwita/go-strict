package cognitive

// SingleAndSequence: base=1, one && sequence=+1 → 2
func SingleAndSequence(a, b, c bool) bool {
	return a && b && c
}

// TwoSequences: base=1, && then || = two sequences = +2 → 3
func TwoSequences(a, b, c bool) bool {
	return a && b || c
}

// ParenGrouped: base=1, (a&&b) and (c||d) are separate contexts = +2 → 3
func ParenGrouped(a, b, c, d bool) bool {
	return (a && b) || (c || d)
}

// InsideIf: base=1, if(n=0)=+1, logical seq(&&)=+1, return(n=1)=+2 → 5
func InsideIf(x, y, z int) bool {
	if x > 0 && y > 0 && z > 0 { // n=0 → +1; one && sequence → +1
		return true // n=1 → +2
	}
	return false
}
