package main

import (
	"io/ioutil"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		log.Println("Usage: brlytlib <input> <output>")
		os.Exit(1)
	}

	input := os.Args[1]
	output := os.Args[2]

	data, err := ParseBRLYT(input)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(output, data, 0666)
	if err != nil {
		panic(err)
	}
}
