package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type MonGen struct {
	MapID   string `json:"map_id"`
	MonName string `json:"mon_name"`
}

type MonItem struct {
	MonName string `json:"mon_name"`
}

var enc = simplifiedchinese.GBK

func ReadMonGens() []*MonGen {
	file, err := os.Open(fmt.Sprintf("%s/Mir200/Envir/MonGen.txt", rootPath))
	if err != nil {
		log.Fatal(err)
	}

	reader := transform.NewReader(file, enc.NewDecoder())
	scanner := bufio.NewScanner(reader)

	monGens := make([]*MonGen, 0)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, ";") || len(line) == 0 {
			continue
		}

		components := strings.Fields(line)
		if len(components) != 7 {
			continue
		}

		monGens = append(monGens, &MonGen{
			MapID:   components[0],
			MonName: components[3],
		})
	}

	return monGens
}

func ReadMonItems() []*MonItem {
	monItem := make([]*MonItem, 0)

	entries, err := os.ReadDir(fmt.Sprintf("%s/Mir200/Envir/MonItems/", rootPath))
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		filename := e.Name()
		name := strings.TrimSuffix(filename, filepath.Ext(filename))
		monItem = append(monItem, &MonItem{
			MonName: name,
		})
	}

	return monItem
}
