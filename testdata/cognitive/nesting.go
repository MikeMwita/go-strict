// Package cognitive contains testdata demonstrating cognitive complexity rules.
package cognitive

import "fmt"

// DeepNesting demonstrates deeply nested control flow.
// Expected complexity: base=1, if×4 at n=0,1,2,3, return at n=4 → 1+1+2+3+4+5=16
func DeepNesting(a, b, c, d bool) {
	if a { // n=0 → +1
		if b { // n=1 → +2
			if c { // n=2 → +3
				if d { // n=3 → +4
					fmt.Println("deep") // return at n=4 → +5 (if it were return)
				}
			}
		}
	}
}

// EarlyReturn shows the benefit of guard clauses (lower complexity).
// Expected complexity: base=1, if at n=0 → +1, branch at n=1 → +2 → total=4
func EarlyReturn(x int) int {
	if x <= 0 { // n=0 → +1
		return 0 // n=1 → +2
	}
	return x * 2
}

// NestedLoops shows range-in-range nesting.
// Expected: base=1, range(n=0)=+1, range(n=1)=+2, if(n=2)=+3, branch(n=3)=+4 → 11
func NestedLoops(matrix [][]int) int {
	sum := 0
	for _, row := range matrix { // n=0 → +1
		for _, val := range row { // n=1 → +2
			if val > 0 { // n=2 → +3
				sum += val
				break // n=3 → +4
			}
		}
	}
	return sum
}
