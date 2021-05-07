package maze

import "fmt"

// Validates the maze.
// If the cell is open to the NORTH, the cell on its NORTH must be open to the SOUTH,
// if the cell is open to the SOUTH, the cell on its SOUTH must be open to the NORTH.
// if the cell is open to the EAST, the cell on its EAST must be open to the WEST, and
// if the cell is open to the WEST, the cell on its WEST must be open to the EAST,
func (maze *Maze) ValidateMaze() {
	fmt.Printf("Validating Maze...\n")
	var point Point
	for y := 0; y < maze.Height; y++ {
		for x := 0; x < maze.Width; x++ {
			point = Point{
				X: x,
				Y: y,
			}
			// If cell is open UP
			if maze.IsWallOpen(&point, Up) {
				next := point.Advance(Up)
				if maze.Contains(next) {
					// Up cell should be open down
					if !maze.IsWallOpen(next, Down) {
						panic("cells did not match")
					}
				}
			}

			// If cell is open DOWN
			if maze.IsWallOpen(&point, Down) {
				next := point.Advance(Down)
				if maze.Contains(next) {
					// Down cell should be open Up
					if !maze.IsWallOpen(next, Up) {
						panic("cells did not match")
					}
				}
			}

			// If cell is open LEFT
			if maze.IsWallOpen(&point, Left) {
				next := point.Advance(Left)
				if maze.Contains(next) {
					// left cell should be open Right
					if !maze.IsWallOpen(next, Right) {
						panic("cells did not match")
					}
				}
			}

			// If cell is open RIGHT
			if maze.IsWallOpen(&point, Right) {
				next := point.Advance(Right)
				if maze.Contains(next) {
					// right cell should be open Left
					if !maze.IsWallOpen(next, Left) {
						panic("cells did not match")
					}
				}
			}
		}
	}
	fmt.Printf("Finished...\n")
	fmt.Printf("Maze is valid.\n")
}
