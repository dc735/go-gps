package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type Access struct {
	Code       string   `json:"code"`
	AccessDesc string   `json:"accessdesc"`
	Company    string   `json:"company"`
	Address    *Address `json:"address,omitempty"`
}

type Address struct {
	Addr1 string `json:"addr1"`
	Addr2 string `json:"addr2"`
	Addr3 string `json:"addr3"`
	City  string `json:"city"`
	Notes string `json:"notes"`
}

func main() {
	var aList []Access
	file, err := os.Open("/home/davec/git/sit/data/access.csv")
	if err != nil {
		log.Println("File not found")
		log.Fatal(err)
	}
	defer file.Close()
	reader := csv.NewReader(bufio.NewReader(file))
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		aList = append(aList, Access{
			Code:       line[0],
			AccessDesc: line[1],
			Company:    line[2],
			Address: &Address{
				Addr1: line[3],
				Addr2: line[4],
				Addr3: line[5],
				City:  line[6],
				Notes: line[7],
			},
		})
		accessJason, _ := json.Marshal(aList)
		fmt.Println(string(accessJason))
	}
	a := Access{}
	code := a.Code
	if code == "CORM" {
		fmt.Println(code)
	}
}
