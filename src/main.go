package main

import (
	"log"

	"github.com/alexflint/go-arg"
)

var args struct {
	RootPath string `arg:"-r,--root-path" help:"MirServer folder path"`
}

var rootPath string

var mapFiles []*MapFile
var mapInfos []*MapInfo
var monGens []*MonGen
var monItems []*MonItem

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
	monGens = ReadMonGens()
	monItems = ReadMonItems()

	validateMaps()
	validateMapInfo()
	validateMonGen()
}

func validateMaps() {
	// Check if all map files are used in the MapInfo
	log.Print("[Unused map files]")
	for _, mapFile := range mapFiles {
		if !Contains(mapInfos, func(mapInfo *MapInfo) bool {
			return mapFile.MapID == mapInfo.MapID
		}) {
			log.Printf("mapfile[%s]", mapFile.MapID)
		}
	}

	// Check if all the maps listed in MapInfo have corresponding map files
	log.Print("[Missing map files]")
	for _, mapInfo := range mapInfos {
		if !Contains(mapFiles, func(mapFile *MapFile) bool {
			return mapFile.MapID == mapInfo.MapID
		}) {
			log.Printf("mapinfo[%s]", mapInfo.MapID)
		}
	}
}

func validateMapInfo() {
	// Check if the MapInfo has no MonGen
	log.Print("[Missing MonGen]")
	for _, mapInfo := range mapInfos {
		if !Contains(monGens, func(monGen *MonGen) bool {
			return mapInfo.MapID == monGen.MapID
		}) {
			log.Printf("mapinfo[%s]", mapInfo.MapID)
		}
	}
}

func validateMonGen() {
	// Check if all the items listed in MonGen have corresponding MapInfo
	log.Print("[Missing MapInfo]")
	for _, monGen := range monGens {
		if !Contains(mapInfos, func(mapInfo *MapInfo) bool {
			return monGen.MapID == mapInfo.MapID
		}) {
			log.Printf("mongen[%s]", monGen.MapID)
		}
	}

	// Check if all the items listed in MonGen have corresponding map file
	log.Print("[Missing map file]")
	for _, monGen := range monGens {
		if !Contains(mapFiles, func(mapFile *MapFile) bool {
			return monGen.MapID == mapFile.MapID
		}) {
			log.Printf("mongen[%s]", monGen.MapID)
		}
	}

	// Check if the MonGen has no MonItem
	log.Print("[Missing MonItem]")
	for _, monGen := range monGens {
		if !Contains(monItems, func(monItem *MonItem) bool {
			return monGen.MonName == monItem.MonName
		}) {
			log.Printf("mongen[%s]", monGen.MonName)
		}
	}

	// Check if all the MonItem have corresponding MonGen
	log.Print("[Missing MonGen]")
	for _, monItem := range monItems {
		if !Contains(monGens, func(monGen *MonGen) bool {
			return monItem.MonName == monGen.MonName
		}) {
			log.Printf("monitem[%s]", monItem.MonName)
		}
	}
}
