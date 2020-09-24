package game

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestComputeMove(t *testing.T) {
	type args struct {
		gameState Game
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
				gameState: Game{
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
				gameState: Game{
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
			got, gotX, gotY := ComputeMove(tt.args.gameState, true)
			t.Log(got, gotX, gotY)
			assert.Equal(t, tt.wantX, gotX)
			assert.Equal(t, tt.wantY, gotY)
		})
	}
}
