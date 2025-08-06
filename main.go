package main

import (
	"fastline/exercises"
	"fmt"
	"os"
)

var exerciseMapper = map[string]func(){
	"exercise1": exercises.Exercise1,
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage : go run main.go <exerciseID>")
		fmt.Println("Available exercises are :")
		for k := range exerciseMapper {
			fmt.Println(k)
		}
		return
	}
	fmt.Println("Checking in github")
	exerciseId := os.Args[1]
	if exercise, exists := exerciseMapper[exerciseId]; exists {
		fmt.Printf("Running %s....", exerciseId)
		exercise()
	} else {
		fmt.Println("Exercise not found")
	}

}
