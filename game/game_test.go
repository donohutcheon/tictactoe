package game

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGame_SetBoard(t *testing.T) {
	type args struct {
		x      int
		y      int
	}
	tests := []struct {
		name             string
		initialGameState TicTacToeState
		args             args
		expErr           error
		expGameState     TicTacToeState
	}{
		{
			name:    "Opening Move",
			initialGameState:  TicTacToeState{
				Board: makeBoard(3),
				Turn:  1,
			},
			args:    args{
				x:      1,
				y:      1,
			},
			expGameState: TicTacToeState{
				Turn:  2,
				Board: [][]SquareState{
					{SquareStateEmpty, SquareStateEmpty, SquareStateEmpty},
					{SquareStateEmpty, SquareStateCross, SquareStateEmpty},
					{SquareStateEmpty, SquareStateEmpty, SquareStateEmpty},
				},
			},
		},
		{
			name:    "Mid-game Move",
			initialGameState:  TicTacToeState{
				Board: [][]SquareState{
					{SquareStateEmpty, SquareStateCross, SquareStateEmpty},
					{SquareStateNaught, SquareStateCross, SquareStateEmpty},
					{SquareStateEmpty, SquareStateNaught, SquareStateEmpty}},
					Turn: 5,
			},
			args:    args{
				x:      0,
				y:      0,
			},
			expGameState: TicTacToeState{
				Turn: 6,
				Board: [][]SquareState{
					{SquareStateCross, SquareStateCross, SquareStateEmpty},
					{SquareStateNaught, SquareStateCross, SquareStateEmpty},
					{SquareStateEmpty, SquareStateNaught, SquareStateEmpty},
				},
			},
		},
		{
			name:    "Out of bounds",
			args:    args{
				x:      3,
				y:      0,
			},
			expGameState: TicTacToeState{},
			expErr:       errors.New("invalid coordinate"),
		},
		{
			name:    "Already occupied",
			initialGameState: TicTacToeState{
				Board: [][]SquareState{{SquareStateNaught}},
				Turn: 2,
			},
			args:    args{
				x:      0,
				y:      0,
			},
			expGameState: TicTacToeState{
				Board: [][]SquareState{{SquareStateNaught}},
				Turn: 2,
			},
			expErr: errors.New("already occupied"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &tt.initialGameState
			err := g.occupyPosition(tt.args.x, tt.args.y)
			assert.Equal(t, tt.expErr, err)
			assert.Equal(t, tt.expGameState, *g)
		})
	}
}

