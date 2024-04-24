package main

import (
	"fmt"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// if we get an argument, change working directory there. Panic if fail.
// Then, list all directory members.
func main() {
	fmt.Print("program ")
	fmt.Println(os.Args[:1])         //first element in args is the program
	suppliedArguments := os.Args[1:] //first element in args is the program
	fmt.Print("arguments: ")
	fmt.Println(suppliedArguments)

	if len(suppliedArguments) > 0 {
		err := os.Chdir(suppliedArguments[0])
		check(err)
	}

	workingDir, err := os.Getwd()
	check(err)
	dirSlice, err := os.ReadDir(workingDir)
	check(err)

	fmt.Print("working in dir ")
	fmt.Println(workingDir)

	fmt.Println()
	fmt.Println(" Files: ")

	for _, entry := range dirSlice {
		fmt.Println(entry.Name()) //fs.DirEntry (interface)
	}

	fmt.Println(len(dirSlice), "entries")
}
