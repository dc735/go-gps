package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ahour(hour string) string {
	switch hour {
	case "00":
		return "a"
	case "01":
		return "b"
	case "02":
		return "c"
	case "03":
		return "d"
	case "04":
		return "e"
	case "05":
		return "f"
	case "06":
		return "g"
	case "07":
		return "h"
	case "08":
		return "i"
	case "09":
		return "j"
	case "10":
		return "k"
	case "11":
		return "l"
	case "12":
		return "m"
	case "13":
		return "n"
	case "14":
		return "o"
	case "15":
		return "p"
	case "16":
		return "q"
	case "17":
		return "r"
	case "18":
		return "s"
	case "19":
		return "t"
	case "20":
		return "u"
	case "21":
		return "v"
	case "22":
		return "w"
	case "23":
		return "x"
	}
	return "ee"
}
func main() {
	var site string
	var y string
	//var test bool
	flag.StringVar(&site, "site", "avln", "Site Code")
	flag.StringVar(&y, "y", "2017", "Year")
	//flag.BoolVar(&test, "test", true, "Print the ouput only")
	flag.Parse()
	basedir := "/archive/gps/nrtdata/"
	basedir = filepath.Join(basedir, y)
	fp := "/home/davec/sites.csv"
	df, err := os.Open(fp)
	if err != nil {
		log.Fatal(err)
	}
	defer df.Close()
	scanner := bufio.NewScanner(df)
	var stn []string
	for scanner.Scan() {
		stn = append(stn, strings.ToLower(scanner.Text()))
	}
	//fmt.Println(stn)
	for _, msite := range stn {
		dayfiles, _ := ioutil.ReadDir(basedir)
		for _, d := range dayfiles {
			hourfiles, _ := ioutil.ReadDir(filepath.Join(basedir, d.Name()))
			for _, h := range hourfiles {
				//dayfile := filepath.Join("", h.Name())
				//fmt.Println("Test" + dayfile)
				//fmt.Println(h.Name())
				alpha := ahour(h.Name())
				//fmt.Println(alpha)
				filename := msite + d.Name() + alpha + ".17d.Z"
				fullpath := filepath.Join(basedir, d.Name(), h.Name(), filename)
				if _, err := os.Stat(fullpath); err != nil {
					fmt.Println("file not found " + fullpath + " - Hour : " + h.Name())
			
				}
			}
		}
	}

}

/*
	if _, err := os.Stat("/archive/gps/nrtdata/2016/033/00"); os.IsNotExist(err) {
		return "Dir Not Found"
	}
	if _, err := os.Stat("/archive/gps/nrtdata/2016/033/00/auck033a.16d.Z"); err != nil {
		return "file not found"
	}
*/
