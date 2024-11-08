package main

import (
    rl "github.com/gen2brain/raylib-go/raylib"
)

const (
    screenWidth  = 1600
    screenHeight = 900
)

func main() {
    rl.InitWindow(screenWidth, screenHeight, "raylib [core] example - basic window")

    rl.SetTargetFPS(60)

    platformWidth := float32(screenWidth / 2)
    platform := Platform{
        Pos:   rl.NewVector2((screenWidth-platformWidth)/2, screenHeight-250), // Centered position of the blocker
        Size:  rl.NewVector2(platformWidth, 100),
        Texture: rl.LoadTexture("./assets/platform-1.png.png"),
    }

    for !rl.WindowShouldClose() {
        rl.BeginDrawing()

        rl.ClearBackground(rl.RayWhite)

        platform.DrawPlatform()

        rl.EndDrawing()
    }

    rl.CloseWindow()
}