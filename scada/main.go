package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	//	"net"
	"net/http"
	"strconv"
	//"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type stnMtr struct {
	//	latency time.Duration
	volts float64
}

func setr(ip string, val string) string {
	response, err := http.Get("http://" + ip + ":8080/cmd.shtm?cmd=setreg&reg=10=10&val=" + val)
	if err != nil {
		return "Not set"
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return "Not Read"
		}
		//result := strings.Split(string(contents), "\n")
		//		fmt.Println(string(contents))
		x := string(contents)
		return x
	}
	return "Fin"
}
func gv(ip, i string) float64 {
	//	var ok bool
	response, err := http.Get("http://" + ip + ":8080/cmd.shtm?cmd=getreg&reg=" + i)
	if err != nil {
		return 0.00 //"Not Found"
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return 0.00 //"volts not read"
		}
		//		result := strings.Split(string(contents), "\n")
		fmt.Println(string(contents))
		y, err := strconv.ParseFloat(string(contents), 64)
		if err != nil {
			return 1.5
		}
		return y

		//return strconv.ParseFloat(string(contents),32)
	}
	return 0.00
}

var ok bool

func main() {
	var ip []string
	file, err := os.Open("./iplist")
	if err != nil {
		log.Println("File not found")
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ip = append(ip, scanner.Text())
	}
	for _, y := range ip {
		for i := 1; i <= 11; i++ {
			c := strconv.Itoa(i)
			x := gv(y, c)
			fmt.Println(c, " - ", x)
		}
		//	set := setr(y, "0")
		//	fmt.Println(set)
	}
}
