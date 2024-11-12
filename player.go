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
	Health              int
	PlayerID            int
	Alive               bool
	PunchHitbox         rl.Rectangle
	attackCooldown      float32
	heavyAttackCooldown float32
	lastAttackTime      float32
	lastHeavyAttackTime float32
	hasDealtDamage      bool
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
		Transform:           playerTransform,
		Box:                 Box{Transform: playerTransform, Size: rl.NewVector2(128, 128), Color: rl.Red},
		StateMachine:        playerStateMachine,
		Health:              100,
		Alive:               true,
		PlayerID:            numPlayer,
		attackCooldown:      0.5, // 0.5 seconds for normal attacks
		heavyAttackCooldown: 1.5, // 1.5 seconds for heavy attacks
		lastAttackTime:      -1,  // Initially set to -1 to indicate no recent attacks
		lastHeavyAttackTime: -1,
		hasDealtDamage:      false,
	}
}

func (p *Player) Damage(attackVal int) {
	damageTaken := attackVal

	switch p.CurrentState.GetName() {
	case BLOCKSTATE:
		damageTaken = damageTaken / 2
	case HEAVYSTATE:
		damageTaken = (damageTaken * 3) / 2
	case JUMPSTATE:
		damageTaken = damageTaken * 2
	}

	p.Health -= damageTaken

	if p.Health <= 0 {
		p.Health = 0
		p.Alive = false
	}

}

func (p *Player) HandlePlayer(p2 *Player) {
	if p.CurrentState.GetName() == ATTACKSTATE {
		if p.CurrentState.GetFrameIndex() <= 4 {
			return
		} else if p.CurrentState.GetFrameIndex() <= 8 {
			p.PerformAttack(p2)
			return
		} else {
			p.ChangeState(IDLESTATE)
		}
	}
	if p.CurrentState.GetName() == HEAVYSTATE {
		if p.CurrentState.GetFrameIndex() <= 4 {
			return
		} else if p.CurrentState.GetFrameIndex() <= 8 {
			p.PerformHeavy(p2)
			return
		} else {
			p.ChangeState(IDLESTATE)
		}
	}
	if p.CurrentState.GetName() == JUMPSTATE {
		if p.CurrentState.GetFrameIndex() >= 4 {
			p.ChangeState(IDLESTATE)
		}
	}
	if p.PlayerID == 1 { // Controls for Player 1
		if rl.IsKeyDown(rl.KeyD) {
			if p.CurrentState.GetName() == IDLESTATE {
				p.ChangeState(WALKSTATE)
			}
			p.Transform.Pos.X += 4
			p.Flip = 1
		} else if rl.IsKeyDown(rl.KeyA) {
			if p.CurrentState.GetName() == IDLESTATE {
				p.ChangeState(WALKSTATE)
			}
			p.Transform.Pos.X -= 4
			p.Flip = -1
		} else if rl.IsKeyDown(rl.KeyS) {
			p.PerformBlock()
			return
		} else {
			p.Box.Vel.X = 0
			if p.CurrentState.GetName() != JUMPSTATE {
				p.ChangeState(IDLESTATE)
			}
		}
		if rl.IsKeyPressed(rl.KeyW) {
			if p.CurrentState.GetName() != JUMPSTATE {
				p.ChangeState(JUMPSTATE)
				p.Box.Vel.Y = -500
			}
		} else if rl.IsKeyPressed(rl.KeyE) {
			p.ChangeState(ATTACKSTATE)
		} else if rl.IsKeyPressed(rl.KeyF) {
			p.ChangeState(HEAVYSTATE)
		}
	} else if p.PlayerID == 2 { // Controls for Player 2
		if rl.IsKeyDown(rl.KeyRight) {
			if p.CurrentState.GetName() == IDLESTATE {
				p.ChangeState(WALKSTATE)
			}
			p.Transform.Pos.X += 4
			p.Flip = 1
		} else if rl.IsKeyDown(rl.KeyLeft) {
			if p.CurrentState.GetName() == IDLESTATE {
				p.ChangeState(WALKSTATE)
			}
			p.Transform.Pos.X -= 4
			p.Flip = -1
		} else if rl.IsKeyDown(rl.KeyDown) {
			p.PerformBlock()
			return
		} else {
			p.Box.Vel.X = 0
			if p.CurrentState.GetName() != JUMPSTATE {
				p.ChangeState(IDLESTATE)
			}
		}

		if rl.IsKeyPressed(rl.KeyUp) {
			if p.CurrentState.GetName() != JUMPSTATE {
				p.ChangeState(JUMPSTATE)
				p.Box.Vel.Y = -500
			}
		} else if rl.IsKeyPressed(rl.KeySpace) {
			p.ChangeState(ATTACKSTATE)
		} else if rl.IsKeyPressed(rl.KeySlash) {
			p.ChangeState(HEAVYSTATE)
		}
	}
}

