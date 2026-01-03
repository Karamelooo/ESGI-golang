package server

import (
	"battleship/game"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type GameServer struct {
	Board      *game.Board
	Port       string
	PlayerName string
	Mux        *http.ServeMux
}

func (s *GameServer) StartServer() {
	s.Mux = http.NewServeMux()
	s.Mux.HandleFunc("/board", s.handleGetBoard)
	s.Mux.HandleFunc("/boats", s.handleGetBoats)
	s.Mux.HandleFunc("/hit", s.handleHit)
	s.Mux.HandleFunc("/hits", s.handleGetHits)
	s.Mux.HandleFunc("/chat", s.handleChat)
	s.Mux.HandleFunc("/profile", s.handleGetProfile)

	go s.startKraken()

	s.Board.Messages = append(s.Board.Messages, fmt.Sprintf("Serveur en écoute sur le port %s...", s.Port))
	if err := http.ListenAndServe(":"+s.Port, s.Mux); err != nil {
		s.Board.Messages = append(s.Board.Messages, fmt.Sprintf("Erreur démarrage serveur: %v", err))
	}
}

func (s *GameServer) startKraken() {
	ticker := time.NewTicker(60 * time.Second)
	for range ticker.C {
		x := rand.Intn(game.BoardSize)
		y := rand.Intn(game.BoardSize)
		s.Board.ReceiveHit(game.Coord{X: x, Y: y})
		s.Board.Messages = append(s.Board.Messages, fmt.Sprintf("ATTAQUE DU KRAKEN en %d,%d !", x, y))
	}
}

func (s *GameServer) handleGetBoard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	grid := s.Board.GetPublicGrid()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(grid)
}

func (s *GameServer) handleGetBoats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	count := s.Board.ShipsAfloat()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"boats_afloat": count})
}

func (s *GameServer) handleHit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req game.HitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	status := s.Board.ReceiveHit(req.Coord)
	if status == "Coule" {
		s.Board.Messages = append(s.Board.Messages, fmt.Sprintf("VOTRE NAVIRE A ETE COULE en %d,%d !", req.Coord.X, req.Coord.Y))
	}

	resp := game.HitResponse{
		Status: status,
		Point:  req.Coord,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *GameServer) handleGetHits(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.Board.ReceivedHits)
}

func (s *GameServer) handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var msg struct {
			Text string `json:"text"`
		}
		if err := json.NewDecoder(r.Body).Decode(&msg); err == nil {
			s.Board.Messages = append(s.Board.Messages, msg.Text)
		}
	} else if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s.Board.Messages)
	}
}

func (s *GameServer) handleGetProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"name": s.PlayerName})
}
