package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/text/transform"
)

type MapFile struct {
	MapID       string `json:"map_id"`
	MapFilename string `json:"map_filename"`
	MapFilePath string `json:"map_filepath"`
}

type MapInfo struct {
	MapID   string `json:"map_id"`
	MapName string `json:"map_name"`
}

type MapConnection struct {
	FromMapID string `json:"from_map_id"`
	ToMapID   string `json:"to_map_id"`
}

func ReadMapFiles() []*MapFile {
	mapFiles := make([]*MapFile, 0)

	dirPath := fmt.Sprintf("%s/Mir200/Map/", rootPath)

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		filename := e.Name()
		ext := filepath.Ext(filename)
		id := strings.TrimSuffix(filename, ext)

		if ext != ".map" {
			continue
		}

		mapFiles = append(mapFiles, &MapFile{
			MapID:       id,
			MapFilename: filename,
			MapFilePath: dirPath + filename,
		})
	}

	return mapFiles
}

func ReadMapInfos() []*MapInfo {
	mapInfos := make([]*MapInfo, 0)

	file, err := os.Open(fmt.Sprintf("%s/Mir200/Envir/MapInfo.txt", rootPath))
	if err != nil {
		log.Fatal(err)
	}

	reader := transform.NewReader(file, enc.NewDecoder())
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, ";") || len(line) == 0 {
			continue
		}

		re := regexp.MustCompile(`\[(.*?)\]`)
		match := re.FindStringSubmatch(line)
		if len(match) > 1 {
			text := match[1]

			components := strings.Fields(text)
			if len(components) != 2 {
				continue
			}

			mapInfos = append(mapInfos, &MapInfo{
				MapID:   components[0],
				MapName: components[1],
			})
		} else {
			continue
		}
	}

	return mapInfos
}

func ReadMapConnections() []*MapConnection {
	mapConnections := make([]*MapConnection, 0)

	file, err := os.Open(fmt.Sprintf("%s/Mir200/Envir/MapInfo.txt", rootPath))
	if err != nil {
		log.Fatal(err)
	}

	reader := transform.NewReader(file, enc.NewDecoder())
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, ";") || len(line) == 0 {
			continue
		}

		components := strings.Split(line, "->")
		if len(components) != 2 {
			continue
		}

		from := strings.Fields(components[0])
		if len(from) != 3 {
			continue
		}
		to := strings.Fields(components[1])
		if len(to) != 3 {
			continue
		}

		mapConnections = append(mapConnections, &MapConnection{
			FromMapID: from[0],
			ToMapID:   to[0],
		})
	}

	return mapConnections
}
