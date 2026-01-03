package game

type CellState int

const (
	Empty CellState = 0
	Miss  CellState = 1
	Hit   CellState = 2
	ShipCell  CellState = 3
)

type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type HitRequest struct {
	Coord Coord `json:"coord"`
}

type HitResponse struct {
	Status string `json:"status"`
	Point  Coord  `json:"point"`
}
