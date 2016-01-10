package main

import (
	"encoding/json"
	// "log"
	"math/rand"
	"time"
)

// addItem(msg.WorldX, msg.WorldY, msg.Item, msg.GlobalId);
type hitbox struct {
	LeftX  int
	RightX int
	LeftY  int
	RightY int
}

type stageitem struct {
	Name     string
	GlobalId int
	Item     int
	WorldX   int
	WorldY   int
	Witdh    int
	Height   int
	HitBox   hitbox
}

type stagemsgall struct {
	Action string
	Data   []stageitem
}

type stagemsg struct {
	Action string
	Data   stageitem
}

const (
	ItemAmount    = 40
	ItemSpriteMax = 35
	X_max         = 32 * 30
	Y_max         = 32 * 30

	FruitWidth  = 20
	FruitHeight = 20
)

var items []stageitem

type removeitemstruct struct {
	Name string
	Id   int
}

var removeitem = make(chan removeitemstruct, 100)

func StageUpdater() {
	for {
		if !(len(items) > 0) {
			fillStage()
			HubHandler.broadcast <- GetStageJSON()
		}

		_item := <-removeitem
		idremove := 0
		approveremove := false

		for it := range items {
			if items[it].GlobalId == _item.Id {
				res, _ := json.Marshal(inputmsg{
					Name:   _item.Name,
					Action: "pickupitem",
					Data:   inputmsgData{Id: _item.Id},
				})

				HubHandler.broadcast <- res
				idremove = it
				approveremove = true
			}
		}

		if approveremove {
			items = append(items[:idremove], items[idremove+1:]...)
			// log.Println(items)
		}
	}
}

func GetStageJSON() []byte {
	_s := stagemsgall{"additem", items}
	res, _ := json.Marshal(_s)
	return res
}

func fillStage() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s1)

	for i := 0; i < ItemAmount; i++ {
		_width := FruitWidth
		_heigth := FruitHeight
		_x := r.Intn(X_max-(64+FruitWidth)) + 32
		_y := r.Intn(200-(64+FruitWidth)) + 32

		items = append(items, stageitem{
			Item:     r.Intn(ItemSpriteMax),
			WorldX:   _x,
			WorldY:   _y,
			Witdh:    _width,
			Height:   _heigth,
			GlobalId: i,
			HitBox: hitbox{
				LeftX:  (_x),
				LeftY:  (_y),
				RightX: (_x + _heigth),
				RightY: (_y + _width),
			},
			Name: "fruitnveg"})
		// log.Println("Added item ", items[i])
	}
}
