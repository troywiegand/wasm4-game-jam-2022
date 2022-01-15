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
	{Id: 3, Name: "Kitchen", Shape: Rectangle, Height: 40, Width: 45, isSpawnable: true, NearbyRooms: []int64{1, 2}},
	{Id: 4, Name: "Closet", Shape: Rectangle, Height: 20, Width: 20, NearbyRooms: []int64{1}},
	{Id: 5, Name: "Staircase", Shape: Rectangle, Height: 40, Width: 15, NearbyRooms: []int64{1, 6}},
	{Id: 6, Name: "Observatory", Shape: Oval, Height: 75, Width: 75, isSpawnable: true, NearbyRooms: []int64{5}},
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

var menuOption = 0
var currentRoom = 1
var previousGamepad byte

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
	w4.PALETTE[1] = 0x7c3f58
	w4.PALETTE[2] = 0xeb6b6f
	w4.PALETTE[3] = 0xf9a875
	var gamepad = *w4.GAMEPAD1
	var pressedThisFrame = gamepad & (gamepad ^ previousGamepad)
	previousGamepad = gamepad

	if smallDungeon[currentRoom].containsW {
		w4.Text("death", 120, 10)
	}

	*w4.DRAW_COLORS = 4

	w4.Text(smallDungeon[currentRoom].Name, 10, 10)

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

	w4.Text("Nearby Rooms", 10, 140)
	for index, i := range smallDungeon[currentRoom].NearbyRooms {
		w4.Text(strconv.FormatInt(smallDungeon[i].Id, 10), 15+15*index, 150)
	}

	*w4.DRAW_COLORS = 3
	w4.Blit(&smiley[0],
		centerX-int(smileySize/2)+mapOffsetX,
		centerY-int(smileySize/2)+mapOffsetY,
		uint(smileySize), uint(smileySize), w4.BLIT_1BPP)

	for index, i := range smallDungeon[currentRoom].NearbyRooms {
		if index == menuOption {
			w4.Text(strconv.FormatInt(smallDungeon[i].Id, 10), 15+15*index, 150)
		}
	}
	if pressedThisFrame&w4.BUTTON_LEFT != 0 && menuOption > 0 {
		menuOption -= 1
	}
	if pressedThisFrame&w4.BUTTON_RIGHT != 0 && menuOption < len(smallDungeon[currentRoom].NearbyRooms)-1 {
		menuOption += 1
	}
	if pressedThisFrame&w4.BUTTON_2 != 0 {
		currentRoom = int(smallDungeon[currentRoom].NearbyRooms[menuOption])
		menuOption = 0
	}
}
