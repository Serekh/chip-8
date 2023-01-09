package utils

// Add8 returns the sum with carry of x, y and carry: sum = x + y + carry.
// The carry input must be 0 or 1; otherwise the behavior is undefined.
// The carryOut output is guaranteed to be 0 or 1.
//
// This function's execution time does not depend on the inputs.
func Add8(x, y, carry uint8) (sum, carryOut uint8) {
	sum16 := uint16(x) + uint16(y) + uint16(carry)
	sum = uint8(sum16)
	carryOut = uint8(sum16 >> 8)
	return
}

// Sub8 returns the difference of x, y and borrow, diff = x - y - borrow.
// The borrow input must be 0 or 1; otherwise the behavior is undefined.
// The borrowOut output is guaranteed to be 0 or 1.
//
// This function's execution time does not depend on the inputs.
func Sub8(x, y, borrow uint8) (diff, borrowOut uint8) {
	diff = x - y - borrow
	// The difference will underflow if the top bit of x is not set and the top
	// bit of y is set (^x & y) or if they are the same (^(x ^ y)) and a borrow
	// from the lower place happens. If that borrow happens, the result will be
	// 1 - 1 - 1 = 0 - 0 - 1 = 1 (& diff).
	borrowOut = ((^x & y) | (^(x ^ y) & diff)) >> 7
	return
}
