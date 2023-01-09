package chip8

import "log"

const (
	stackSize uint16 = 16
)

type Stack struct {
	sp        uint16 // Stack Pointer
	addresses [stackSize]uint16
}

func NewStack() *Stack {
	return &Stack{}
}

func (stack *Stack) Clear() {
	stack.sp = 0x00
	for i := range stack.addresses {
		stack.addresses[i] = 0x00
	}
}

func (stack *Stack) Push(value uint16) {
	if stack.sp >= stackSize {
		log.Fatalf("Stack overflow")
	}

	stack.addresses[stack.sp] = value
	stack.sp += 1
}

func (stack *Stack) Pop() uint16 {
	if stack.sp == 0 {
		log.Fatalf("Stack overflow")
	}

	stack.sp -= 1
	return stack.addresses[stack.sp]
}
