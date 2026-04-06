package cognitive

// Map applies fn to every element of xs and returns the results.
// With the generics rule enabled: base=1, generic=+1, range(n=0)=+1 → 3
func Map[T, U any](xs []T, fn func(T) U) []U {
	out := make([]U, len(xs))
	for i, x := range xs { // n=0 → +1
		out[i] = fn(x)
	}
	return out
}

// Filter returns elements of xs for which keep returns true.
// generic=+1, range(n=0)=+1, if(n=1)=+2 → base+4=5
func Filter[T any](xs []T, keep func(T) bool) []T {
	var out []T
	for _, x := range xs { // n=0 → +1
		if keep(x) { // n=1 → +2
			out = append(out, x)
		}
	}
	return out
}
