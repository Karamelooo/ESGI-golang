package game

import (
	"math/rand"
	"time"
)

const BoardSize = 10

type Board struct {
	Grid         [BoardSize][BoardSize]CellState
	Ships        []*Ship
	ReceivedHits []HitResponse
	Messages     []string
}

type Ship struct {
	Name   string
	Length int
	Hits   int
	Coords []Coord
}

func NewBoard() *Board {
	b := &Board{
		Ships: []*Ship{},
	}
	b.placeShips()
	return b
}

func (b *Board) placeShips() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	shipsToPlace := []struct {
		Name   string
		Length int
	}{
		{"Carrier", 5},
		{"Battleship", 4},
		{"Cruiser", 3},
		{"Submarine", 3},
		{"Destroyer", 2},
	}

	for _, s := range shipsToPlace {
		placed := false
		for !placed {
			horizontal := r.Intn(2) == 0
			x := r.Intn(BoardSize)
			y := r.Intn(BoardSize)

			if b.canPlace(x, y, s.Length, horizontal) {
				newShip := &Ship{
					Name:   s.Name,
					Length: s.Length,
					Coords: []Coord{},
				}
				for i := 0; i < s.Length; i++ {
					cx, cy := x, y
					if horizontal {
						cx += i
					} else {
						cy += i
					}
					b.Grid[cy][cx] = ShipCell
					newShip.Coords = append(newShip.Coords, Coord{X: cx, Y: cy})
				}
				b.Ships = append(b.Ships, newShip)
				placed = true
			}
		}
	}
}

func (b *Board) canPlace(x, y, length int, horizontal bool) bool {
	for i := 0; i < length; i++ {
		cx, cy := x, y
		if horizontal {
			cx += i
		} else {
			cy += i
		}

		if cx >= BoardSize || cy >= BoardSize {
			return false
		}
		if b.Grid[cy][cx] != Empty {
			return false
		}
	}
	return true
}

func (b *Board) ReceiveHit(c Coord) string {
	if c.X < 0 || c.X >= BoardSize || c.Y < 0 || c.Y >= BoardSize {
		return "Invalid"
	}

	currentState := b.Grid[c.Y][c.X]
	result := "Rate"

	if currentState == ShipCell {
		b.Grid[c.Y][c.X] = Hit
		result = "Touche"
		for _, s := range b.Ships {
			for _, sc := range s.Coords {
				if sc.X == c.X && sc.Y == c.Y {
					s.Hits++
					if s.Hits == s.Length {
						result = "Coule"
					}
					break
				}
			}
		}
	} else if currentState == Empty {
		b.Grid[c.Y][c.X] = Miss
		result = "Rate"
	} else {
		if currentState == Hit {
			result = "Touche"
		} else {
			result = "Rate"
		}
	}

	b.ReceivedHits = append(b.ReceivedHits, HitResponse{
		Status: result,
		Point:  c,
	})

	return result
}

func (b *Board) GetPublicGrid() [BoardSize][BoardSize]int {
	var publicGrid [BoardSize][BoardSize]int
	for y := 0; y < BoardSize; y++ {
		for x := 0; x < BoardSize; x++ {
			s := b.Grid[y][x]
			if s == ShipCell {
				publicGrid[y][x] = int(Empty)
			} else {
				publicGrid[y][x] = int(s)
			}
		}
	}
	return publicGrid
}

func (b *Board) ShipsAfloat() int {
	count := 0
	for _, s := range b.Ships {
		if s.Hits < s.Length {
			count++
		}
	}
	return count
}
