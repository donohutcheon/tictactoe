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

type SquareState rune
const (
	SquareStateEmpty  SquareState = 0
	SquareStateCross  SquareState = 'X'
	SquareStateNaught SquareState = '0'
)

type Result int
const (
	ResultNone      Result = iota
	ResultNInARow   Result = iota
	ResultStalemate Result = iota
)

type TicTacToeState struct {
	Board [][]SquareState `json:"board"`
	Turn  int             `json:"-"`
}

type TicTacToeStateResponse struct {
	Board      [][]SquareState   `json:"board"`
	Result     Result            `json:"result,omitempty"`
	WinningRow [][]SquareState   `json:"winningRow,omitempty"`
	Turn       int               `json:"turn"`
	NextPlayer rune              `json:"nextPlayer"`
}

// TicTacToeStateHandler accepts a TicTacToeState representing the
// current state of the game and responds with a TicTacToeStateResponse
// describing the new state of the game.
func TicTacToeStateHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Parapgraph #1
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeHTTPError(w, http.StatusBadRequest, "could read request", err)
		return
	}

	// Parapgraph #2
	req := &TicTacToeState{
		Board: makeBoard(3),
	}
	err = json.Unmarshal(b, req)
	if err != nil {
		writeHTTPError(w, http.StatusBadRequest, "could not interpret request", err)
		return
	}

	req.initialize()

	// Parapgraph #3
	result, _ := req.getGameResult()
	if result == ResultNone {
		_, x, y := computeMove(*req, true)
		err = req.occupyPosition(x, y)
		if err != nil {
			writeHTTPError(w, http.StatusInternalServerError, "failed to set board", err)
			return
		}
	}

	// Parapgraph #4
	result, winningRow := req.getGameResult()
	resp := TicTacToeStateResponse{
		Board:      req.Board,
		Result:     result,
		Turn:       req.Turn,
		WinningRow: winningRow,
		NextPlayer: req.playersTurn(),
	}

	// Parapgraph #5
	b, err = json.Marshal(resp)
	if err != nil {
		writeHTTPError(w, http.StatusInternalServerError, "failed to marshal response", err)
		return
	}

	// Parapgraph #6
	_, err = w.Write(b)
	if err != nil {
		log.Fatal(err)
		return
	}
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

func (t *TicTacToeState) initialize() {
	turn := 1
	for _, y := range t.Board {
		for _, x := range y {
			if x != SquareStateEmpty {
				turn++
			}
		}
	}

	t.Turn = turn
}

func (t *TicTacToeState) playersTurn() rune {
	if t.Turn % 2 == 1 {
		return 'X'
	}

	return '0'
}

func (t *TicTacToeState) isOccupied(x, y int) bool {
	return t.Board[y][x] != SquareStateEmpty
}

func (t *TicTacToeState) occupyPosition(x, y int) error {
	if x < 0 || x > len(t.Board) || y < 0 || y > len(t.Board) {
		return errors.New("invalid coordinate")
	}
	if t.Board[y][x] != SquareStateEmpty {
		return errors.New("already occupied")
	}

	player := t.playersTurn()
	t.Turn++
	if player == 'X' {
		t.Board[y][x] = SquareStateCross
		return nil
	}
	t.Board[y][x] = SquareStateNaught

	return nil
}

// getGameResult calculates the current state of the game returning the result
//  and the row that concluded the game if there is a complete row, nil otherwise.
func (t *TicTacToeState) getGameResult() (Result, [][]SquareState) {
	n := len(t.Board)
	var rowOfN [][]SquareState

	// Check diagonal
	rowOfN = makeBoard(n)
	var i int
	for i = 0; i < n - 1 && t.Board[i][i] == t.Board[i+1][i+1] && t.Board[i][i] != SquareStateEmpty; i++ {
		rowOfN[i][i] = t.Board[i][i]
		rowOfN[i+1][i+1] = t.Board[i+1][i+1]
	}
	if i == n - 1 {
		return ResultNInARow, rowOfN
	}

	// Check anti-diagonal
	rowOfN = makeBoard(n)
	for i = 0; i < n - 1 && t.Board[n-1-i][i] == t.Board[n-i-2][i+1] && t.Board[n-i-2][i+1] != SquareStateEmpty; i++ {
		rowOfN[n-1-i][i] = t.Board[n-1-i][i]
		rowOfN[n-i-2][i+1] = t.Board[n-i-2][i+1]
	}
	if i == n - 1 {
		return ResultNInARow, rowOfN
	}

	var j int
	for j = 0; j < n; j++ {
		// Check columns
		rowOfN = makeBoard(n)
		for i = 0; i < n-1 && t.Board[i][j] == t.Board[i+1][j] && t.Board[i][j] != SquareStateEmpty; i++ {
			rowOfN[i][j] = t.Board[i][j]
			rowOfN[i+1][j] = t.Board[i+1][j]
		}
		if i == n-1 {
			return ResultNInARow, rowOfN
		}

		// Check rows
		rowOfN = makeBoard(n)
		for i = 0; i < n-1 && t.Board[j][i] == t.Board[j][i+1] && t.Board[j][i] != SquareStateEmpty; i++ {
			rowOfN[j][i] = t.Board[j][i]
			rowOfN[j][i+1] = t.Board[j][i+1]
		}
		if i == n-1 {
			return ResultNInARow, rowOfN
		}
	}

	// Check for stalemate
	if t.Turn > n * n {
		return ResultStalemate, nil
	}

	return ResultNone, nil
}

func computeMove(gameState TicTacToeState, isMax bool) (int, int, int) {
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

			gs := TicTacToeState{
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
