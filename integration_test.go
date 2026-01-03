package main

import (
	"battleship/client"
	"battleship/game"
	"battleship/server"
	"fmt"
	"testing"
	"time"
)

func TestGameIntegration(t *testing.T) {
	// Setup Player 1
	board1 := game.NewBoard()
	srv1 := &server.GameServer{
		Board:      board1,
		Port:       "8090", // Test Port 1
		PlayerName: "Tester1",
	}
	go srv1.StartServer()

	// Setup Player 2
	board2 := game.NewBoard()
	srv2 := &server.GameServer{
		Board:      board2,
		Port:       "8091", // Test Port 2
		PlayerName: "Tester2",
	}
	go srv2.StartServer()

	// Wait for servers to start
	time.Sleep(1 * time.Second)

	// Setup Clients
	cli1 := client.NewClient([]string{"http://localhost:8091"}) // P1 targets P2
	cli2 := client.NewClient([]string{"http://localhost:8090"}) // P2 targets P1

	t.Run("Test Shooting", func(t *testing.T) {
		// P1 shoots P2 at 0,0
		resp, err := cli1.FireShot("http://localhost:8091", 0, 0)
		if err != nil {
			t.Fatalf("P1 failed to shoot P2: %v", err)
		}
		if resp.Status != "Miss" && resp.Status != "Hit" {
			t.Errorf("Unexpected status: %s", resp.Status)
		}
		fmt.Printf("P1 received shot result: %s\n", resp.Status)

		// P2 shoots P1 at 5,5
		resp2, err := cli2.FireShot("http://localhost:8090", 5, 5)
		if err != nil {
			t.Fatalf("P2 failed to shoot P1: %v", err)
		}
		fmt.Printf("P2 received shot result: %s\n", resp2.Status)
	})

	t.Run("Test Chat", func(t *testing.T) {
		msg := "Hello P2"
		err := cli1.SendChat("http://localhost:8091", msg)
		if err != nil {
			t.Fatalf("Failed to send chat: %v", err)
		}

		// Verify P2 received it (checking server state directly for simplicity)
		found := false
		for _, m := range board2.Messages {
			if m == msg {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("P2 did not receive chat message")
		}
	})

	t.Run("Test Boats Afloat", func(t *testing.T) {
		count, err := cli1.GetOpponentBoats("http://localhost:8091")
		if err != nil {
			t.Fatalf("Failed to get boats: %v", err)
		}
		if count != 5 { // Starts with 5 ships
			t.Errorf("Expected 5 boats, got %d", count)
		}
	})
}
