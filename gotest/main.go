package main

import "testing"

func TestSum(t *testing.T) {
	total := Sum(5, 4)
	if total != 10 {
		t.Errorf("Sum was incorrect, got: %d, want: %d.", total, 10)
	}
}
func Sum(x int, y int) int {
	return x + y
}

func main() {
	Sum(5, 5)
}
