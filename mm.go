package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func dieIf(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

/*
if we get an argument, change working directory to there or die.

then, rummage through any .zip archives, looking for kicad_sym or kicad_mod
if they're found, make a new directory named after the .zip (remove LIB_ prefix if present)
*/
func main() {
	suppliedArguments := os.Args[1:] //first element in args is the program

	if len(suppliedArguments) > 0 {
		err := os.Chdir(suppliedArguments[0])
		dieIf(err)
	}

	workingDir, err := os.Getwd()
	dieIf(err)
	dirSlice, err := os.ReadDir(workingDir)
	dieIf(err)
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
		dieIf(err)

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
				dieIf(err)
				fileBytes, err := io.ReadAll(ioReadCloser)
				dieIf(err)
				err = os.WriteFile(filepath.Join(partDirName, filepath.Base(potentialFileOfInterest.Name)), fileBytes, 0655)
				dieIf(err)
				defer ioReadCloser.Close()
			}
		}

		archiveContents.Close()
	}
}
