package client

import (
	"battleship/game"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type TUI struct {
	Board  *game.Board
	Client *Client
	Target string
	Reader *bufio.Reader
}

func NewTUI(b *game.Board, c *Client) *TUI {
	t := &TUI{
		Board:  b,
		Client: c,
		Reader: bufio.NewReader(os.Stdin),
	}
	if len(c.Targets) > 0 {
		t.Target = c.Targets[0]
	}
	return t
}

func (t *TUI) Start() {
	for {
		t.render()
		t.handleInput()
	}
}

func (t *TUI) render() {
	fmt.Print("\033[H\033[2J")

	fmt.Println("=== COMMANDANT DE BATEAU DE GUERRE ===")
	fmt.Printf("Navires à flot: %d\n", t.Board.ShipsAfloat())

	fmt.Println("\nVOTRE PLATEAU:")
	t.printOwnBoard()

	fmt.Printf("\nADVERSAIRE (%s):\n", t.Target)
	if t.Target != "" {
		oppGrid, err := t.Client.GetOpponentBoard(t.Target)
		if err != nil {
			fmt.Printf("Erreur récupération plateau adverse: %v\n", err)
		} else {
			t.printGrid(oppGrid, false)
			boats, err := t.Client.GetOpponentBoats(t.Target)
			if err == nil && boats == 0 {
				fmt.Println("\n!!! VOUS AVEZ GAGNÉ !!! TOUS LES NAVIRES ENNEMIS SONT COULÉS !!!")
			}
		}
	} else {
		fmt.Println("Aucune cible sélectionnée.")
	}

	if t.Board.ShipsAfloat() == 0 {
		fmt.Println("\n!!! VOUS AVEZ PERDU !!! VOTRE FLOTTE EST DÉTRUITE !!!")
	}

	fmt.Println("\nJournal des Tirs Reçus (5 derniers):")
	cnt := len(t.Board.ReceivedHits)
	start := cnt - 5
	if start < 0 {
		start = 0
	}
	for i := start; i < cnt; i++ {
		h := t.Board.ReceivedHits[i]
		fmt.Printf("- %d,%d : %s\n", h.Point.X, h.Point.Y, h.Status)
	}
}

func (t *TUI) printOwnBoard() {
	fmt.Println("   0 1 2 3 4 5 6 7 8 9")
	fmt.Println("  ______________________ X")
	for y := 0; y < game.BoardSize; y++ {
		fmt.Printf("%d |", y)
		for x := 0; x < game.BoardSize; x++ {
			cell := t.Board.Grid[y][x]
			symbol := "~"
			switch cell {
			case game.Empty:
				symbol = "~"
			case game.Miss:
				symbol = "o"
			case game.Hit:
				symbol = "X"
			case game.ShipCell:
				symbol = "S"
			}
			fmt.Printf("%s ", symbol)
		}
		fmt.Println()
	}
	fmt.Println("Y")
}

func (t *TUI) printGrid(grid [game.BoardSize][game.BoardSize]int, isOwn bool) {
	fmt.Println("   0 1 2 3 4 5 6 7 8 9")
	fmt.Println("  ______________________ X")
	for y := 0; y < game.BoardSize; y++ {
		fmt.Printf("%d |", y)
		for x := 0; x < game.BoardSize; x++ {
			val := game.CellState(grid[y][x])
			symbol := "?"
			switch val {
			case game.Empty:
				symbol = "~"
			case game.Miss:
				symbol = "o"
			case game.Hit:
				symbol = "X"
			}
			fmt.Printf("%s ", symbol)
		}
		fmt.Println()
	}
	fmt.Println("Y")
}

func (t *TUI) handleInput() {
	if t.Board.ShipsAfloat() == 0 {
		fmt.Println("\n!!! VOUS AVEZ PERDU - GAME OVER !!!")
		fmt.Println("(Appuyez sur 'q' pour quitter)")
		fmt.Print("> ")
		input, _ := t.Reader.ReadString('\n')
		if strings.TrimSpace(input) == "q" {
			os.Exit(0)
		}
		return
	}

	fmt.Println("\nMessages Serveur / Kraken:");
	for _, m := range t.Board.Messages {
		fmt.Println(m)
	}

	fmt.Println("\nCommandes: [x y] tir, [croix x y] tir-croix, [cercle x y] tir-cercle, [ligne x1 y1 x2 y2] tir-ligne, [chat <msg>] message, [target <url>] cible, [add <url>] ajouter, [q] quitter")
	fmt.Print("> ")
	input, _ := t.Reader.ReadString('\n')
	input = strings.TrimSpace(input)

	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}

	cmd := parts[0]
	if cmd == "q" {
		os.Exit(0)
	}

	if cmd == "target" && len(parts) > 1 {
		t.Target = parts[1]
		return
	}

	if cmd == "chat" && len(parts) > 1 {
		if t.Target == "" {
			fmt.Println("Pas de cible")
			return
		}
		msg := strings.Join(parts[1:], " ")
		t.Client.SendChat(t.Target, msg)
		fmt.Println("Envoyé:", msg)
		return
	}

	if cmd == "croix" && len(parts) == 3 {
		x, _ := strconv.Atoi(parts[1])
		y, _ := strconv.Atoi(parts[2])
		coords := []struct{ x, y int }{{x, y}, {x + 1, y}, {x - 1, y}, {x, y + 1}, {x, y - 1}}
		for _, c := range coords {
			t.fireAndPrint(c.x, c.y)
			time.Sleep(200 * time.Millisecond)
		}
		return
	}

	if cmd == "cercle" && len(parts) == 3 {
		cx, _ := strconv.Atoi(parts[1])
		cy, _ := strconv.Atoi(parts[2])
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				if dx == 0 && dy == 0 {
					continue
				}
				t.fireAndPrint(cx+dx, cy+dy)
				time.Sleep(200 * time.Millisecond)
			}
		}
		return
	}

	if cmd == "ligne" && len(parts) == 5 {
		x1, _ := strconv.Atoi(parts[1])
		y1, _ := strconv.Atoi(parts[2])
		x2, _ := strconv.Atoi(parts[3])
		y2, _ := strconv.Atoi(parts[4])

		dx := x2 - x1
		dy := y2 - y1
		steps := 0
		if abs(dx) > abs(dy) {
			steps = abs(dx)
		} else {
			steps = abs(dy)
		}

		if steps > 4 {
			fmt.Println("Erreur: La longueur maximale du tir en ligne est de 4 cases.")
			return
		}

		if steps == 0 {
			steps = 1
		}

		Xinc := float64(dx) / float64(steps)
		Yinc := float64(dy) / float64(steps)

		currX := float64(x1)
		currY := float64(y1)

		for i := 0; i <= steps; i++ {
			tx := int(currX + 0.5)
			ty := int(currY + 0.5)

			t.fireAndPrint(tx, ty)
			time.Sleep(200 * time.Millisecond)

			currX += Xinc
			currY += Yinc
		}
		return
	}

	if cmd == "add" && len(parts) > 1 {
		t.Client.Targets = append(t.Client.Targets, parts[1])
		if t.Target == "" {
			t.Target = parts[1]
		}
		return
	}

	if len(parts) == 2 {
		x, err1 := strconv.Atoi(parts[0])
		y, err2 := strconv.Atoi(parts[1])
		if err1 == nil && err2 == nil {
			t.fireAndPrint(x, y)
			time.Sleep(2 * time.Second)
			return
		}
	}

	fmt.Println("Commande inconnue")
	time.Sleep(1 * time.Second)
}

func (t *TUI) fireAndPrint(x, y int) {
	if t.Target == "" {
		fmt.Println("Pas de cible!")
		return
	}

	resp, err := t.Client.FireShot(t.Target, x, y)
	if err != nil {
		fmt.Printf("Echec tir %d,%d: %v\n", x, y, err)
		return
	}

	fmt.Printf("Tir en %d,%d: %s\n", x, y, resp.Status)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
