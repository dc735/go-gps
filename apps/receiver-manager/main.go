package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"./commands"

	"github.com/delta/meta"
)

var admin = "http://admin:cuseeme@"
var FwVersion = "/prog/Show?FirmwareVersion"
var FwInstall = "/prog/Upload?FirmwareFile"
var CloneInstall = "/cgi-bin/clone_fileUpload.html"
var CloneFile = "/home/davec/GNS_NETR9_5.10.xml"
var FirmwareFile = "/home/davec/NetR9_V5.10.timg"

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "make noise")
	var data string
	flag.StringVar(&data, "data", "/home/davec/go/src/github.com/delta", "base data directory")
	var site string
	flag.StringVar(&site, "site", "TEST", "base site code")
	var test bool
	flag.BoolVar(&test, "test", true, "Are we just testing")

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
	site = strings.ToUpper(site)
	// load antenna details into a map ...
	antennamap := make(map[string]meta.InstalledAntenna)
	{
		var antennas meta.InstalledAntennaList
		if err := meta.LoadList(filepath.Join(data, "install/antennas.csv"), &antennas); err != nil {
			panic(err)
		}

		for _, n := range antennas {
			antennamap[n.MarkCode] = n
		}
	}
	// load station details
	markmap := make(map[string]meta.Mark)
	{
		var marks meta.MarkList
		if err := meta.LoadList(filepath.Join(data, "network/marks.csv"), &marks); err != nil {
			panic(err)
		}

		for _, m := range marks {
			markmap[m.Code] = m
		}
	}
	///////////////////////////////////////////////////////
	// sort the akeys on output
	var akeys []string
	for k, _ := range antennamap {
		akeys = append(akeys, k)
	}
	sort.Strings(akeys)
	// a simple loop and print
	for _, k := range akeys {
		v, ok := antennamap[k]
		if !ok {
			panic("invalid mark key: " + k)
		}
		if v.MarkCode == site {
			{
				n, ok := antennamap[v.MarkCode]
				if !ok {
					panic("unable to find network: " + v.MarkCode)
				}
				j, err := json.MarshalIndent(n, "", "  ")
				if err != nil {
					panic(err)
				}
				fmt.Println(string(j))
			}
		}
	}
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
			j, err := json.MarshalIndent(v, "", "  ")
			if err != nil {
				panic(err)
			}
			fmt.Println(string(j))
		}
	}
	/////_///////////////////////////////////////////////////
	SiteIP := commands.FindIP(site)
	fmt.Println(SiteIP.String())
	fmt.Println(site)
	var Rip = "10.100.59.150"
	//#####
	if test == false {
		go func() {
			for i := 0; i < 360; i++ {
				s := commands.Status(admin + Rip)
				fmt.Println(s)
				time.Sleep(time.Second * 5)
				if i > 5 {
					if strings.Contains(s, "Done") {
						break
					}
				}
			}
			//	fmt.Println("We Gave up")
			wg.Done()
		}()
		go func() {
			c := commands.Upload(admin+Rip+FwInstall, FirmwareFile)
			fmt.Println(c)
			wg.Done()
		}()
		// Wait for the goroutines to finish.
		fmt.Println("Waiting For Upgrade To Finish")
		wg.Wait()
	}
	//	fmt.Println("Installing Clone File")
	//	c := commands.Upload(admin+Rip+CloneInstall,CloneFile)
	//  fmt.Println(c)
}

/*
	//	s := commands.FindStatus("http://admin:cuseeme@10.100.59.151")
	//	c := commands.Upload("http://admin:cuseeme@10.100.59.151/cgi-bin/clone_fileUpload.html","/home/davec/GNS_NETR9_5.10.xml")
		c := commands.Upload("http://admin:cuseeme@10.100.59.151/prog/Upload?FirmwareFile","/home/davec/NetR9_V5.10.timg")
			for i := 0; i < 100; i++ {
				s := commands.FindStatus("http://admin:cuseeme@10.100.59.151")
				fmt.Println(s)
				time.Sleep(time.Second * 1)
			}
			fmt.Println(c)
			fmt.Println(e.IP)
*/
