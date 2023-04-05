package main

import (
	"log"
	"os"
	"strings"

	"github.com/alexflint/go-arg"
)

var args struct {
	RootPath string `arg:"-r,--root-path" help:"MirServer folder path"`
}

var rootPath string

var mapFiles []*MapFile
var mapInfos []*MapInfo
var mapConnections []*MapConnection
var monGens []*MonGen
var monItems []*MonItem

const (
	SourceMaskUnknown       uint = 0x00
	SourceMaskMapFile       uint = 0x01 << 0
	SourceMaskMapInfo       uint = 0x01 << 1
	SourceMaskMapConnection uint = 0x01 << 2
	SourceMaskMonGen        uint = 0x01 << 3
	SourceMaskMonItem       uint = 0x01 << 4
)

func main() {
	arg.MustParse(&args)

	rootPath = args.RootPath
	if len(rootPath) == 0 {
		panic("root path empty")
	}

	UnifyMaps()
	UnifyMonItems()

	mapFiles = ReadMapFiles()
	mapInfos = ReadMapInfos()
	mapConnections = ReadMapConnections()
	monGens = ReadMonGens()
	monItems = ReadMonItems()

	checkForUnusedMapFiles()
	checkForMissingMapFiles()
	checkForIdleMapFiles()
	checkForMissingMonItems()
	checkForUnusedMonItems()
}

func checkForUnusedMapFiles() {
	log.Print("[Unused map files]")

	mapIDs := make(map[*MapFile]uint)

	for _, mapFile := range mapFiles {
		if !Contains(mapInfos, func(mapInfo *MapInfo) bool {
			return mapFile.MapID == mapInfo.MapID
		}) {
			mapIDs[mapFile] |= SourceMaskMapInfo
		}
	}

	for _, mapFile := range mapFiles {
		if !Contains(mapConnections, func(mapConnection *MapConnection) bool {
			return mapFile.MapID == mapConnection.FromMapID || mapFile.MapID == mapConnection.ToMapID
		}) {
			mapIDs[mapFile] |= SourceMaskMapConnection
		}
	}

	for _, mapFile := range mapFiles {
		if !Contains(monGens, func(monGen *MonGen) bool {
			return mapFile.MapID == monGen.MapID
		}) {
			mapIDs[mapFile] |= SourceMaskMonGen
		}
	}

	for mapFile, usage := range mapIDs {
		files := make([]string, 0)
		if usage&SourceMaskMapInfo != 0x0 {
			files = append(files, "MapInfo")
		}
		if usage&SourceMaskMapConnection != 0x0 {
			files = append(files, "MapConnection")
		}
		if usage&SourceMaskMonGen != 0x0 {
			files = append(files, "MonGen")
		}

		log.Printf("mapfile[%s] is not used by [%s]", mapFile.MapID, strings.Join(files, "|"))

		if usage == SourceMaskMapInfo|SourceMaskMapConnection|SourceMaskMonGen {
			os.Remove(mapFile.MapFilePath)
		}
	}
}

func checkForMissingMapFiles() {
	log.Print("[Missing map files]")

	mapIDs := make(map[string]uint)

	for _, mapInfo := range mapInfos {
		if !Contains(mapFiles, func(mapFile *MapFile) bool {
			return mapFile.MapID == mapInfo.MapID
		}) {
			mapIDs[mapInfo.MapID] |= SourceMaskMonGen
		}
	}

	for _, mapConnection := range mapConnections {
		if !Contains(mapFiles, func(mapFile *MapFile) bool {
			return mapConnection.FromMapID == mapFile.MapID
		}) {
			mapIDs[mapConnection.FromMapID] |= SourceMaskMapConnection
		}
	}
	for _, mapConnection := range mapConnections {
		if !Contains(mapFiles, func(mapFile *MapFile) bool {
			return mapConnection.ToMapID == mapFile.MapID
		}) {
			mapIDs[mapConnection.ToMapID] |= SourceMaskMapConnection
		}
	}

	for _, monGen := range monGens {
		if !Contains(mapFiles, func(mapFile *MapFile) bool {
			return monGen.MapID == mapFile.MapID
		}) {
			mapIDs[monGen.MapID] |= SourceMaskMonGen
		}
	}

	for mapID, usage := range mapIDs {
		files := make([]string, 0)
		if usage&SourceMaskMapInfo != 0x0 {
			files = append(files, "MapInfo")
		}
		if usage&SourceMaskMapConnection != 0x0 {
			files = append(files, "MapConnection")
		}
		if usage&SourceMaskMonGen != 0x0 {
			files = append(files, "MonGen")
		}

		log.Printf("mapfile[%s] is expected by [%s]", mapID, strings.Join(files, "|"))
	}
}

func checkForIdleMapFiles() {
	log.Print("[Idle map files]")

	for _, mapFile := range mapFiles {
		if !Contains(monGens, func(monGen *MonGen) bool {
			return mapFile.MapID == monGen.MapID
		}) {
			log.Printf("mapfile[%s] is idle", mapFile.MapID)
		}
	}
}

func checkForMissingMonItems() {
	log.Print("[Missing mon items]")

	for _, monGen := range monGens {
		if !Contains(monItems, func(monItem *MonItem) bool {
			return monGen.MonName == monItem.MonName
		}) {
			log.Printf("monitem[%s] is expected by [MonGen]", monGen.MonName)
		}
	}
}

func checkForUnusedMonItems() {
	log.Print("[Unused mon items]")

	for _, monItem := range monItems {
		if !Contains(monGens, func(monGen *MonGen) bool {
			return monItem.MonName == monGen.MonName
		}) {
			log.Printf("monitem[%s] is not used", monItem.MonName)
		}
	}
}
