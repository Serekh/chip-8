package main

import (
	"chip-8/pkg/chip8"
	"chip-8/pkg/utils"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	"log"
	"os"
	"path"
	"path/filepath"
)

const (
	scale         int = 15
	width             = int(chip8.ScreenW) * scale
	height            = int(chip8.ScreenH) * scale
	ticksPerFrame int = 10
	beepFilename      = "beep.wav"
)

var keys = map[ebiten.Key]uint8{
	ebiten.KeyDigit1: 0x1, ebiten.KeyDigit2: 0x2, ebiten.KeyDigit3: 0x3, ebiten.KeyDigit4: 0xC,
	ebiten.KeyQ: 0x4, ebiten.KeyW: 0x5, ebiten.KeyE: 0x6, ebiten.KeyR: 0xD,
	ebiten.KeyA: 0x7, ebiten.KeyS: 0x8, ebiten.KeyD: 0x9, ebiten.KeyF: 0xE,
	ebiten.KeyZ: 0xA, ebiten.KeyX: 0x0, ebiten.KeyC: 0xB, ebiten.KeyV: 0xF,
}

type Game struct {
	Emulator *chip8.Emulator
}

func NewGame(beeper *utils.Beeper) *Game {
	return &Game{
		Emulator: chip8.NewChip8Emulator(beeper),
	}
}

func (game *Game) Update() error {
	// Quit
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}
	// Input
	for key, keyCode := range keys {
		if inpututil.IsKeyJustPressed(key) {
			game.Emulator.KeyDown(keyCode)
		} else if inpututil.IsKeyJustReleased(key) {
			game.Emulator.KeyUp(keyCode)
		}
	}
	// Tick
	for i := 0; i < ticksPerFrame; i++ {
		game.Emulator.Tick()
	}
	game.Emulator.TickTimers()
	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	//screen.Clear()
	screenBuffer := game.Emulator.GetDisplayScreen()
	for i, pixel := range screenBuffer {
		if pixel != 0 {
			// Convert our 1D array's index into a 2D (x,y) position
			x := i % int(chip8.ScreenW)
			y := i / int(chip8.ScreenW)
			// Draw a rectangle at (x,y), scaled up by our SCALE value
			ebitenutil.DrawRect(screen, float64(x*scale), float64(y*scale), float64(scale), float64(scale), color.White)
		}
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
}

func (game *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	// Check params
	if len(os.Args) < 2 {
		log.Fatal("Expected 'path to game' that will be loaded")
	}

	// Read rom
	exeDirectory, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal("Unexpected error: ", err)
	}

	filename := path.Join(exeDirectory, "roms/", os.Args[1])
	var romData []uint8
	romData, err = os.ReadFile(filename)
	if err != nil {
		log.Fatal("Unexpected error: ", err)
	}

	// Initialize beeper
	var beepSound []uint8
	beepSound, err = os.ReadFile(path.Join(exeDirectory, beepFilename))
	if err != nil {
		log.Fatal("Unexpected error: ", err)
	}
	beeper := utils.NewBeeper(beepSound)
	defer beeper.Close()

	// Initialize emulator
	game := NewGame(beeper)
	game.Emulator.Load(romData)

	// Exec
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Chip-8")
	if err = ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
