package main

import (
	"testing"
)

func TestCollisionItem(t *testing.T) {

	hit1 := hitbox{
		LeftX:  10,
		LeftY:  10,
		RightX: 20,
		RightY: 20,
	}

	hit2 := hitbox{
		LeftX:  30,
		LeftY:  30,
		RightX: 40,
		RightY: 40,
	}
	if collisionItem(hit1, hit2) {
		t.Error("hit1 cross hit2")
	}

	hit1 = hitbox{
		LeftX:  10,
		LeftY:  10,
		RightX: 20,
		RightY: 20,
	}

	hit2 = hitbox{
		LeftX:  15,
		LeftY:  15,
		RightX: 50,
		RightY: 50,
	}

	if !collisionItem(hit1, hit2) {
		t.Error("hit1 not cross hit2")
	}
}

func TestCollisionStage(t *testing.T) {
	x := 10
	y := 10
	collisionStage(&x, &y)
	if x != 10 && y != 10 {
		t.Error("negative border. x or y changed !!!")
	}

	x = -10
	y = -10
	collisionStage(&x, &y)
	if x == -10 && y == -10 {
		t.Error("negative border. x or y value still -10 !!!")
	}
	if x != 0 && y != 0 {
		t.Error("negative border. x or y not zero !!!")
	}

	x = X_max + 10
	y = Y_max + 10
	collisionStage(&x, &y)
	if x != X_max && y != Y_max {
		t.Error("x or y not max value !!!")
	}
}
