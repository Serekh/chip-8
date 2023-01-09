package chip8

import "log"

const (
	keysCount uint8 = 16
)

type Input struct {
	keys [keysCount]bool
}

func NewInput() *Input {
	return &Input{}
}
func (input *Input) Clear() {
	for i := range input.keys {
		input.keys[i] = false
	}
}

func (input *Input) KeyPressed(key uint8, pressed bool) {
	if key >= keysCount {
		log.Fatalf("Keys array index out of bounds, 0x%02X", key)
	}

	input.keys[key] = pressed
}

func (input *Input) IsKeyPressed(key uint8) bool {
	if key >= keysCount {
		log.Fatalf("Keys array index out of bounds, 0x%02X", key)
	}
	return input.keys[key]
}

func (input *Input) IsAnyKeyPressed() (uint8, bool) {
	for key, pressed := range input.keys {
		if pressed {
			return uint8(key), pressed
		}
	}

	return keysCount, false
}
