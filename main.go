package main

import (
	"cart/w4"
	"math/rand"
	"strconv"
)

var smallDungeon = []Room{
	{Id: 0, Name: "Front Porch", isSpawnable: false, NearbyRooms: []int64{1}},
	{Id: 1, Name: "Main Foyer", Description: "a ragged entryway", Shape: Oval, Height: 60, Width: 75, NearbyRooms: []int64{2, 3, 4, 5}},
	{Id: 2, Name: "Side Hallway", Description: "mesmeric tiling", Shape: Rectangle, Height: 72, Width: 10, NearbyRooms: []int64{1, 3, 8}},
	{Id: 3, Name: "Kitchen", Description: "mystery meat", Shape: Rectangle, Height: 40, Width: 45, isSpawnable: true, NearbyRooms: []int64{1, 2, 8}},
	{Id: 4, Name: "Closet", Description: "a few hung coats", Shape: Rectangle, Height: 20, Width: 20, NearbyRooms: []int64{1, 12}},
	{Id: 5, Name: "Staircase", Description: "a carved railing", Shape: Rectangle, Height: 40, Width: 15, NearbyRooms: []int64{1, 6}},
	{Id: 6, Name: "Observatory", Description: "taurus mocks you", Shape: Oval, Height: 75, Width: 75, isSpawnable: true, NearbyRooms: []int64{5, 7}},
	{Id: 7, Name: "Balcony", Description: "crickets chirp", Shape: Rectangle, Height: 15, Width: 70, NearbyRooms: []int64{6, 11}},
	{Id: 8, Name: "Main Hallway", Description: "chipping paint", Shape: Rectangle, Height: 65, Width: 30, isSpawnable: true, NearbyRooms: []int64{2, 3, 9, 11}},
	{Id: 9, Name: "Small Bedroom", Description: "a broken toy", Shape: Rectangle, Height: 42, Width: 38, isSpawnable: true, NearbyRooms: []int64{8, 10}},
	{Id: 10, Name: "Small Bathroom", Description: "a dripping sink", Shape: Rectangle, Height: 28, Width: 32, isSpawnable: true, NearbyRooms: []int64{8, 9, 12}},
	{Id: 11, Name: "Study", Description: "a dried inkwell", Shape: Oval, Height: 45, Width: 38, isSpawnable: true, NearbyRooms: []int64{7, 8}},
	{Id: 12, Name: "Crawlspace", Description: "a dirth of webs", Shape: Oval, Height: 12, Width: 80, isSpawnable: true, NearbyRooms: []int64{4, 10}},
}

var currentSeed int64 = 0
var previousWumpusRoom = -1

func randomizeDungeon(n int) {

	for _, i := range smallDungeon {
		i.containsW = false
	}

	seeds := []int64{24555, 23, -79, 48, 45, 86, -432}
	currentSeed = seeds[n%len(seeds)]*int64(n) + int64(n) + int64(previousWumpusRoom)
	rnd := rand.New(rand.NewSource(currentSeed))

	RandomOrder := rnd.Perm(len(smallDungeon))

	for _, i := range RandomOrder {
		if smallDungeon[i].isSpawnable && i != previousWumpusRoom {
			smallDungeon[i].containsW = true
			previousWumpusRoom = i
			break
		}
	}
}

func randomWumpusWarning(n int) string {

	warnings := []string{"Muffled Growls", "Ominous Feeling", "Evil Was Here", "Lingering Dread", "Foul Odors", "Sounding Snarlings"}
	rnd := rand.New(rand.NewSource(int64(n)))

	RandomOrder := rnd.Intn(len(warnings))

	return warnings[RandomOrder]
}

var smiley = [8]byte{
	0b11000011,
	0b10000001,
	0b00010010,
	0b00110110,
	0b00000000,
	0b00100100,
	0b10111101,
	0b11000011,
}

