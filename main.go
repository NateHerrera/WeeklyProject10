package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 1600
	screenHeight = 900
)

func main() {
	rl.InitWindow(screenWidth, screenHeight, "raylib [core] example - seamless parallax background")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	// Load textures for each layer
	backgroundTexture := rl.LoadTexture("./assets/backgroundproj10.png")
	midgroundTexture := rl.LoadTexture("./assets/middlegroundproj10.png")
	foregroundTexture := rl.LoadTexture("./assets/foregroundproj10.png")

	// Load the platform
	platformWidth := float32(screenWidth / 2)
	platform := Platform{
		Pos:   rl.NewVector2((screenWidth-platformWidth)/2, screenHeight-250),
		Size:  rl.NewVector2(platformWidth, 100),
		Color: rl.DarkGreen,
	}

	// Set background scale for resizing
	backgroundScale := float32(screenHeight) / float32(backgroundTexture.Height)

	// Scrolling variables for each parallax layer
	var scrollingBack, scrollingMid, scrollingFore float32 = 0.0, 0.0, 0.0

	for !rl.WindowShouldClose() {
		// Update positions for parallax effect
		scrollingBack -= 1.0 // Slowest layer
		scrollingMid -= 2.0  // Middle layer speed
		scrollingFore -= 4.5 // Foreground layer speed

		// Reset positions for seamless effect
		if scrollingBack <= -float32(backgroundTexture.Width)*backgroundScale {
			scrollingBack = 0
		}
		if scrollingMid <= -float32(midgroundTexture.Width)*backgroundScale {
			scrollingMid = 0
		}
		if scrollingFore <= -float32(foregroundTexture.Width)*backgroundScale {
			scrollingFore = 0
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		// Draw each layer twice to ensure seamless scrolling
		drawSeamlessLayer(backgroundTexture, scrollingBack, backgroundScale)
		drawSeamlessLayer(midgroundTexture, scrollingMid, backgroundScale)
		drawSeamlessLayer(foregroundTexture, scrollingFore, backgroundScale)

		// Draw the platform
		platform.DrawPlatform()

		rl.EndDrawing()
	}
}

// Helper function to draw seamless parallax layer with scaling
func drawSeamlessLayer(texture rl.Texture2D, positionX float32, scale float32) {
	// Draw the texture twice to ensure continuous scrolling effect
	rl.DrawTextureEx(texture, rl.NewVector2(positionX, 0), 0, scale, rl.White)
	rl.DrawTextureEx(texture, rl.NewVector2(positionX+float32(texture.Width)*scale, 0), 0, scale, rl.White)
}
