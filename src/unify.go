package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func UnifyMaps() {
	dirPath := fmt.Sprintf("%s/Mir200/Map/", rootPath)

	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		newExt := strings.ToLower(filepath.Ext(file.Name()))
		newName := strings.ToUpper(strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())))
		os.Rename(dirPath+file.Name(), dirPath+newName+newExt)
		// log.Print("mapfile unify from: ", dirPath+""+file.Name(), " to: ", dirPath+""+newName+newExt)
	}
}

func UnifyMonItems() {
	dirPath := fmt.Sprintf("%s/Mir200/Envir/MonItems/", rootPath)

	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		newExt := strings.ToLower(filepath.Ext(file.Name()))
		newName := strings.ToUpper(strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())))
		os.Rename(dirPath+file.Name(), dirPath+newName+newExt)
		// log.Print("monitem unify from: ", dirPath+""+file.Name(), " to: ", dirPath+""+newName+newExt)
	}
}
