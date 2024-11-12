package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 1600
	screenHeight = 900
)

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Stickman Fighters")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	// Load textures for each layer
	backgroundTexture := rl.LoadTexture("./assets/backgroundproj10.png")
	midgroundTexture := rl.LoadTexture("./assets/middlegroundproj10.png")
	foregroundTexture := rl.LoadTexture("./assets/foregroundproj10.png")

	// Load the platform
	platformWidth := float32(screenWidth / 2)
	platform1 := Platform{
		Pos:   rl.NewVector2((screenWidth-platformWidth)/2, screenHeight-250),
		Size:  rl.NewVector2(platformWidth, 50),
		Color: rl.White,
	}

	platform2 := Platform{
		Pos:   rl.NewVector2((screenWidth-platformWidth)/2-150, screenHeight-450),
		Size:  rl.NewVector2(platformWidth/2-100, 30),
		Color: rl.White,
	}

	platform3 := Platform{
		Pos:   rl.NewVector2((screenWidth-platformWidth)/2+650, screenHeight-450),
		Size:  rl.NewVector2(platformWidth/2-100, 30),
		Color: rl.White,
	}

	platform4 := Platform{
		Pos:   rl.NewVector2((screenWidth-platformWidth)/2+225, screenHeight-650),
		Size:  rl.NewVector2(platformWidth/2-50, 30),
		Color: rl.White,
	}

	//mainPlatform := Platform{Pos: rl.NewVector2(screenWidth/2, 700), Size: rl.NewVector2(platformWidth, 30), Color: rl.White}

	player1 := NewPlayer(1)
	gravity := rl.NewVector2(0, 980)
	player2 := NewPlayer(2)
	player2.Transform.Pos = rl.NewVector2(1300, 400)
	player1HealthBar := NewProgressBar(295, 5, 500, 50, 5, 5)
	player2HealthBar := NewProgressBar(805, 5, 500, 50, 5, 5)
	// Set background scale for resizing
	backgroundScale := float32(screenHeight) / float32(backgroundTexture.Height)

	// Scrolling variables for each parallax layer
	var scrollingBack, scrollingMid, scrollingFore float32 = 0.0, 0.0, 0.0

	for !rl.WindowShouldClose() {
		if !player1.Alive || !player2.Alive {
			rl.ClearBackground(rl.Black)
			rl.BeginDrawing()
			rl.DrawText("KO!", 200, 200, 200, rl.Red)
			if !player2.Alive {
				rl.DrawText("Player 1 Wins!", screenWidth/2, screenHeight/2, 100, rl.White)
			} else {
				rl.DrawText("Player 2 Wins!", screenWidth/2, screenHeight/2, 100, rl.White)
			}
			rl.DrawText("Press R to restart!", screenWidth/2, (screenHeight/2)+200, 50, rl.White)
			if rl.IsKeyPressed(rl.KeyR) {
				player1.Transform.Pos = rl.NewVector2(300, 400)
				player2.Transform.Pos = rl.NewVector2(1300, 400)
				player1.Vel = rl.Vector2Zero()
				player2.Vel = rl.Vector2Zero()
				player1.Health = 100
				player2.Health = 100
				player1.Alive = true
				player2.Alive = true
				player1.ChangeState(IDLESTATE)
				player2.ChangeState(IDLESTATE)
			}
			rl.EndDrawing()
			continue
		}
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

		// Draw the platforms
		platform1.DrawPlatform()
		platform2.DrawPlatform()
		platform3.DrawPlatform()
		platform4.DrawPlatform()

		// mainPlatform.DrawPlatform()

		CheckCollision(&player1.Box, platform1)
		CheckCollision(&player1.Box, platform2)
		CheckCollision(&player1.Box, platform3)
		CheckCollision(&player1.Box, platform4)
		CheckCollision(&player2.Box, platform1)
		CheckCollision(&player2.Box, platform2)
		CheckCollision(&player2.Box, platform3)
		CheckCollision(&player2.Box, platform4)

		// CheckCollision(&player1.Box, mainPlatform)
		// CheckCollision(&player2.Box, mainPlatform)
		// Update and draw the player
		if rl.IsKeyPressed(rl.KeyQ) {
			player1.Damage(10)
			player2.Damage(10)
		}

		player1.UpdatePlayer(gravity, screenWidth, &player2)
		player2.UpdatePlayer(gravity, screenWidth, &player1)

		player1HealthBar.SetProgress(float32(player1.Health) / 100)
		player2HealthBar.SetProgress(float32(player2.Health) / 100)
		player1HealthBar.DrawBarReverse()
		player2HealthBar.DrawBar()

		rl.EndDrawing()
	}

	rl.CloseWindow()
}

// Helper function to draw seamless parallax layer with scaling
func drawSeamlessLayer(texture rl.Texture2D, positionX float32, scale float32) {
	// Draw the texture twice to ensure continuous scrolling effect
	rl.DrawTextureEx(texture, rl.NewVector2(positionX, 0), 0, scale, rl.White)
	rl.DrawTextureEx(texture, rl.NewVector2(positionX+float32(texture.Width)*scale, 0), 0, scale, rl.White)
}
