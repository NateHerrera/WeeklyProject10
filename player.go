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
    Health int
}

func NewPlayer(numPlayer int) Player {
    var playerIdle rl.Texture2D
    var playerWalk rl.Texture2D
    var playerJump rl.Texture2D
    var playerBlock rl.Texture2D
    var playerAttack rl.Texture2D
    switch numPlayer {
    case 1:
        playerIdle = rl.LoadTexture("./assets/firstplayerright.png")
        playerWalk = rl.LoadTexture("./assets/firstplayerrunright.png")
        playerJump = rl.LoadTexture("./assets/firstplayerjump.png")
        playerBlock = rl.LoadTexture("./assets/firstplayerright.png")
        playerAttack = rl.LoadTexture("./assets/firstplayerpunch.png")
    case 2:
        playerIdle = rl.LoadTexture("./assets/secondplayerright.png")
        playerWalk = rl.LoadTexture("./assets/secondplayerrunright.png")
        playerJump = rl.LoadTexture("./assets/secondplayerjump.png")
        playerBlock = rl.LoadTexture("./assets/secondplayerright.png")
        playerAttack = rl.LoadTexture("./assets/secondplayerpunch.png")
    }
    playerTransform := NewTransform(rl.NewVector2(300, 400))
    idleState := NewAnimation(playerTransform, playerIdle, .2, IDLESTATE)
    walkState := NewAnimation(playerTransform, playerWalk, .2, WALKSTATE)
    jumpState := NewAnimation(playerTransform, playerJump, .2, JUMPSTATE)
    blockState := NewAnimation(playerTransform, playerBlock, .2, BLOCKSTATE)
    attackState := NewAnimation(playerTransform, playerAttack, .1, ATTACKSTATE)
    playerStateMachine := NewStateMachine(&idleState)
    playerStateMachine.AddState(&walkState)
    playerStateMachine.AddState(&jumpState)
    playerStateMachine.AddState(&blockState)
    playerStateMachine.AddState(&attackState)
    return Player{
        Transform:    playerTransform,
        Box:          Box{Transform: playerTransform, Size: rl.NewVector2(64, 64), Color: rl.Red},
        StateMachine: playerStateMachine,
    }
}

func (p *Player) HandlePlayer() {
    if rl.IsKeyDown(rl.KeyD) {
        p.ChangeState(WALKSTATE)
        p.Transform.Pos.X += 4 // Move 4 pixel to the right per frame
        p.Flip = 1
    } else if rl.IsKeyDown(rl.KeyA) {
        p.ChangeState(WALKSTATE)
        p.Transform.Pos.X -= 4 // Move 4 pixel to the left per frame
        p.Flip = -1
	} else if rl.IsKeyDown(rl.KeyW) {
		p.ChangeState(JUMPSTATE)
		p.Transform.Pos.Y -= 10 // Move 4 pixel up per frame
	} else if rl.IsKeyDown(rl.KeyS) {
		p.ChangeState(BLOCKSTATE)
	} else if rl.IsKeyDown(rl.KeyE) {	
		p.ChangeState(ATTACKSTATE)	
    } else {
        p.ChangeState(IDLESTATE)
    }

}

func (p *Player) UpdatePlayer(g rl.Vector2, screenWidth float32) {
    p.ApplyGravity(g)
    p.HandlePlayer()

    // Ensure player stays within screen boundaries
    if p.Transform.Pos.X < 0 {
        p.Transform.Pos.X = 0
    } else if p.Transform.Pos.X + p.Scale.X > screenWidth {
        p.Transform.Pos.X = screenWidth - p.Scale.X
    }

    p.UpdateBox()
    p.Tick()
}