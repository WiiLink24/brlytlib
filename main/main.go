package main

import (
	"brlytlib"
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

	brlyt, err := brlytlib.ParseBRLYT(input)
	if err != nil {
		panic(err)
		return
	}
	if err != nil {
		panic(err)
		log.Fatal(err)
		return
	}

	err = ioutil.WriteFile(output, brlyt, 0600)
	if err != nil {
		panic(err)
	}
}
