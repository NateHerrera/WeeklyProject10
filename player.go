package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Box struct {
	*Transform
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

type Player struct {
	*Transform
	Box
	StateMachine
	Health   int
	PlayerID int
}

func NewPlayer(numPlayer int) Player {
	var playerIdle rl.Texture2D
	var playerWalk rl.Texture2D
	var playerJump rl.Texture2D
	var playerBlock rl.Texture2D
	var playerAttack rl.Texture2D
	var playerHeavy rl.Texture2D
	switch numPlayer {
	case 1:
		playerIdle = rl.LoadTexture("./assets/firstplayerright.png")
		playerWalk = rl.LoadTexture("./assets/firstplayerrunright.png")
		playerJump = rl.LoadTexture("./assets/firstplayerjump.png")
		playerBlock = rl.LoadTexture("./assets/firstplayerblock.png")
		playerAttack = rl.LoadTexture("./assets/firstplayerpunch.png")
		playerHeavy = rl.LoadTexture("./assets/firstpersonheavy.png")
	case 2:
		playerIdle = rl.LoadTexture("./assets/secondplayerright.png")
		playerWalk = rl.LoadTexture("./assets/secondplayerrunright.png")
		playerJump = rl.LoadTexture("./assets/secondplayerjump.png")
		playerBlock = rl.LoadTexture("./assets/secondplayerblock.png")
		playerAttack = rl.LoadTexture("./assets/secondplayerpunch.png")
		playerHeavy = rl.LoadTexture("./assets/secondpersonheavy.png")
	}

	playerTransform := NewTransform(rl.NewVector2(300, 400))

	// Set initial Flip based on player ID
	if numPlayer == 1 {
		playerTransform.Flip = 1
	} else {
		playerTransform.Flip = -1
	}

	idleState := NewAnimation(playerTransform, playerIdle, .2, IDLESTATE)
	walkState := NewAnimation(playerTransform, playerWalk, .2, WALKSTATE)
	jumpState := NewAnimation(playerTransform, playerJump, .2, JUMPSTATE)
	blockState := NewAnimation(playerTransform, playerBlock, .1, BLOCKSTATE)
	attackState := NewAnimation(playerTransform, playerAttack, .1, ATTACKSTATE)
	heavyState := NewAnimation(playerTransform, playerHeavy, .2, HEAVYSTATE)
	playerStateMachine := NewStateMachine(&idleState)
	playerStateMachine.AddState(&walkState)
	playerStateMachine.AddState(&jumpState)
	playerStateMachine.AddState(&blockState)
	playerStateMachine.AddState(&attackState)
	playerStateMachine.AddState(&heavyState)

	return Player{
		Transform:    playerTransform,
		Box:          Box{Transform: playerTransform, Size: rl.NewVector2(128, 128), Color: rl.Red},
		StateMachine: playerStateMachine,
		PlayerID:     numPlayer,
	}
}

func (p *Player) HandlePlayer() {
	if p.PlayerID == 1 { // Controls for Player 1
		if rl.IsKeyDown(rl.KeyD) {
			if p.CurrentState.GetName() == IDLESTATE {
				p.ChangeState(WALKSTATE)
			}
			p.Box.Vel.X = 500
			p.Flip = 1
		} else if rl.IsKeyDown(rl.KeyA) {
			if p.CurrentState.GetName() == IDLESTATE {
				p.ChangeState(WALKSTATE)
			}
			p.Box.Vel.X = -500
			p.Flip = -1
		} else {
			p.Box.Vel.X = 0
			if p.CurrentState.GetName() != JUMPSTATE {
				p.ChangeState(IDLESTATE)
			}
		}
		if rl.IsKeyPressed(rl.KeyW) {
			p.ChangeState(JUMPSTATE)
			p.Box.Vel.Y = -500
		} else if rl.IsKeyDown(rl.KeyS) {
			p.PerformBlock()
		} else if rl.IsKeyDown(rl.KeyE) {
			p.PerformAttack()
		} else if rl.IsKeyDown(rl.KeyF) {
			p.PerformHeavy()
		}
	} else if p.PlayerID == 2 { // Controls for Player 2
		if rl.IsKeyDown(rl.KeyRight) {
			p.ChangeState(WALKSTATE)
			p.Transform.Pos.X += 4
			p.Flip = 1
		} else if rl.IsKeyDown(rl.KeyLeft) {
			p.ChangeState(WALKSTATE)
			p.Transform.Pos.X -= 4
			p.Flip = -1
		} else if rl.IsKeyDown(rl.KeyUp) {
			p.ChangeState(JUMPSTATE)
			p.Transform.Pos.Y -= 10
		} else if rl.IsKeyDown(rl.KeyDown) {
			p.PerformBlock()
		} else if rl.IsKeyDown(rl.KeySpace) {
			p.PerformAttack()
		} else if rl.IsKeyDown(rl.KeySlash) {
			p.PerformHeavy()
		} else {
			p.ChangeState(IDLESTATE)
		}
	}
}

// Now we need to implement the PerformAttack and PerformBlock functions
// To flip the character whoich ever way theyre facing
func (p *Player) PerformAttack() {
	// Set attack offset based on current Flip direction
	attackOffset := rl.NewVector2(30*float32(p.Flip), 0) // Use Flip as set in NewPlayer

	// Calculate the attack position based on offset
	attackPos := rl.Vector2Add(p.Transform.Pos, attackOffset)

	// Draw a hitbox for the attack (visual cue)
	rl.DrawRectangle(int32(attackPos.X), int32(attackPos.Y), 20, 10, rl.Red)

	p.ChangeState(ATTACKSTATE)
}

// same for heavy attack
func (p *Player) PerformHeavy() {
	// Set heavy attack offset based on current Flip direction
	attackOffset := rl.NewVector2(40*float32(p.Flip), 0) // Slightly larger offset for heavy attack

	// Calculate the attack position based on offset
	attackPos := rl.Vector2Add(p.Transform.Pos, attackOffset)

	// Draw a hitbox for the heavy attack (visual cue)
	rl.DrawRectangle(int32(attackPos.X), int32(attackPos.Y), 30, 15, rl.DarkGreen) // Larger and darker hitbox for heavy attack

	p.ChangeState(HEAVYSTATE)
}

// Same for the block
func (p *Player) PerformBlock() {
	// Show block position based on Flip
	blockOffset := rl.NewVector2(15, 0) // Offset for block positioning
	if p.Flip == -1 {                   // Facing left
		blockOffset.X = -15
	}
	blockPos := rl.Vector2Add(p.Transform.Pos, blockOffset)

	// Draw the block hitbox or stance
	rl.DrawRectangle(int32(blockPos.X), int32(blockPos.Y), 25, 15, rl.Blue) // Visual cue for block

	p.ChangeState(BLOCKSTATE)
}

func (p *Player) UpdatePlayer(g rl.Vector2, screenWidth float32) {
	p.ApplyGravity(g)
	p.HandlePlayer()

	// Ensure player stays within screen boundaries
	if p.Transform.Pos.X < 0 {
		p.Transform.Pos.X = 0
	} else if p.Transform.Pos.X+p.Scale.X > screenWidth {
		p.Transform.Pos.X = screenWidth - p.Scale.X
	}
	p.DrawBox()
	p.Tick()
	p.UpdateBox()

}
