package main

import (
	"encoding/json"
	"flag"
	"fmt"
)

func main() {
	// I couldn't figure out a nice and clean way for figuring out automatically
	// what the volume name of the current directory actually was -- just ask
	// the user
	drive := flag.String("drive", "NA", "REQUIRED: The name of the current drive being searched.")
	rootDir := flag.String("root", ".", "The name of the directory being searched.")
	flag.Parse()
	if *drive == "NA" {
		fmt.Println("Error: Drive name required")
		return
	}

	files := WalkPath(*rootDir, *drive)
	for _, file := range files {
		data, _ := json.Marshal(file)
		fmt.Println(string(data))
	}

}
