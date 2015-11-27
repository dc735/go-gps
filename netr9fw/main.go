package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"github.com/delta/meta"
//	"bytes"
)

func main() {

	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "make noise")

	var data string
	flag.StringVar(&data, "data", "/home/davec/go/src/github.com/delta/data", "base data directory")

	var site string
	flag.StringVar(&site, "site", "TEST", "base site code")

	
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "An example program to examine ...\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}

	flag.Parse()

	// load network details into a map ...
	antennamap := make(map[string]meta.InstalledAntenna)
	{
		var antennas meta.InstalledAntennas
		if err := meta.LoadList(filepath.Join(data, "antennas/trimble/installs.csv"), &antennas); err != nil {
			panic(err)
		}

		for _, n := range antennas {
			antennamap[n.Mark] = n
		}
	}
	//fmt.Println(netmap)

	// load station details
	markmap := make(map[string]meta.Mark)
	{
		var marks meta.Marks
		if err := meta.LoadList(filepath.Join(data, "marks.csv"), &marks); err != nil {
			panic(err)
		}

		for _, m := range marks {
			markmap[m.Code] = m
		}
	}
	//fmt.Println(netmap)

	// sort the keys on output
	var keys []string
	for k, _ := range markmap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// a simple loop and print
	for _, k := range keys {
		v, ok := markmap[k]
		if !ok {
			panic("invalid mark key: " + k)
		}
		if v.Code == site {
		{
			n, ok := antennamap[v.Code]
			if !ok {
				panic("unable to find network: " + v.Code)
			}
			j, err := json.MarshalIndent(n, "", "  ")
			if err != nil {
				panic(err)
			}
			fmt.Println(string(j))
		}
//		fmt.Println(v.Code)

			j, err := json.MarshalIndent(v, "", "  ")
			if err != nil {
				panic(err)
			}
			fmt.Println(string(j))
		}
	}

}