package main

// import (
// 	"log"
// )

func collisionItem(_it1 hitbox, _it2 hitbox) bool {
	if _it1.LeftX > _it2.RightX || _it2.LeftX > _it1.RightX {
		return false
	}

	if _it1.LeftY > _it2.RightY || _it2.LeftY > _it1.RightY {
		return false
	}

	return true
}

func collisionStage(_x *int, _y *int) {
	if *_x < 0 {
		*_x = 0
	} else if *_x > X_max {
		*_x = X_max
	}
	if *_y < 0 {
		*_y = 0
	} else if *_y > Y_max {
		*_y = Y_max
	}
}

func collisionBlock(_hit hitbox) bool {
	for it := range tilemap {
		if tilemap[it].tiletype != 0 {
			ithitbox := hitbox{
				LeftX:  tilemap[it].p.x,
				LeftY:  tilemap[it].p.y,
				RightX: tilemap[it].p.x + 32,
				RightY: tilemap[it].p.y + 32,
			}
			if collisionItem(ithitbox, _hit) {
				// log.Println("tile - ", tilemap[it].p)
				return true
			}
		}
	}
	return false
}

//32 16
func collision(d *actionmsg) {
	collisionStage(&d.In.Data.Dx, &d.In.Data.Dy)

	_player := hitbox{
		LeftX:  d.In.Data.Dx + (16 - 4),
		LeftY:  d.In.Data.Dy + (16 - 4),
		RightX: d.In.Data.Dx + (16 + 4),
		RightY: d.In.Data.Dy + (16 + 4),
	}
	if collisionBlock(_player) {
		d.In.Data.Dx = *d.Dx
		d.In.Data.Dy = *d.Dy
	} else {
		*d.Dx = d.In.Data.Dx
		*d.Dy = d.In.Data.Dy
	}

	_player = hitbox{
		LeftX:  d.In.Data.Dx,
		LeftY:  d.In.Data.Dy,
		RightX: d.In.Data.Dx + 32,
		RightY: d.In.Data.Dy + 32,
	}

	for it := range items {
		if collisionItem(items[it].HitBox, _player) {
			removeitem <- removeitemstruct{Name: d.In.Name, Id: items[it].GlobalId}
		}
	}
}
