package game

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGame_String(t *testing.T) {
	type fields struct {
		Board [][]SquareState
	}
	tests := []struct {
		name   string
		fields fields
		expected string
	}{
		{
			name:   "Golden",
			fields: fields{
				Board: [][]SquareState{
					{SquareStateEmpty,SquareStateNaught,SquareStateNaught},
					{SquareStateCross,SquareStateCross,SquareStateNaught},
					{SquareStateNaught,SquareStateCross,SquareStateCross},
				},
			},
			expected: " |0|0\n-+-+-\nX|X|0\n-+-+-\n0|X|X\n",
		},
		{
			name:   "Zero",
			fields: fields{
				Board: nil,
			},
			expected: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				Board: tt.fields.Board,
			}
			str := g.String()
			assert.Equal(t, tt.expected, str)
		})
	}
}

func TestGame_SetBoard(t *testing.T) {
	type fields struct {
		Board [3][3]SquareState
	}
	type args struct {
		x      int
		y      int
	}
	tests := []struct {
		name    string
		initialGameState  Game
		args    args
		expErr  error
		expGameState Game
	}{
		{
			name:    "Opening Move",
			initialGameState:  Game{
				Board: MakeBoard(3),
				Turn: 1,
			},
			args:    args{
				x:      1,
				y:      1,
			},
			expGameState: Game{
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
			initialGameState:  Game{
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
			expGameState: Game{
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
			expGameState: Game{},
			expErr: errors.New("invalid coordinate"),
		},
		{
			name:    "Already occupied",
			initialGameState: Game{
				Board: [][]SquareState{{SquareStateNaught}},
				Turn: 2,
			},
			args:    args{
				x:      0,
				y:      0,
			},
			expGameState: Game{
				Board: [][]SquareState{{SquareStateNaught}},
				Turn: 2,
			},
			expErr: errors.New("already occupied"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &tt.initialGameState
			err := g.SetBoard(tt.args.x, tt.args.y)
			assert.Equal(t, tt.expErr, err)
			assert.Equal(t, tt.expGameState, *g)
		})
	}
}

func TestSquareState_Opponent(t *testing.T) {
	tests := []struct {
		name        string
		player      rune
		expOpponent rune
	}{
		{
			name:        "Cross",
			player:      '0',
			expOpponent: 'X',
		},
		{
			name:        "Naught",
			player:      'X',
			expOpponent: '0',
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opponent := opponent(tt.player)
			assert.Equal(t, tt.expOpponent, opponent)
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
				Turn: 10,
				Board: MakeBoard(3),
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
			g := &Game{
				Turn:  tt.fields.Turn,
				Board: tt.fields.Board,
			}
			got, gotWinningRow := g.GetGameResult()
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.expWinningRow, gotWinningRow)
		})
	}
}