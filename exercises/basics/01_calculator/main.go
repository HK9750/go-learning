package main

import (
	"fmt"
	"os"
	"strconv"
)

func add(val1 float64, val2 float64) float64 {
	return val1 + val2
}

func sub(val1 float64, val2 float64) float64 {
	return val1 - val2
}

func mul(val1 float64, val2 float64) float64 {
	return val1 * val2
}

func div(val1 float64, val2 float64) float64 {
	if val2 == 0 {
		fmt.Println("Cannot divide by zero")
		return 0
	}
	return val1 / val2
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Please add the arguments")
		os.Exit(1)
	}

	operation := os.Args[1]
	val1,err1 := strconv.ParseFloat(os.Args[2], 64)
	val2,err2 := strconv.ParseFloat(os.Args[3], 64)

	if err1 != nil || err2 != nil {
		fmt.Println("Please add valid arguments")
	}

	switch operation {
	case "add":
		fmt.Println(add(val1, val2))
	case "sub":
		fmt.Println(sub(val1, val2))
	case "mul":
		fmt.Println(mul(val1, val2))
	case "div":
		fmt.Println(div(val1, val2))
	default:
		fmt.Println("Invalid operation")
	}
}