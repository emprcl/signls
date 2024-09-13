package ui

type viewport struct {
	offsetX, offsetY int
	Width, Height    int
}

func (v *viewport) Update(cursorX, cursorY, gridWidth, gridHeight int) {
	if cursorX < v.offsetX {
		v.offsetX = cursorX
	} else if cursorX >= v.offsetX+v.Width {
		v.offsetX = cursorX - v.Width + 1
	}

	if cursorY < v.offsetY {
		v.offsetY = cursorY
	} else if cursorY >= v.offsetY+v.Height {
		v.offsetY = cursorY - v.Height + 1
	}

	if v.offsetX < 0 {
		v.offsetX = 0
	} else if v.offsetX > gridWidth-v.Width {
		v.offsetX = gridWidth - v.Width
	}

	if v.offsetY < 0 {
		v.offsetY = 0
	} else if v.offsetY > gridHeight-v.Height {
		v.offsetY = gridHeight - v.Height
	}
}