var deadSmiley = [8]byte{
	0b11000011,
	0b10000001,
	0b00100100,
	0b00010010,
	0b00000000,
	0b00000000,
	0b10000001,
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
var lookMode = false
var shootOption = 0
var menuOption = 0
var currentRoom = 1
var previousGamepad byte
var pressedThisFrame byte
var alive bool = true
var won bool = false
var showMenus bool = true
var showEnding bool = false
var gameNumber = 0

func handleRestart() bool {
	showMenus = false
	selectedMenu = Actions
	w4.Text("Play Again?", 10, 120)
	w4.Text("Yes", 15, 130)
	w4.Text("No", 15+45, 130)
	*w4.DRAW_COLORS = 3
	if menuOption == 0 {
		w4.Text(">", 8, 130)
		if pressedThisFrame&w4.BUTTON_1 != 0 {
			gameNumber += 1
			w4.Text("Restarting", 15, 150)
			randomizeDungeon(gameNumber)
			alive = true
			showMenus = true
			won = false
			selectedMenu = Rooms
			currentRoom = 1
			menuOption = -1
			showEnding = false
			return true
		}
	} else {
		w4.Text(">", 53, 130)
		if pressedThisFrame&w4.BUTTON_1 != 0 {
			showEnding = true
		}
	}
	if showEnding {
		w4.Text("Game Over Yeah!", 15, 150)
	}
	return false
}

//go:export start
func start() {
	randomizeDungeon(gameNumber)
}

//go:export update
func update() {

	w4.PALETTE[0] = 0x6969
	w4.PALETTE[1] = 0x000
	w4.PALETTE[2] = 0xeb6b6f
	w4.PALETTE[3] = 0xf9a875
	var gamepad = *w4.GAMEPAD1
	pressedThisFrame = gamepad & (gamepad ^ previousGamepad)
	previousGamepad = gamepad

	*w4.DRAW_COLORS = 4
	// w4.Text(strconv.FormatInt(currentSeed, 10), 120, 1)
	// w4.Text(strconv.FormatInt(int64(menuOption), 10), 100, 1)
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

	if !alive {
		*w4.DRAW_COLORS = 3
		w4.Blit(&deadSmiley[0],
			centerX-int(smileySize/2)+mapOffsetX,
			centerY-int(smileySize/2)+mapOffsetY,
			uint(smileySize), uint(smileySize), w4.BLIT_1BPP)
		w4.Text("death", 120, 1)
		w4.Text("The foul Wumpus \nhas slayed you!", 13, 13)
		handleRestart()
	} else {
		if won {
			*w4.DRAW_COLORS = 3
			w4.Blit(&smiley[0],
				centerX-int(smileySize/2)+mapOffsetX,
				centerY-int(smileySize/2)+mapOffsetY,
				uint(smileySize), uint(smileySize), w4.BLIT_1BPP)
			*w4.DRAW_COLORS = 4
			w4.Text("win", 120, 1)
			w4.Text("You have slayed \nthe foul Wumpus!", 13, 13)
			showActionText = false
			if handleRestart() {
				currentRoom = 1
			}
		}

		if smallDungeon[currentRoom].containsW {
			alive = false
		}

		*w4.DRAW_COLORS = 4

		w4.Text(smallDungeon[currentRoom].Name, 10, 1)

		if showMenus {
			w4.Text("Actions", 10, 120)
			w4.Text("Shoot", 15, 130)
			w4.Text("Look", 15+45, 130)

			w4.Text("Nearby Rooms", 10, 140)
			for index, i := range smallDungeon[currentRoom].NearbyRooms {
				w4.Text(strconv.FormatInt(smallDungeon[i].Id, 10), 15+15*index, 150)
			}
		}

		*w4.DRAW_COLORS = 3
		w4.Blit(&smiley[0],
			centerX-int(smileySize/2)+mapOffsetX,
			centerY-int(smileySize/2)+mapOffsetY,
			uint(smileySize), uint(smileySize), w4.BLIT_1BPP)

		if showMenus {
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
		}

		if showActionText {
			if menuOption == 0 {
				if shootMode {
					*w4.DRAW_COLORS = 2
					w4.Text("Shooting into ", 13, 13)
					*w4.DRAW_COLORS = 3
					w4.Text(strconv.FormatInt(smallDungeon[smallDungeon[currentRoom].NearbyRooms[shootOption]].Id, 10), 120, 13)
				} else {
					w4.Text("Shoot?", 13, 13)
				}
			} else {
				if lookMode {
					nearbyWump := false
					for _, i := range smallDungeon[currentRoom].NearbyRooms {
						if smallDungeon[i].containsW {
							nearbyWump = true
						}
					}
					if nearbyWump {
						w4.Text(randomWumpusWarning(gameNumber*currentRoom+currentRoom+int(smallDungeon[currentRoom].Width)), 13, 13)
					} else {
						w4.Text(smallDungeon[currentRoom].Description, 13, 13)
					}
				} else {
					w4.Text("Look around?", 13, 13)
				}
			}
		}

		if pressedThisFrame&w4.BUTTON_1 != 0 {
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
					} else {
						lookMode = true
					}
				}
			}

		}
		if pressedThisFrame&w4.BUTTON_UP != 0 {
			lookMode = false
			selectedMenu = Actions
			showActionText = true
			menuOption = 0
		}
		if pressedThisFrame&w4.BUTTON_DOWN != 0 {
			lookMode = false
			selectedMenu = Rooms
			menuOption = 0
			showActionText = false
		}
	}
	if pressedThisFrame&w4.BUTTON_LEFT != 0 {
		lookMode = false
		if !shootMode && menuOption > 0 {
			menuOption -= 1
		}
		if shootMode && shootOption > 0 {
			shootOption -= 1
		}
	}
	if pressedThisFrame&w4.BUTTON_RIGHT != 0 {
		lookMode = false
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

	if menuOption < 0 {
		menuOption = 0
	}
}
