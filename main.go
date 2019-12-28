package main

import "fmt"

func main() {
	files := WalkPath(".")
	for _, file := range files {
		fmt.Println(file)
	}

}
