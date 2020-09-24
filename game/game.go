package game

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

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
	Board [][]SquareState
}

type TicTacToeStateResponse struct {
	Board      [][]SquareState   `json:"board"`
	Result     Result            `json:"result,omitempty"`
	WinningRow [][]SquareState   `json:"winningRow,omitempty"`
	Turn       int               `json:"turn"`
	NextPlayer rune              `json:"nextPlayer"`
}

func MakeBoard(n int) [][]SquareState {
	board := make([][]SquareState, n)
	for i := 0; i < n; i++ {
		board[i] = make([]SquareState, n)
	}
	return board
}

func CopyBoard(src [][]SquareState) [][]SquareState {
	n := len(src)
	board := make([][]SquareState, n)
	for j := 0; j < n; j++ {
		board[j] = make([]SquareState, n)
		copy(board[j], src[j])
	}
	return board
}

func NewGameFromRequest(req *TicTacToeStateRequest) *Game {
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

func (g *Game) String() string{
	sb := new(strings.Builder)
	for y := range g.Board {
		for x, state := range g.Board[y] {
			switch state {
			case SquareStateEmpty:
				sb.WriteString(" ")
			case SquareStateCross:
				sb.WriteString("X")
			case SquareStateNaught:
				sb.WriteString("0")
			}
			if x < 2 {
				sb.WriteString("|")
			}
		}
		sb.WriteString("\n")
		if y < 2 {
			sb.WriteString("-+-+-\n")
		}
	}

	return sb.String()
}

func (g *Game) SetBoard(x, y int) error {
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

func (g *Game) GetGameResult() (Result, [][]SquareState) {
	n := len(g.Board)
	var rowOfN [][]SquareState
	// Check diagonal
	rowOfN = MakeBoard(n)
	var i int
	for i = 0; i < n - 1 && g.Board[i][i] == g.Board[i+1][i+1] && g.Board[i][i] != SquareStateEmpty; i++ {
		rowOfN[i][i] = g.Board[i][i]
		rowOfN[i+1][i+1] = g.Board[i+1][i+1]
	}
	if i == n - 1 {
		return ResultNInARow, rowOfN
	}
	// Check anti-diagonal
	rowOfN = MakeBoard(n)
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
		rowOfN = MakeBoard(n)
		for i = 0; i < n-1 && g.Board[i][j] == g.Board[i+1][j] && g.Board[i][j] != SquareStateEmpty; i++ {
			rowOfN[i][j] = g.Board[i][j]
			rowOfN[i+1][j] = g.Board[i+1][j]
		}
		if i == n-1 {
			return ResultNInARow, rowOfN
		}

		// Check rows
		rowOfN = MakeBoard(n)
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

func SetGameState(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		_, wErr := w.Write([]byte("Oh dear, something messed up"))
		if wErr != nil {
			log.Fatal(wErr)
		}
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	req := &TicTacToeStateRequest{
		Board: MakeBoard(3),
	}
	err = json.Unmarshal(b, req)
	if err != nil {
		_, wErr := w.Write([]byte("could not interpret request"))
		if wErr != nil {
			log.Fatal(wErr)
		}
		w.WriteHeader(http.StatusBadRequest)
	}

	g := NewGameFromRequest(req)
	_, x, y := ComputeMove(*g, true)
	g.SetBoard(x, y)

	result, winningRow := g.GetGameResult()
	resp := TicTacToeStateResponse{
		Board:      g.Board,
		Result:     result,
		Turn:       g.Turn,
		WinningRow: winningRow,
		NextPlayer: g.playersTurn(),
	}

	b, err = json.Marshal(resp)
	if err != nil {
		_, wErr := w.Write([]byte("Oh dear, something messed up"))
		if wErr != nil {
			log.Fatal(wErr)
		}
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(b)
}