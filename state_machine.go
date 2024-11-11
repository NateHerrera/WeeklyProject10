package main

import (
    rl "github.com/gen2brain/raylib-go/raylib"
)

type State interface {
    TickState()
    GetName() string
    ResetTime()
    GetFrameIndex() int
}

type StateMachine struct {
    CurrentState State
    StateMap     map[string]State
    StateTimer   float32 // Timer to track duration in the current state
}

func NewStateMachine(initialState State) StateMachine {
    newMachine := StateMachine{
        CurrentState: initialState,
        StateMap:     make(map[string]State),
        StateTimer:   0, // Initialize the timer
    }
    newMachine.AddState(initialState)
    return newMachine
}

func (sm *StateMachine) AddState(newState State) {
    sm.StateMap[newState.GetName()] = newState
}

func (sm *StateMachine) ChangeState(stateName string) {
    if sm.CurrentState.GetName() == stateName {
        return
    }

    sm.CurrentState.ResetTime()
    sm.CurrentState = sm.StateMap[stateName]
    sm.StateTimer = 0 // Reset the timer when changing states
    sm.CurrentState.TickState()
}

func (sm *StateMachine) Tick() {
    sm.StateTimer += float32(rl.GetFrameTime()) // Increment timer with frame time
    sm.CurrentState.TickState()
}
