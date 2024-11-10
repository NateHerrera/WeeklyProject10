package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	IDLESTATE   = "idle"
	WALKSTATE   = "walk"
	JUMPSTATE   = "jump"
	BLOCKSTATE  = "block"
	ATTACKSTATE = "attack"
	HEAVYSTATE  = "heavy"
)

type Transform struct {
	Pos   rl.Vector2
	Flip  int
	Scale rl.Vector2
}

func NewTransform(newPos rl.Vector2) *Transform {
	return &Transform{Pos: newPos, Flip: 1, Scale: rl.NewVector2(128, 128)}
}

type Animation struct {
	*Transform
	SpriteSheet  rl.Texture2D
	MaxIndex     int
	CurrentIndex int
	Timer        float32
	SwitchTime   float32
	Name         string
}

func (a *Animation) TickState() {
	a.UpdateTime()
	a.DrawAnimation()
}

func (a *Animation) GetName() string {
	return a.Name
}

func (a *Animation) ResetTime() {
	a.Timer = 0
	a.CurrentIndex = 0
}

func NewAnimation(newTransform *Transform, newSheet rl.Texture2D, newTime float32, newName string) Animation {
	spriteDimension := newSheet.Height
	frames := int(newSheet.Width / spriteDimension)
	newAnimation := Animation{
		Transform:    newTransform,
		SpriteSheet:  newSheet,
		MaxIndex:     frames - 1,
		CurrentIndex: 0,
		Timer:        0,
		SwitchTime:   newTime,
		Name:         newName,
	}
	return newAnimation
}

func (a *Animation) UpdateTime() {
	a.Timer += rl.GetFrameTime()
	if a.Timer > a.SwitchTime {
		a.Timer = 0
		a.CurrentIndex++
		if a.CurrentIndex > a.MaxIndex {
			a.CurrentIndex = 0
		}
	}
}

func (a Animation) DrawAnimation() {
	sourceRect := rl.NewRectangle(float32(64*a.CurrentIndex), 0, 64*float32(a.Flip), 64)
	destRect := rl.NewRectangle(a.Pos.X, a.Pos.Y, float32(a.Scale.X), float32(a.Scale.Y))
	//origin := rl.NewVector2(a.Scale.X/2, a.Scale.Y/2)
	rl.DrawTexturePro(a.SpriteSheet, sourceRect, destRect, rl.Vector2Zero(), 0, rl.White)
}
