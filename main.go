package main

import (
	"cart/w4"
	"math/rand"
	"strconv"
)

var smallDungeon = []Room{
	{Id: 0, Name: "Placeholder Because I'm Lazy", isSpawnable: false},
	{Id: 1, Name: "Main Foyer", Shape: Oval, Height: 60, Width: 75, NearbyRooms: []int64{2, 3, 4, 5}},
	{Id: 2, Name: "Side Hallway", Shape: Rectangle, Height: 72, Width: 10, NearbyRooms: []int64{1, 3}},
	{Id: 3, Name: "Kitchen", Shape: Rectangle, Height: 40, Width: 45, isSpawnable: true, NearbyRooms: []int64{1, 2, 8}},
	{Id: 4, Name: "Closet", Shape: Rectangle, Height: 20, Width: 20, NearbyRooms: []int64{1}},
	{Id: 5, Name: "Staircase", Shape: Rectangle, Height: 40, Width: 15, NearbyRooms: []int64{1, 6}},
	{Id: 6, Name: "Observatory", Shape: Oval, Height: 75, Width: 75, isSpawnable: true, NearbyRooms: []int64{5, 7}},
	{Id: 7, Name: "Balcony", Shape: Rectangle, Height: 15, Width: 70, NearbyRooms: []int64{6}},
	{Id: 8, Name: "Main Hallway", Shape: Rectangle, Height: 65, Width: 30, NearbyRooms: []int64{2, 3}},
}

var r = rand.New(rand.NewSource(2))
var RandomOrder = r.Perm(len(smallDungeon) + 1)

var smiley = [8]byte{
	0b11000011,
	0b10000001,
	0b00100100,
	0b00100100,
	0b00000000,
	0b00100100,
	0b10011001,
	0b11000011,
}
var smileySize = 8

var centerX = 80
var centerY = 80

var mapOffsetX = 0
var mapOffsetY = -15

type MenuTypes bool

const (
	Rooms   MenuTypes = true
	Actions           = false
)

var selectedMenu = Rooms
var showActionText = false
var shootMode = false
var shootOption = 0
var menuOption = 0
var currentRoom = 1
var previousGamepad byte
var alive bool = true
var won bool = false

//go:export start
func start() {
	for _, i := range RandomOrder {
		if smallDungeon[i].isSpawnable {
			smallDungeon[i].containsW = true
			break
		}
	}
}

//go:export update
func update() {

	w4.PALETTE[0] = 0x6969
	w4.PALETTE[1] = 0x000
	w4.PALETTE[2] = 0xeb6b6f
	w4.PALETTE[3] = 0xf9a875
	var gamepad = *w4.GAMEPAD1
	var pressedThisFrame = gamepad & (gamepad ^ previousGamepad)
	previousGamepad = gamepad

	if !alive {
		w4.Text("death", 120, 1)
	}

	if won {
		w4.Text("win", 120, 1)
		showActionText = false
	}

	if smallDungeon[currentRoom].containsW {
		alive = false
	}

	*w4.DRAW_COLORS = 4

	w4.Text(smallDungeon[currentRoom].Name, 10, 1)

	if smallDungeon[currentRoom].isRectangular() {
		w4.Rect(
			centerX-int(smallDungeon[currentRoom].Width/2)+mapOffsetX,
			centerY-int(smallDungeon[currentRoom].Height/2)+mapOffsetY,
			smallDungeon[currentRoom].Width,
			smallDungeon[currentRoom].Height)
	} else {
		w4.Oval(
			centerX-int(smallDungeon[currentRoom].Width/2)+mapOffsetX,
			centerY-int(smallDungeon[currentRoom].Height/2)+mapOffsetY,
			smallDungeon[currentRoom].Width,
			smallDungeon[currentRoom].Height)
	}

	w4.Text("Actions", 10, 120)
	w4.Text("Shoot", 15, 130)
	w4.Text("Look", 15+45, 130)

	w4.Text("Nearby Rooms", 10, 140)
	for index, i := range smallDungeon[currentRoom].NearbyRooms {
		w4.Text(strconv.FormatInt(smallDungeon[i].Id, 10), 15+15*index, 150)
	}

	*w4.DRAW_COLORS = 3
	w4.Blit(&smiley[0],
		centerX-int(smileySize/2)+mapOffsetX,
		centerY-int(smileySize/2)+mapOffsetY,
		uint(smileySize), uint(smileySize), w4.BLIT_1BPP)

	if selectedMenu {
		w4.Text(">", 3, 140)
		for index, i := range smallDungeon[currentRoom].NearbyRooms {
			if index == menuOption {
				w4.Text(strconv.FormatInt(smallDungeon[i].Id, 10), 15+15*index, 150)
			}
		}
	} else {
		w4.Text(">", 3, 120)
		if menuOption == 0 {
			w4.Text("Shoot", 15, 130)
		} else {
			w4.Text("Look", 15+45, 130)
		}
	}

	if showActionText {
		if menuOption == 0 {
			if shootMode {
				*w4.DRAW_COLORS = 2
				w4.Text("Shooting into ", 13, 13)
				w4.Text(strconv.FormatInt(smallDungeon[smallDungeon[currentRoom].NearbyRooms[shootOption]].Id, 10), 120, 13)
			} else {
				w4.Text("Shoot?", 13, 13)
			}
		} else {
			w4.Text("Look around?", 13, 13)
		}
	}

	if pressedThisFrame&w4.BUTTON_LEFT != 0 {
		if !shootMode && menuOption > 0 {
			menuOption -= 1
		}
		if shootMode && shootOption > 0 {
			shootOption -= 1
		}
	}
	if pressedThisFrame&w4.BUTTON_RIGHT != 0 {
		if shootMode {
			if shootOption < len(smallDungeon[currentRoom].NearbyRooms)-1 {
				shootOption += 1
			}
		} else {
			if selectedMenu {
				if menuOption < len(smallDungeon[currentRoom].NearbyRooms)-1 {
					menuOption += 1
				}
			} else {
				if menuOption < 1 {
					menuOption += 1
				}
			}
		}
	}
	if pressedThisFrame&w4.BUTTON_2 != 0 {
		if shootMode {
			if smallDungeon[int(smallDungeon[currentRoom].NearbyRooms[shootOption])].containsW {
				won = true
			}
			shootMode = false
		} else {
			if selectedMenu {
				currentRoom = int(smallDungeon[currentRoom].NearbyRooms[menuOption])
				menuOption = 0
			} else {
				if menuOption == 0 {
					shootMode = true
				}
			}
		}

	}
	if pressedThisFrame&w4.BUTTON_UP != 0 {
		selectedMenu = Actions
		showActionText = true
		menuOption = 0
	}
	if pressedThisFrame&w4.BUTTON_DOWN != 0 {
		selectedMenu = Rooms
		menuOption = 0
		showActionText = false
	}
}
