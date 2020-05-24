package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

var clear map[string]func()

func init() {
	clear = make(map[string]func())
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	}
}

//Gets the units from a text file and reads them into a slice
func getUnits(fileName string) (map[string][]string, []string) {
	CallClear()

	conversion := make(map[string][]string)
	units := []string{}

	var counter = 0

	file, err := os.Open(fileName)

	if err != nil {
		log.Fatalf("Failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		temp := strings.Split(scanner.Text(), "\t")
		conversion[temp[0]] = temp[1:]
		units = append(units, temp[0])
		counter++
	}

	file.Close()

	return conversion, units
}

func convert(conversion map[string][]string, units []string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Choose Unit to be Converted")

	for count := 0; count < len(units); count++ {
		fmt.Println(strconv.Itoa(count + 1) + ". " + units[count])
	}

	var ind int
	var ind2 int
	var err error

	for {
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		ind, err = strconv.Atoi(text)

		if err == nil && (0 < ind && ind <= len(units)) {
			break
		}
	}

	CallClear()

	fmt.Println("Choose Unit to Convert " + units[ind - 1] + " to")
	for count := 0; count < len(units); count++ {
		if (count != ind - 1) {
			fmt.Println(strconv.Itoa(count + 1) + ". " + units[count])
		}
	}

	for {
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		ind2, err = strconv.Atoi(text)

		if err == nil && ind2 != ind && (0 < ind2 && ind2 <= len(units)) {
			break
		}
	}

	CallClear()

	fmt.Println("How many " + units[ind - 1] + " to convert? (input \"quit\" to exit)")

	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)

	for {

		if text == "quit" || text == "Quit" {
			os.Exit(0)
		}

		f, err := strconv.ParseFloat(text, 64)

		fS := fmt.Sprintf("%g", f)

		if err != nil {
			continue
		}

		iC, err := strconv.ParseFloat(conversion[units[ind - 1]][1], 64)
		sC, err := strconv.ParseFloat(conversion[units[ind2 - 1]][0], 64)
		cF, _ := strconv.ParseFloat(conversion[units[ind2 - 1]][2], 64)

		fC := f * iC * sC
		fC = math.Round(fC * cF) / cF


		fCS := fmt.Sprintf("%g", fC)

		CallClear()

		fmt.Println(fS + " " + units[ind - 1] + " ------> " + fCS + " " + units[ind2 - 1])
		fmt.Print("\n\n")

		fmt.Println("Enter Another Value of " + units[ind - 1] + " or \"quit\" to exit")

		for {
			text, _ = reader.ReadString('\n')
			text = strings.Replace(text, "\n", "", -1)

			if text == "quit" || text == "Quit" {
				os.Exit(0)
			} else {
				_, err := strconv.ParseFloat(text, 64)

				if err != nil {
					continue
				} else {
					break
				}
			}
		}
	}
}

func chooseType() string{
	var FILENAME string
	const fileName = "units.txt"
	var counter = 0
	reader := bufio.NewReader(os.Stdin)
	var ind int
	var m  = make(map[string]string)
	options := []string{}

	file, err := os.Open(fileName)

	if err != nil {
		log.Fatalf("Failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	fmt.Println("What would you like to convert?")

	for scanner.Scan() {
		temp := strings.Split(scanner.Text(), "\t")
		m[temp[0]] = temp[1]
		options = append(options, temp[0])
		fmt.Println(strconv.Itoa(counter + 1) + ". " + temp[0])
		counter++
	}

	file.Close()

	for {
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		ind, err = strconv.Atoi(text)

		if err == nil && (0 < ind && ind <= counter) {
			break
		}
	}

	FILENAME = m[options[ind - 1]]

	fmt.Println(FILENAME)

	return FILENAME
}

//Main function
func main()  {
	//Clears screen for better readability
	CallClear()

	FILENAME := chooseType()

	convert(getUnits(FILENAME))

}
