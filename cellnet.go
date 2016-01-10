package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type point struct {
	x int
	y int
}

type tile struct {
	p        point
	width    int
	height   int
	tiletype int
}

const (
	tilewidth  = 32
	tileheight = 32
)

var tilemap []tile
var jsonmap jsonobject

type jsonobject struct {
	Height      int
	Layers      []LayersObj
	Orientation string
	Tileheight  int
	Tilesets    []TilesetsObj
	Tilewidth   int
	Version     int
	Width       int
}

type LayersObj struct {
	Data    []int
	Height  int
	Name    string
	Opacity int
	Type    string
	Visible bool
	Width   int
	X       int
	Y       int
}

type TilesetsObj struct {
	Firstgid    int
	Image       string
	Imageheight int
	Imagewidth  int
	Margin      int
	Name        string
	Spacing     int
	tileheight  int
	tilewidth   int
}

func ParseMap() {

	file, e := ioutil.ReadFile("/home/drpsycho/js/jsgame/stuff/sprites/map.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	// fmt.Printf("%s\n", string(file))

	json.Unmarshal(file, &jsonmap)
}

func FillTileMap() {
	AmountTileX := jsonmap.Layers[0].Width
	AmountTileY := jsonmap.Layers[0].Height
	// tileByY := Y_max/tileheight
	globalit := 0
	for _y := 0; _y < AmountTileY; _y++ {
		for _x := 0; _x < AmountTileX; _x++ {
			tilemap = append(tilemap, tile{
				width:    tilewidth,
				height:   tileheight,
				tiletype: jsonmap.Layers[0].Data[globalit],
				p: point{
					x: _x * 32, //tilewidth,
					y: _y * 32, //tileheight,
				},
			})
			globalit += 1
		}
	}
}
