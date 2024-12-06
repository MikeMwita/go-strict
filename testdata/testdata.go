package main

import "fmt"

func main() {
	fmt.Println("Hello, world!")

	if true {
		fmt.Println("This condition is always true.")
		if false {
			fmt.Println("This block adds to complexity.")
		} else if true {
			fmt.Println("Another condition.")
		}
	}

	for i := 0; i < 5; i++ {
		fmt.Println("Iteration:", i)
		if i > 2 {
			fmt.Println("Complex loop.")
		}
	}

	switch x := 2; x {
	case 1:
		fmt.Println("Case 1")
	case 2:
		fmt.Println("Case 2")
		if true {
			fmt.Println("Nested in switch")
		}
	default:
		fmt.Println("Default case")
	}
}

func complexFunction() {
	if true {
		fmt.Println("Nested condition")
		if false {
			fmt.Println("This block adds to complexity.")
		}
	}

	for i := 0; i < 3; i++ {
		fmt.Println("Nested loop:", i)
		if i == 2 {
			fmt.Println("Complex condition inside loop.")
		}
	}
}
