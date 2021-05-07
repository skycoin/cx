package maze

import "fmt"

// PrintMaze prints the entire maze onto cli with +,-, *space*, and | characters.
func (maze *Maze) PrintMaze() {
	hWall := []byte("+---")
	hOpen := []byte("+   ")
	vWall := []byte("|   ")
	vOpen := []byte("    ")
	rightCorner := []byte("+\n")
	rightWall := []byte("|\n")
	var b []byte

	for y := 0; y < maze.Height; y++ {
		for z := 0; z < 3; z++ {
			for x := 0; x < maze.Width; x++ {
				switch z {
				case 0:
					// Top
					if y == 0 {
						// Top wall
						b = append(b, hWall...)
						// End of top
						if x == (maze.Width)-1 {
							b = append(b, rightCorner...)
						}
					}

					if y > 0 {
						if !maze.IsWallOpen(&Point{X: x, Y: y}, Up) {
							b = append(b, hWall...)
						} else {
							b = append(b, hOpen...)
						}
						// End of top
						if x == (maze.Width)-1 {
							b = append(b, rightWall...)
						}
					}

				case 1:
					// Middle
					if x == 0 {
						b = append(b, vWall...)
					}

					if !maze.IsWallOpen(&Point{X: x, Y: y}, Right) {
						// End of middle
						if x == (maze.Width)-1 {
							b = append(b, rightWall...)
						} else {
							b = append(b, vWall...)
						}
					} else {
						b = append(b, vOpen...)
					}

				case 2:
					// Bottom
					if y == (maze.Height)-1 {
						b = append(b, hWall...)
						if x == (maze.Width)-1 {
							b = append(b, rightCorner...)
						}
					}
				}
			}
		}
	}
	fmt.Print(string(b))
}
