package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Platform struct {
	Pos   rl.Vector2
	Size  rl.Vector2
	Color rl.Color
}

func (p Platform) DrawPlatform() {
	rl.DrawRectangle(int32(p.Pos.X), int32(p.Pos.Y), int32(p.Size.X), int32(p.Size.Y), p.Color)
}

func CheckCollision(box *Box, platform Platform) {
	if rl.CheckCollisionRecs( //Raylib let's us quickly check overlap with the rectangle class.
		rl.NewRectangle(box.Pos.X, box.Pos.Y, box.Size.X, box.Size.Y),
		rl.NewRectangle(platform.Pos.X, platform.Pos.Y, platform.Size.X, platform.Size.Y),
	) {
		if box.Pos.Y+box.Size.Y > platform.Pos.Y && box.Vel.Y > 0 { //now check which side to stop the velocity
			box.Pos.Y = platform.Pos.Y - box.Size.Y //move box in case of overlap
			box.Vel.Y = 0                          //stop the box from moving further
		}
		if box.Pos.Y < platform.Pos.Y+platform.Size.Y && box.Vel.Y < 0 {
			box.Pos.Y = platform.Pos.Y + platform.Size.Y
			box.Vel.Y = 0
		}
		if box.Pos.X+box.Size.X > platform.Pos.X && box.Vel.X > 0 {
			box.Pos.X = platform.Pos.X - box.Size.X
			box.Vel.X = 0
		}
		if box.Pos.X < platform.Pos.X+platform.Size.X && box.Vel.X < 0 {
			box.Pos.X = platform.Pos.X + platform.Size.X
			box.Vel.X = 0
		}
	}
}
