package main

import rl "github.com/gen2brain/raylib-go/raylib"

type ProgressBar struct {
	X            int32
	Y            int32
	Width        int32
	Height       int32
	BorderWidth  int32
	BorderHeight int32
	progress     float32
	BaseColor    rl.Color
	AccentColor  rl.Color
	BorderColor  rl.Color
}

func (pb *ProgressBar) SetProgress(newProgress float32) {
	pb.progress = newProgress
	if pb.progress < 0 {
		pb.progress = 0
	}
	if pb.progress > 1 {
		pb.progress = 1
	}
}

func (pb ProgressBar) DrawBar() {
	rl.DrawRectangle(pb.X-pb.BorderWidth, pb.Y-pb.BorderHeight, pb.Width+(2*pb.BorderWidth), pb.Height+(2*pb.BorderHeight), pb.BorderColor)
	rl.DrawRectangle(pb.X, pb.Y, pb.Width, pb.Height, pb.BaseColor)
	rl.DrawRectangle(pb.X, pb.Y, int32(pb.progress*float32(pb.Width)), pb.Height, pb.AccentColor)
}

func (pb ProgressBar) DrawBarReverse() {
	rl.DrawRectangle(pb.X-pb.BorderWidth, pb.Y-pb.BorderHeight, pb.Width+(2*pb.BorderWidth), pb.Height+(2*pb.BorderHeight), pb.BorderColor)
	rl.DrawRectangle(pb.X, pb.Y, pb.Width, pb.Height, pb.BaseColor)
	rl.DrawRectangle(pb.X+(pb.Width-int32(pb.progress*float32(pb.Width))), pb.Y, int32(pb.progress*float32(pb.Width)), pb.Height, pb.AccentColor)
}

func NewProgressBar(newX, newY, newWidth, newHeight, newBorderWidth, newBorderHeight int32) ProgressBar {
	pb := ProgressBar{X: newX, Y: newY, Width: newWidth, Height: newHeight, BorderWidth: newBorderWidth, BorderHeight: newBorderHeight, BaseColor: rl.Red, AccentColor: rl.Green, BorderColor: rl.Black}
	pb.progress = 1
	return pb
}
