package main

import (
	"encoding/xml"
	brlyt "github.com/WiiLink24/brlytlib"
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
		file, err := os.ReadFile(input)
		if err != nil {
			log.Fatalln(err)
		}

		root, err := brlyt.ParseBRLYT(file)
		if err != nil {
			log.Fatalln(err)
		}

		theXML, err := xml.MarshalIndent(root, "", "\t")
		if err != nil {
			log.Fatalln(err)
		}

		err = os.WriteFile(output, theXML, 0666)
		if err != nil {
			log.Fatalln(err)
		}
	case "toBRLYT":
		file, err := os.ReadFile(input)
		if err != nil {
			log.Fatalln(err)
		}

		data, err := brlyt.WriteBRLYT(file)
		if err != nil {
			log.Fatalln(err)
		}

		err = os.WriteFile(output, data, 0666)
		if err != nil {
			log.Fatalln(err)
		}
	default:
		log.Println("Usage: brlytlib [toXML|toBRLYT] <input> <output>")
		os.Exit(1)
	}

}
