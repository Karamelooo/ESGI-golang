package client

import (
	"battleship/game"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	Targets []string
}

func NewClient(targets []string) *Client {
	return &Client{
		Targets: targets,
	}
}

func (c *Client) GetOpponentBoard(targetURL string) ([game.BoardSize][game.BoardSize]int, error) {
	resp, err := http.Get(targetURL + "/board")
	var grid [game.BoardSize][game.BoardSize]int
	if err != nil {
		return grid, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return grid, fmt.Errorf("bad status: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&grid); err != nil {
		return grid, err
	}
	return grid, nil
}

func (c *Client) FireShot(targetURL string, x, y int) (game.HitResponse, error) {
	req := game.HitRequest{
		Coord: game.Coord{X: x, Y: y},
	}
	body, _ := json.Marshal(req)

	resp, err := http.Post(targetURL+"/hit", "application/json", bytes.NewBuffer(body))
	var hitResp game.HitResponse
	if err != nil {
		return hitResp, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return hitResp, fmt.Errorf("bad status: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&hitResp); err != nil {
		return hitResp, err
	}
	return hitResp, nil
}

func (c *Client) SendChat(targetURL, text string) error {
	msg := struct {
		Text string `json:"text"`
	}{Text: text}
	body, _ := json.Marshal(msg)
	resp, err := http.Post(targetURL+"/chat", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (c *Client) GetOpponentBoats(targetURL string) (int, error) {
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Get(targetURL + "/boats")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result map[string]int
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}
	return result["boats_afloat"], nil
}
