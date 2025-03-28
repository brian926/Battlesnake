package main

// Welcome to
// __________         __    __  .__                               __
// \______   \_____ _/  |__/  |_|  |   ____   ______ ____ _____  |  | __ ____
//  |    |  _/\__  \\   __\   __\  | _/ __ \ /  ___//    \\__  \ |  |/ // __ \
//  |    |   \ / __ \|  |  |  | |  |_\  ___/ \___ \|   |  \/ __ \|    <\  ___/
//  |________/(______/__|  |__| |____/\_____>______>___|__(______/__|__\\_____>
//
// This file can be a nice home for your Battlesnake logic and helper functions.
//
// To get you started we've included code to prevent your Battlesnake from moving backwards.
// For more info see docs.battlesnake.com

import (
	"log"
	"math/rand"
	"strings"
)

// info is called when you create your Battlesnake on play.battlesnake.com
// and controls your Battlesnake's appearance
// TIP: If you open your Battlesnake URL in a browser you should see this data
func info() BattlesnakeInfoResponse {
	log.Println("INFO")

	return BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "",            // TODO: Your Battlesnake username
		Color:      "#20b2aa",     // TODO: Choose color
		Head:       "gamer",       // TODO: Choose head
		Tail:       "rbc-necktie", // TODO: Choose tail
	}
}

// start is called when your Battlesnake begins a game
func start(state GameState) {
	log.Println("GAME START")
}

// end is called when your Battlesnake finishes a game
func end(state GameState) {
	log.Printf("GAME OVER\n\n")
}

func isHeadAvoidingBody(newHead Coord, body []Coord) bool {
	for _, segment := range body {
		if newHead == segment {
			return false
		}
	}
	return true
}

// move is called on every turn and returns your next move
// Valid moves are "up", "down", "left", or "right"
// See https://docs.battlesnake.com/api/example-move for available data
func move(state GameState) BattlesnakeMoveResponse {

	isMoveSafe := map[string]bool{
		"up":    true,
		"down":  true,
		"left":  true,
		"right": true,
	}

	// We've included code to prevent your Battlesnake from moving backwards
	myHead := state.You.Body[0] // Coordinates of your head
	myNeck := state.You.Body[1] // Coordinates of your "neck"

	if myNeck.X < myHead.X { // Neck is left of head, don't move left
		isMoveSafe["left"] = false
		log.Print("LEFT isn't safe, neck is in the way")

	} else if myNeck.X > myHead.X { // Neck is right of head, don't move right
		isMoveSafe["right"] = false
		log.Print("RIGHT isn't safe, neck is in the way")

	} else if myNeck.Y < myHead.Y { // Neck is below head, don't move down
		isMoveSafe["down"] = false
		log.Print("DOWN isn't safe, neck is in the way")

	} else if myNeck.Y > myHead.Y { // Neck is above head, don't move up
		isMoveSafe["up"] = false
		log.Print("UP isn't safe, neck is in the way")
	}

	// TODO: Step 1 - Prevent your Battlesnake from moving out of bounds
	boardWidth, boardHeight := state.Board.Width, state.Board.Height // Width and Height of board e.g. 0,0 10,10
	outOfBoard := map[string]bool{
		"up":    myHead.Y+1 >= boardHeight,
		"down":  myHead.Y-1 < 0,
		"left":  myHead.X-1 < 0,
		"right": myHead.X+1 >= boardWidth,
	}

	for move, isOut := range outOfBoard {
		if isOut {
			isMoveSafe[move] = false
			log.Printf("%s isn't safe, off of board", strings.ToUpper(move))
		}
	}

	// TODO: Step 2 - Prevent your Battlesnake from colliding with itself
	myBody := state.You.Body

	moves := map[string]Coord{
		"left":     {X: myHead.X - 1, Y: myHead.Y},
		"right":    {X: myHead.X + 1, Y: myHead.Y},
		"up":       {X: myHead.X, Y: myHead.Y + 1},
		"down":     {X: myHead.X, Y: myHead.Y - 1},
		"twoleft":  {X: myHead.X - 2, Y: myHead.Y},
		"tworight": {X: myHead.X + 2, Y: myHead.Y},
		"twoup":    {X: myHead.X, Y: myHead.Y + 2},
		"twodown":  {X: myHead.X, Y: myHead.Y - 2},
	}

	two := map[string]string{
		"twoup":    "up",
		"twodown":  "down",
		"twoleft":  "left",
		"tworight": "right",
	}

	for move, coord := range moves {
		if !isHeadAvoidingBody(coord, myBody) {
			value, exists := two[move]
			if exists {
				isMoveSafe[value] = false
				log.Printf("Two %s isn't safe, body in way", strings.ToUpper(value))
			} else {
				isMoveSafe[move] = false
				log.Printf("%s isn't safe, body in way", strings.ToUpper(move))
			}
		}
	}

	// TODO: Step 3 - Prevent your Battlesnake from colliding with other Battlesnakes
	opponents := state.Board.Snakes
	for _, snakes := range opponents {
		for move, coord := range moves {
			if !isHeadAvoidingBody(coord, snakes.Body) {
				value, exists := two[move]
				if exists {
					isMoveSafe[value] = false
					log.Printf("%s isn't safe, another snake in way", strings.ToUpper(value))
				} else {
					isMoveSafe[move] = false
					log.Printf("%s isn't safe, another snake in way", strings.ToUpper(move))
				}
			}
		}
	}

	// Are there any safe moves left?
	safeMoves := []string{}
	for move, isSafe := range isMoveSafe {
		if isSafe {
			safeMoves = append(safeMoves, move)
		}
	}

	if len(safeMoves) == 0 {
		log.Printf("MOVE %d: No safe moves detected! Moving down\n", state.Turn)
		return BattlesnakeMoveResponse{Move: "down"}
	}

	// Choose a random move from the safe ones
	nextMove := safeMoves[rand.Intn(len(safeMoves))]

	// TODO: Step 4 - Move towards food instead of random, to regain health and survive longer
	// food := state.Board.Food

	log.Printf("MOVE %d: %s X:%d Y:%d\n", state.Turn, nextMove, myHead.X, myHead.Y)
	return BattlesnakeMoveResponse{Move: nextMove}
}

func main() {
	RunServer()
}
