package chip8

const (
	ScreenW uint16 = 64
	ScreenH uint16 = 32
)

type Display struct {
	screen [ScreenW * ScreenH]uint8
}

func NewDisplay() *Display {
	return &Display{}
}

func (display *Display) GetScreen() [ScreenW * ScreenH]uint8 {
	return display.screen
}

func (display *Display) Clear() {
	for i := range display.screen {
		display.screen[i] = 0
	}
}

func (display *Display) Draw(x, y uint16) uint8 {
	// Sprites should wrap around screen, so apply modulo
	x %= ScreenW
	y %= ScreenH
	// Get our pixel's index for our 1D screen array
	index := x + ScreenW*y
	// Return old value and
	old := display.screen[index]
	display.screen[index] ^= 1
	return old
}
