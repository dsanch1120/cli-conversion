/*
	Created by Daniel Sanchez
	May 21st, 2020
	Allows user to make unit conversions on a terminal
 */

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
//Map that will clear the terminal screen
var clear map[string]func()
//Boolean that allows user to restart the program with a new conversion
var cont = true

//Allows the program to clear the terminal in Windows and Unix operating systems
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
//Clears the Screen
func CallClear() {
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	}
}
//Gets the units from a text file and reads them into a slice
func getUnits(fileName string) (map[string][]string, []string) {
	//Clears terminal for readability
	CallClear()
	//Variables to be used throughout the function
	conversion := make(map[string][]string)
	units := []string{}
	var counter = 0
	//Opens file
	file, err := os.Open(fileName)
	//Checks that file opened without error
	if err != nil {
		log.Fatalf("Failed opening file: %s", err)
	}
	//Scanner to read through file
	scanner := bufio.NewScanner(file)
	//Reads through file line by line
	for scanner.Scan() {
		//Splits each line into a slice
		temp := strings.Split(scanner.Text(), "\t")
		conversion[temp[0]] = temp[1:]
		units = append(units, temp[0])
		counter++
	}

	file.Close()

	return conversion, units
}
//Handles the conversions
func convert(conversion map[string][]string, units []string) {
	//Variables used throughout function
	var ind int		//Index of first unit
	var ind2 int	//Index of second unit
	var err error	//Handles potential errors
	//Reader to be used throughout function. Gets input from user
	reader := bufio.NewReader(os.Stdin)
	//Prompts user to choose unit to convert
	fmt.Println("Choose Unit to be Converted")
	//Prints each unit that can be converted
	for count := 0; count < len(units); count++ {
		fmt.Println(strconv.Itoa(count + 1) + ". " + units[count])
	}
	//Loop that ensures user enters a correct choice
	for {
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		ind, err = strconv.Atoi(text)
		//Checks that user entered a number in the correct range
		if err == nil && (0 < ind && ind <= len(units)) {
			break
		}
	}
	//Clears terminal for increased readability
	CallClear()
	//Prompts user to choose unit to be converted to
	fmt.Println("Choose Unit to Convert " + units[ind - 1] + " to")
	//Displays conversion options to user (not including their original choice)
	for count := 0; count < len(units); count++ {
		if (count != ind - 1) {
			fmt.Println(strconv.Itoa(count + 1) + ". " + units[count])
		}
	}
	//Ensures user gave a correct input
	for {
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		ind2, err = strconv.Atoi(text)
		//Checks that user entered a number in the correct range
		if err == nil && ind2 != ind && (0 < ind2 && ind2 <= len(units)) {
			break
		}
	}
	//Clears terminal for better readability
	CallClear()
	//Prompts user to enter a number for conversion
	fmt.Println("How many " + units[ind - 1] + " to convert? (input \"quit\" to exit)")
	//Handles user input
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	//Calculates the conversion
	for {
		//Ends program if user input "quit"
		if text == "quit" || text == "Quit" {
			CallClear()
			os.Exit(0)
		}
		//Converts user's input to a floating point
		f, err := strconv.ParseFloat(text, 64)
		fS := fmt.Sprintf("%g", f)
		//If user gave an incorrect input, loop will iterate until they give a correct one
		if err != nil {
			text, _ := reader.ReadString('\n')
			text = strings.Replace(text, "\n", "", -1)
			continue
		}
		//Changes conversion rates from strings to float64 variables
		iC, err := strconv.ParseFloat(conversion[units[ind - 1]][1], 64)
		sC, err := strconv.ParseFloat(conversion[units[ind2 - 1]][0], 64)
		cF, _ := strconv.ParseFloat(conversion[units[ind2 - 1]][2], 64)
		//Calculates the conversion
		fC := f * iC * sC
		fC = math.Round(fC * cF) / cF
		//Creates a string value based upon conversion result
		fCS := fmt.Sprintf("%g", fC)
		//Clears terminal screen for increased readability
		CallClear()
		//Displays result to user
		fmt.Println(fS + " " + units[ind - 1] + " ------> " + fCS + " " + units[ind2 - 1])
		fmt.Print("\n\n")
		//Prompts user to enter another value for conversion, restart, or end the program
		fmt.Println("Enter Another Value of " + units[ind - 1] + " for another conversion")
		fmt.Println("Enter \"Restart\" to make a new conversion")
		fmt.Println("Enter \"Quit\" to end program")
		//Handles user's input
		for {
			//Gets user's input
			text, _ = reader.ReadString('\n')
			text = strings.Replace(text, "\n", "", -1)
			// Handles user's input
			if text == "quit" || text == "Quit" {
				CallClear()
				os.Exit(0)
			} else if text == "restart" || text == "Restart" {
				cont = true
				return
			} else {
				//Checks if user input a float64 varaible
				_, err := strconv.ParseFloat(text, 64)
				//Reloops this loop if user gave a bad input
				if err != nil {
					continue
				} else {	//Exits this loop if user gave a good input
					break
				}
			}
		}
	}
}
//Allows the user to choose what type of conversion the program will do
func chooseType() string{
	//Variables to be used throughout the function
	var FILENAME string						//Name of the file (to be assigned)
	const fileName = "units.txt"			//Name of the file containing the file names
	var counter = 0							//Counter to be used throughout
	reader := bufio.NewReader(os.Stdin)		//Bufio Reader to be used throughout
	var ind int								//Index to help choose file
	var m  = make(map[string]string)		//Maps the description to the file names
	options := []string{}					//Slice of strings with file descriptions
	//Opens the file
	file, err := os.Open(fileName)
	//Checks if file opened properly
	if err != nil {
		log.Fatalf("Failed opening file: %s", err)
	}
	//Scanner to read through file
	scanner := bufio.NewScanner(file)
	//Prompts user to choose type of conversion
	fmt.Println("What would you like to convert?")
	//For loop that iterates through file, adds to variables, displays file descriptions
	for scanner.Scan() {
		temp := strings.Split(scanner.Text(), "\t")
		m[temp[0]] = temp[1]
		options = append(options, temp[0])
		fmt.Println(strconv.Itoa(counter + 1) + ". " + temp[0])
		counter++
	}

	file.Close()
	//Loop that iterates until user correctly chooses a conversion type
	for {
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		ind, err = strconv.Atoi(text)
		//Ensures user gives a correct input
		if err == nil && (0 < ind && ind <= counter) {
			break
		}
	}
	//Uses the map to assign the file name
	FILENAME = m[options[ind - 1]]

	return FILENAME
}
//Main function
func main()  {
	//Loop that runs until the user inputs "quit" in the convert function
	for cont {
		//Clears screen for better readability
		CallClear()
		//Gets the filename from the chooseType function
		FILENAME := chooseType()
		//Calls the functions that handle the program
		convert(getUnits(FILENAME))
	}
}