// Now we need to implement the PerformAttack and PerformBlock functions
// To flip the character whoich ever way theyre facing
func (p *Player) PerformAttack(p2 *Player) {
	// Check if enough time has passed since the last attack
	// if p.StateMachine.StateTimer < p.attackCooldown {
	// 	return // Exit if still in cooldown
	// }

	// Set attack offset based on current Flip direction
	attackOffset := rl.NewVector2(30*float32(p.Flip), 0)

	// Calculate the attack position based on offset
	attackPos := rl.Vector2Add(rl.Vector2Add(p.Transform.Pos, rl.NewVector2(p.Scale.X/2, p.Scale.Y/2)), attackOffset)
	p.PunchHitbox = rl.NewRectangle(attackPos.X, attackPos.Y, 20, 10)

	// Check if this attack hitbox overlaps with the opponent's hitbox
	if !p.hasDealtDamage && rl.CheckCollisionRecs(p.PunchHitbox, rl.NewRectangle(p2.Transform.Pos.X, p2.Transform.Pos.Y, p2.Scale.X, p2.Scale.Y)) {
		p2.Damage(10)
		knockbackForce := 5 * float32(p.Flip)
		p2.Transform.Pos.X += knockbackForce
		p.hasDealtDamage = true // Ensure only one damage per attack instance
	}

	p.ChangeState(ATTACKSTATE)
}

func (p *Player) CheckHit(hitBox rl.Rectangle) bool {
	return rl.CheckCollisionRecs(rl.NewRectangle(p.Pos.X, p.Pos.Y, p.Scale.X, p.Scale.Y), hitBox)
}

// same for heavy attack
func (p *Player) PerformHeavy(p2 *Player) {
	// Check if enough time has passed since the last heavy attack
	// if p.StateMachine.StateTimer < p.heavyAttackCooldown {
	// 	return // Exit if still in cooldown
	// }

	// Set heavy attack offset based on current Flip direction
	attackOffset := rl.NewVector2(40*float32(p.Flip), 0)

	// Calculate the attack position based on offset
	attackPos := rl.Vector2Add(rl.Vector2Add(p.Transform.Pos, rl.NewVector2(p.Scale.X/2, p.Scale.Y/2)), attackOffset)
	p.PunchHitbox = rl.NewRectangle(attackPos.X, attackPos.Y, 30, 15)

	// Check if this heavy attack hitbox overlaps with the opponent's hitbox
	if !p.hasDealtDamage && rl.CheckCollisionRecs(p.PunchHitbox, rl.NewRectangle(p2.Transform.Pos.X, p2.Transform.Pos.Y, p2.Scale.X, p2.Scale.Y)) {
		p2.Damage(20)
		knockbackForce := 10 * float32(p.Flip)
		p2.Transform.Pos.X += knockbackForce
		p.hasDealtDamage = true // Ensure only one damage per attack instance
	}

	p.ChangeState(HEAVYSTATE)
}

// Same for the block
func (p *Player) PerformBlock() {
	// Show block position based on Flip
	blockOffset := rl.NewVector2(15, 0) // Offset for block positioning
	if p.Flip == -1 {                   // Facing left
		blockOffset.X = -15
	}
	blockPos := rl.Vector2Add(rl.Vector2Add(p.Transform.Pos, rl.NewVector2(p.Scale.X/2, p.Scale.Y/2)), blockOffset)

	// Draw the block hitbox or stance
	rl.DrawRectangle(int32(blockPos.X), int32(blockPos.Y), 25, 15, rl.Blue) // Visual cue for block

	p.ChangeState(BLOCKSTATE)
}

func (p *Player) UpdatePlayer(g rl.Vector2, screenWidth float32, p2 *Player) {
	p.ApplyGravity(g)
	p.HandlePlayer(p2)

	if p.Transform.Pos.X < 0 {
		p.Transform.Pos.X = 0
	} else if p.Transform.Pos.X+p.Scale.X > screenWidth {
		p.Transform.Pos.X = screenWidth - p.Scale.X
	}

	// Reset the single-hit flag at the end of each attack animation
	if p.CurrentState.GetName() != ATTACKSTATE && p.CurrentState.GetName() != HEAVYSTATE {
		p.hasDealtDamage = false
	}

	p.StateMachine.Tick()
	p.UpdateBox()

	if p.Transform.Pos.Y >= screenHeight {
		p.Health = 0
		p.Alive = false
	}
}

func (p *Player) EnableAttackHitbox() {
	p.PunchHitbox = rl.NewRectangle((p.Transform.Pos.X+(p.Scale.X/2))+30*float32(p.Flip), (p.Transform.Pos.Y + (p.Scale.Y / 2)), 50, 20)
}

func (p *Player) DisableAttackHitbox() {
	p.PunchHitbox = rl.NewRectangle(0, 0, 0, 0) // Reset hotbox
}
