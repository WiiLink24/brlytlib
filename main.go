package main

import (
	"io/ioutil"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 4 {
		log.Println("Usage: brlytlib [toXML|toBRLYT] <input> <output>")
		os.Exit(1)
	}

	action := os.Args[1]
	input := os.Args[2]
	output := os.Args[3]

	switch action {
	case "toXML":
		data, err := ParseBRLYT(input)
		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile(output, data, 0666)
		if err != nil {
			panic(err)
		}
	case "toBRLYT":
		file, err := ioutil.ReadFile(input)
		if err != nil {
			return
		}

		data, err := WriteBRLYT(file)
		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile(output, data, 0666)
		if err != nil {
			panic(err)
		}
	}

}
