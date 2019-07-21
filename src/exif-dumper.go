package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(
			os.Stderr,
			"error: bad usage.\n" +
			"usage: prog filename\n")
		os.Exit(1)
	}

	res := getRecords(os.Args[1])

	fmt.Printf(
		"DateTime: '%s'\n" +
		"DateTimeOriginal: '%s'\n" +
		"DateTimeDigitized: '%s'\n" +
		"Make: '%s'\n" +
		"Model: '%s'\n",
		res.DateTime,
		res.DateTimeOriginal,
		res.DateTimeDigitized,
		res.Make,
		res.Model)
}
