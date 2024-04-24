package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// if we get an argument, change working directory there. Panic if fail.
// Then, list all directory members.
func main() {
	suppliedArguments := os.Args[1:] //first element in args is the program

	if len(suppliedArguments) > 0 {
		err := os.Chdir(suppliedArguments[0])
		check(err)
	}

	workingDir, err := os.Getwd()
	check(err)
	dirSlice, err := os.ReadDir(workingDir)
	check(err)
	fmt.Println("working in", workingDir)

	var zipFilenames []string

	for _, entry := range dirSlice {
		if filepath.Ext(entry.Name()) == ".zip" {
			zipFilenames = append(zipFilenames, entry.Name())
		}
	}

	var archiveContents *zip.ReadCloser

	fmt.Println(len(dirSlice), "files", len(zipFilenames), "zipfiles")
	for _, archiveFilename := range zipFilenames {
		fmt.Println(archiveFilename)

		archiveContents, err = zip.OpenReader(archiveFilename)
		check(err)

		for _, potentialFileOfInterest := range archiveContents.File {
			var partDirName string
			dirMade := false

			extention := filepath.Ext(potentialFileOfInterest.Name)
			if extention == ".kicad_sym" || extention == ".kicad_mod" {
				if !dirMade {
					dirMade = true
					if archiveFilename[:4] == "LIB_" && len(archiveFilename) > 8 { //len("LIB_.zip") is 8
						partDirName = archiveFilename[4 : len(archiveFilename)-4]
					} else {
						partDirName = archiveFilename[:len(archiveFilename)-4]
					}
					os.Mkdir(partDirName, 0755)
				}

				ioReadCloser, err := potentialFileOfInterest.Open()
				check(err)
				fileBytes, err := io.ReadAll(ioReadCloser)
				check(err)
				err = os.WriteFile(filepath.Join(partDirName, filepath.Base(potentialFileOfInterest.Name)), fileBytes, 0655)
				check(err)
				defer ioReadCloser.Close()
			}
		}

		archiveContents.Close()
	}
}
