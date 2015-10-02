package main

import (
	"fmt"
	"log"
	"bufio"
    "os"
	"strings"
)

import (
//        "github.com/zone"
)
func check(e error) {
    if e != nil {
        panic(e)
    }
}
func main() {
	file, err := os.Open("/home/davec/go/src/github.com/go-gps/delta-ng/DELTA_MARK.csv")
	if err != nil {
	    log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result := strings.Split(scanner.Text(), ",")
		var fn = "/home/davec/go/src/github.com/go-gps/delta-ng/sites/"+strings.ToLower(result[0])+".toml"
		var f, err = os.Create(fn)
    	check(err)
		defer f.Close()
		w := bufio.NewWriter(f)
		o := fmt.Sprintf("Code = \"%s\"\nName = \"%s\"\nLatitude = %s\nLongitude =%s\nHeight = %s\nGroundRelationship = %s\nOpened = %s\nClosed = %s\nDOME = \"%s\"\nNetwork = \"%s\"\nMonument = \"%s\"\n",result[0],result[1],result[2],result[3],result[4],result[5],result[6],result[7],result[8],result[9],result[10])
		_, err = w.WriteString(o)
		check(err)
		w.Flush()
	}
	if err := scanner.Err(); err != nil {
    	log.Fatal(err)
	}
}