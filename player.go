package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Box struct {
	Pos   rl.Vector2
	Vel   rl.Vector2
	Size  rl.Vector2
	Color rl.Color
}

func (b *Box) ApplyGravity(g rl.Vector2) {
	b.Vel = rl.Vector2Add(b.Vel, rl.Vector2Scale(g, rl.GetFrameTime()))
}

func (b *Box) UpdateBox() {
	b.Pos = rl.Vector2Add(b.Pos, rl.Vector2Scale(b.Vel, rl.GetFrameTime()))
}

func (b Box) DrawBox() {
	rl.DrawRectangle(int32(b.Pos.X), int32(b.Pos.Y), int32(b.Size.X), int32(b.Size.Y), b.Color)
}
