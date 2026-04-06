package cognitive

import "sort"

// SortWithClosure shows a closure adding to nesting.
// base=1, closure(n=0)=+1, if inside closure(n=1)=+2, branch(n=2)=+3 → 7
func SortWithClosure(items []int) []int {
	sorted := make([]int, len(items))
	copy(sorted, items)
	sort.Slice(sorted, func(i, j int) bool { // closure at n=0 → +1
		if sorted[i] == sorted[j] { // n=1 → +2
			return false // n=2 → +3
		}
		return sorted[i] < sorted[j]
	})
	return sorted
}

// GoroutineWithClosure shows a go statement with a closure.
// base=1, go(n=0)=+1, range inside goroutine(n=1)=+2 → 4
func GoroutineWithClosure(items []int, out chan int) {
	go func() { // n=0 → +1
		for _, v := range items { // n=1 → +2
			out <- v
		}
	}()
}

// DeferWithClosure shows defer with a closure.
// base=1, defer(n=0)=+1, if inside defer(n=1)=+2 → 4
func DeferWithClosure(r interface{ Close() error }) {
	defer func() { // n=0 → +1
		if err := r.Close(); err != nil { // n=1 → +2
			_ = err
		}
	}()
}
