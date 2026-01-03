package main

import (
	"battleship/client"
	"battleship/game"
	"battleship/server"
	"flag"
	"strings"
)

func main() {
	port := flag.String("port", "8080", "Port to listen on")
	peers := flag.String("peers", "", "Comma-separated list of opponent URLs (e.g. http://localhost:8081)")
	name := flag.String("name", "Player", "Player Name")
	flag.Parse()

	board := game.NewBoard()

	srv := &server.GameServer{
		Board:      board,
		Port:       *port,
		PlayerName: *name,
	}

	go srv.StartServer()

	var peerList []string
	if *peers != "" {
		peerList = strings.Split(*peers, ",")
	}

	cli := client.NewClient(peerList)

	ui := client.NewTUI(board, cli)
	ui.Start()
}
