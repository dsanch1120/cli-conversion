package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

//Gets the units from a text file and reads them into a slice
func getUnits(fileName string) [][]string {
	units := [][]string{}
	var counter int = 0

	file, err := os.Open(fileName)

	if err != nil {
		log.Fatalf("Failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		temp := strings.Split(scanner.Text(), "\t")
		for count := 0; count < len(temp); count++ {
			units[counter] = append(units[counter], temp[count])
		}
		counter++
	}

	file.Close()

	return units
}

func convert(units [][]string) {
	fmt.Println("Choose Unit to be Converted")

	//var convertIndex1 int
	//var convertIndex2 int

	for count := 0; count < len(units); count++ {
		fmt.Println(strconv.Itoa(count + 1) + ". " + units[count][0])
	}


}

func main()  {
	//Variables to be used throughout program
	const FILENAME = "units.txt"

	convert(getUnits(FILENAME))

}