func TestGame_CheckGameOver(t *testing.T) {
	type fields struct {
		Turn  int
		Board [][]SquareState
	}
	tests := []struct {
		name          string
		fields        fields
		want          Result
		expWinningRow [][]SquareState
	}{
		{
			name:   "Diagonal",
			fields: fields{
				Board: [][]SquareState{
					{SquareStateCross,SquareStateNaught,SquareStateEmpty},
					{SquareStateNaught,SquareStateCross,SquareStateEmpty},
					{SquareStateEmpty,SquareStateEmpty,SquareStateCross},
				},
			},
			want: ResultNInARow,
			expWinningRow: [][]SquareState{
				{SquareStateCross,SquareStateEmpty,SquareStateEmpty},
				{SquareStateEmpty,SquareStateCross,SquareStateEmpty},
				{SquareStateEmpty,SquareStateEmpty,SquareStateCross},
			},
		},
		{
			name:   "Anti-diagonal",
			fields: fields{
				Board: [][]SquareState{
					{SquareStateCross,SquareStateCross,SquareStateNaught},
					{SquareStateCross,SquareStateNaught,SquareStateCross},
					{SquareStateNaught,SquareStateNaught,SquareStateCross},
				},
			},
			want: ResultNInARow,
			expWinningRow: [][]SquareState{
				{SquareStateEmpty,SquareStateEmpty,SquareStateNaught},
				{SquareStateEmpty,SquareStateNaught,SquareStateEmpty},
				{SquareStateNaught,SquareStateEmpty,SquareStateEmpty},
			},
		},
		{
			name:   "Row",
			fields: fields{
				Board: [][]SquareState{
					{SquareStateNaught,SquareStateEmpty,SquareStateEmpty},
					{SquareStateCross,SquareStateCross,SquareStateCross},
					{SquareStateEmpty,SquareStateNaught,SquareStateEmpty},
				},
			},
			want: ResultNInARow,
			expWinningRow: [][]SquareState{
				{SquareStateEmpty,SquareStateEmpty,SquareStateEmpty},
				{SquareStateCross,SquareStateCross,SquareStateCross},
				{SquareStateEmpty,SquareStateEmpty,SquareStateEmpty},
			},
		},
		{
			name:   "Column",
			fields: fields{
				Board: [][]SquareState{
					{SquareStateCross,SquareStateNaught,SquareStateEmpty},
					{SquareStateCross,SquareStateNaught,SquareStateCross},
					{SquareStateEmpty,SquareStateNaught,SquareStateEmpty},
				},
			},
			want: ResultNInARow,
			expWinningRow: [][]SquareState{
				{SquareStateEmpty,SquareStateNaught, SquareStateEmpty},
				{SquareStateEmpty,SquareStateNaught, SquareStateEmpty},
				{SquareStateEmpty,SquareStateNaught, SquareStateEmpty},
			},
		},
		{
			name:   "Stalemate",
			fields: fields{
				Turn:  10,
				Board: makeBoard(3),
			},
			want:   ResultStalemate,
		},
		{
			name:   "Result None",
			fields: fields{
				Board: [][]SquareState{
					{SquareStateCross,SquareStateEmpty,SquareStateEmpty},
					{SquareStateEmpty,SquareStateNaught,SquareStateEmpty},
					{SquareStateEmpty,SquareStateEmpty,SquareStateEmpty},
				},
			},
			want:   ResultNone,
		},
		{
			name:   "Result None 2",
			fields: fields{
				Board: [][]SquareState{
					{SquareStateEmpty,SquareStateNaught,SquareStateEmpty},
					{SquareStateCross,SquareStateEmpty,SquareStateEmpty},
					{SquareStateEmpty,SquareStateEmpty,SquareStateEmpty},
				},
			},
			want:   ResultNone,
		},
		{
			name:   "Zero",
			want:   ResultNone,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &TicTacToeState{
				Turn:  tt.fields.Turn,
				Board: tt.fields.Board,
			}
			got, gotWinningRow := g.getGameResult()
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.expWinningRow, gotWinningRow)
		})
	}
}

func TestComputeMove(t *testing.T) {
	type args struct {
		gameState TicTacToeState
		player    rune
		result    int
	}
	tests := []struct {
		name  string
		args  args
		want  int
		wantX int
		wantY int
	}{
		{
			name:  "Naughts Win",
			args:  args{
				gameState: TicTacToeState{
					Turn:  6,
					Board: [][]SquareState{
						{SquareStateCross, SquareStateEmpty, SquareStateEmpty},
						{SquareStateCross, SquareStateCross, SquareStateNaught},
						{SquareStateEmpty, SquareStateEmpty, SquareStateNaught},
					},
				},
				player:    '0',
				result:    0,
			},
			want:  0,
			wantX: 2,
			wantY: 0,
		},
		{
			name:  "Naughts Defend",
			args:  args{
				gameState: TicTacToeState{
					Turn:  4,
					Board: [][]SquareState{
						{SquareStateNaught, SquareStateEmpty, SquareStateEmpty},
						{SquareStateCross, SquareStateCross, SquareStateEmpty},
						{SquareStateEmpty, SquareStateEmpty, SquareStateEmpty},
					},
				},
				player:    '0',
				result:    0,
			},
			want:  0,
			wantX: 2,
			wantY: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotX, gotY := computeMove(tt.args.gameState, true)
			t.Log(got, gotX, gotY)
			assert.Equal(t, tt.wantX, gotX)
			assert.Equal(t, tt.wantY, gotY)
		})
	}
}