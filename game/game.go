package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type SquareState int
const (
	SquareStateEmpty  SquareState = 0
	SquareStateCross  SquareState = 'X'
	SquareStateNaught SquareState = '0'
)

type Result int
const (
	ResultNone      Result = 0
	ResultNInARow   Result = 1
	ResultStalemate Result = 2
)

type Game struct {
	Board [][]SquareState
	Turn int
}

type TicTacToeStateRequest struct {
	Board [][]SquareState `json:"board"`
}

type TicTacToeStateResponse struct {
	Board      [][]SquareState   `json:"board"`
	Result     Result            `json:"result,omitempty"`
	WinningRow [][]SquareState   `json:"winningRow,omitempty"`
	Turn       int               `json:"turn"`
	NextPlayer rune              `json:"nextPlayer"`
}

func makeBoard(n int) [][]SquareState {
	board := make([][]SquareState, n)
	for i := 0; i < n; i++ {
		board[i] = make([]SquareState, n)
	}

	return board
}

func copyBoard(src [][]SquareState) [][]SquareState {
	n := len(src)
	board := make([][]SquareState, n)
	for j := 0; j < n; j++ {
		board[j] = make([]SquareState, n)
		copy(board[j], src[j])
	}

	return board
}

func newGameFromRequest(req *TicTacToeStateRequest) *Game {
	turn := 1
	for _, y := range req.Board {
		for _, x := range y {
			if x != SquareStateEmpty {
				turn++
			}
		}
	}

	return &Game{
		Board: req.Board,
		Turn: turn,
	}
}

func opponent(player rune) rune {
	return player ^ 'h'
}

func (g *Game) playersTurn() rune {
	if g.Turn % 2 == 1 {
		return 'X'
	}

	return '0'
}

func (g *Game) isOccupied(x, y int) bool {
	return g.Board[y][x] != SquareStateEmpty
}

func (g *Game) occupyPosition(x, y int) error {
	if x < 0 || x > 2 || y < 0 || y > 2 {
		return errors.New("invalid coordinate")
	}
	if g.Board[y][x] != SquareStateEmpty {
		return errors.New("already occupied")
	}

	player := g.playersTurn()
	g.Turn++
	if player == 'X' {
		g.Board[y][x] = SquareStateCross
		return nil
	}
	g.Board[y][x] = SquareStateNaught

	return nil
}

// getGameResult calculates the current state of the game returning the result
//  and the row that concluded the game if there is a complete row, nil otherwise.
func (g *Game) getGameResult() (Result, [][]SquareState) {
	n := len(g.Board)
	var rowOfN [][]SquareState

	// Check diagonal
	rowOfN = makeBoard(n)
	var i int
	for i = 0; i < n - 1 && g.Board[i][i] == g.Board[i+1][i+1] && g.Board[i][i] != SquareStateEmpty; i++ {
		rowOfN[i][i] = g.Board[i][i]
		rowOfN[i+1][i+1] = g.Board[i+1][i+1]
	}
	if i == n - 1 {
		return ResultNInARow, rowOfN
	}

	// Check anti-diagonal
	rowOfN = makeBoard(n)
	for i = 0; i < n - 1 && g.Board[n-1-i][i] == g.Board[n-i-2][i+1] && g.Board[n-i-2][i+1] != SquareStateEmpty; i++ {
		rowOfN[n-1-i][i] = g.Board[n-1-i][i]
		rowOfN[n-i-2][i+1] = g.Board[n-i-2][i+1]
	}
	if i == 2 {
		return ResultNInARow, rowOfN
	}

	var j int
	for j = 0; j < n; j++ {
		// Check columns
		rowOfN = makeBoard(n)
		for i = 0; i < n-1 && g.Board[i][j] == g.Board[i+1][j] && g.Board[i][j] != SquareStateEmpty; i++ {
			rowOfN[i][j] = g.Board[i][j]
			rowOfN[i+1][j] = g.Board[i+1][j]
		}
		if i == n-1 {
			return ResultNInARow, rowOfN
		}

		// Check rows
		rowOfN = makeBoard(n)
		for i = 0; i < n-1 && g.Board[j][i] == g.Board[j][i+1] && g.Board[j][i] != SquareStateEmpty; i++ {
			rowOfN[j][i] = g.Board[j][i]
			rowOfN[j][i+1] = g.Board[j][i+1]
		}
		if i == n-1 {
			return ResultNInARow, rowOfN
		}
	}

	// Check for stalemate
	if g.Turn > n * n {
		return ResultStalemate, nil
	}

	return ResultNone, nil
}

func writeHTTPError(w http.ResponseWriter, statusCode int, description string, err error) {
	message := fmt.Sprintf("%s : %v", description, err)
	_, wErr := w.Write([]byte(message))
	if wErr != nil {
		log.Fatal(wErr)
	}
	log.Fatal(description, err)
	w.WriteHeader(statusCode)
}

func computeMove(gameState Game, isMax bool) (int, int, int) {
	optimalX := 0
	optimalY := 0
	multiplier := 1
	if !isMax {
		multiplier = -1
	}
	threshold := math.MaxInt32 * -1 * multiplier

	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			if gameState.isOccupied(x, y) {
				continue
			}

			gs := Game{
				Board: copyBoard(gameState.Board),
				Turn:  gameState.Turn,
			}

			err := gs.occupyPosition(x, y)
			if err != nil {
				log.Fatal(err)
				continue
			}
			result, _ := gs.getGameResult()
			if result == ResultNInARow {
				return 1 * multiplier, x, y
			} else if result == ResultStalemate {
				return 0, x, y
			}

			r, _, _ := computeMove(gs, !isMax)

			if (isMax && r > threshold) || (!isMax && r < threshold) {
				threshold = r
				optimalX = x
				optimalY = y
			}
		}
	}

	return threshold, optimalX, optimalY
}

// TicTacToeStateHandler accepts a TicTacToeStateRequest representing the
// current state of the game and responds with a TicTacToeStateResponse
// describing the new state of the game.
func TicTacToeStateHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeHTTPError(w, http.StatusBadRequest, "could read request", err)
		return
	}

	req := &TicTacToeStateRequest{
		Board: makeBoard(3),
	}
	err = json.Unmarshal(b, req)
	if err != nil {
		writeHTTPError(w, http.StatusBadRequest, "could not interpret request", err)
		return
	}

	g := newGameFromRequest(req)
	result, _ := g.getGameResult()
	if result == ResultNone {
		_, x, y := computeMove(*g, true)
		err = g.occupyPosition(x, y)
		if err != nil {
			writeHTTPError(w, http.StatusInternalServerError, "failed to set board", err)
			return
		}
	}

	result, winningRow := g.getGameResult()
	resp := TicTacToeStateResponse{
		Board:      g.Board,
		Result:     result,
		Turn:       g.Turn,
		WinningRow: winningRow,
		NextPlayer: g.playersTurn(),
	}

	b, err = json.Marshal(resp)
	if err != nil {
		writeHTTPError(w, http.StatusInternalServerError, "failed to marshal response", err)
		return
	}

	_, err = w.Write(b)
	if err != nil {
		log.Fatal(err)
		return
	}
}
