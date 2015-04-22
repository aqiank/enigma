package bytesutil

// Shift bytes by certain position (currently only one way)
func Shift(b []byte, n int) {
	if len(b) == 0 || n == 0 {
		return
	}
	var start, end, dir int
	if n > 0 {
		start = 0
		end = len(b) - 1
		dir = 1
	} else {
		start = len(b) - 1
		end = 0
		dir = -1
	}
	last := b[start]
	for i := start; i < end; i += dir {
		b[i] = b[i+n]
	}
	b[end] = last
}

