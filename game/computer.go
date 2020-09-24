package game

import "log"

func ComputeMove(gameState Game, isMax bool) (int, int, int) {
	optimalX := 0
	optimalY := 0
	multiplier := 1
	if !isMax {
		multiplier = -1
	}
	threshold := 32767 * -1 * multiplier

	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			if gameState.isOccupied(x, y) {
				continue
			}

			gs := Game{
				Board: CopyBoard(gameState.Board),
				Turn:  gameState.Turn,
			}

			player := gameState.playersTurn()
			err := gs.SetBoard(x, y)
			if err != nil {
				log.Fatal(err)
				continue
			}
			result, _ := gs.GetGameResult()
			if result == ResultNInARow {
				return 1 * multiplier, x, y
			} else if result == ResultStalemate {
				return 0, x, y
			}

			r, _, _ := ComputeMove(gs, !isMax)

			if (isMax && r > threshold) || (!isMax && r < threshold) {
				threshold = r
				optimalX = x
				optimalY = y
			}
		}
	}

	return threshold, optimalX, optimalY
}

