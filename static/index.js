'use strict'
const canvas = document.getElementById('myCanvas')
const instructions = document.getElementById('instructions')
var size
var gameState = [[0, 0, 0], [0, 0, 0], [0, 0, 0]]
var winningRow = [[0, 0, 0], [0, 0, 0], [0, 0, 0]]
var gameHasResult = false
redraw()
canvas.addEventListener('mousedown', onMouseDown)
window.addEventListener('resize', redraw)

function redraw () {
  resize(canvas)
  drawBoard(canvas)
}

function onResetGame() {
  gameState = [[0, 0 ,0], [0, 0, 0], [0, 0, 0]]
  winningRow = [[0, 0, 0], [0, 0, 0], [0, 0, 0]]
  gameHasResult = false
  instructions.innerText = 'Next Player: X'
  drawBoard(canvas)
}

function resize (canvas) {
  const displayWidth = canvas.clientWidth
  const displayHeight = canvas.clientHeight
  if (canvas.width !== displayWidth ||
        canvas.height !== displayHeight) {
    canvas.width = displayWidth
    canvas.height = displayHeight
  }
  const minExtent = Math.min(canvas.clientWidth, canvas.clientHeight)
  size = minExtent / 4
}

function drawBoard (canvas) {
  const ctx = canvas.getContext('2d')
  ctx.clearRect(0, 0, canvas.width, canvas.height)
  ctx.beginPath()
  ctx.lineWidth = Math.min(10, size * 0.1)
  ctx.strokeStyle = 'rgb(0, 0, 0, 1)'
  ctx.moveTo(size, 0)
  ctx.lineTo(size, 3 * size)
  ctx.moveTo(2 * size, 0)
  ctx.lineTo(2 * size, 3 * size)
  ctx.moveTo(0, size)
  ctx.lineTo(3 * size, size)
  ctx.moveTo(0, 2 * size)
  ctx.lineTo(3 * size, 2 * size)
  ctx.stroke()

  drawState(canvas, gameState, winningRow)
}

function drawPlayer (canvas, x, y, player, emphasize) {
  const ctx = canvas.getContext('2d')
  const midSqX = (x + 0.5) * size
  const midSqY = y * size
  var color = 'rgb(255, 0, 0, 1)'
  ctx.clearRect((x + 0.1) * size, (y + 0.1) * size, 0.8 * size, 0.8 * size)
  if (player === 'X'.charCodeAt(0)) {
    if (!emphasize) {
      color = 'rgba(0, 196, 64, 1)'
    }
    ctx.beginPath()
    ctx.strokeStyle = color
    ctx.moveTo(midSqX - (0.3 * size), midSqY + (0.2 * size))
    ctx.lineTo(midSqX + (0.3 * size), midSqY + (0.8 * size))
    ctx.moveTo(midSqX - (0.3 * size), midSqY + (0.8 * size))
    ctx.lineTo(midSqX + (0.3 * size), midSqY + (0.2 * size))
    ctx.stroke()
  } else {
    if (!emphasize) {
      color = 'rgba(0, 0, 196, 1)'
    }
    ctx.beginPath()
    ctx.strokeStyle = color
    ctx.arc(midSqX, midSqY + (0.5 * size), 0.3 * size, 0, 2 * Math.PI)
    ctx.stroke()
  }
}

function getMouseClickCoords (canvas, event) {
  const rect = canvas.getBoundingClientRect()
  const x = event.clientX - rect.left
  const y = event.clientY - rect.top
  return [x, y]
}

function getTranslatedCoords (canvas, x, y) {
  const tx = Math.floor(x / size)
  const ty = Math.floor(y / size)
  if (tx < 0 || tx > 2 || ty < 0 || ty > 2) {
    return [false, null, null]
  }

  return [true, tx, ty]
}

function drawState (canvas, gameState, gameStateDone) {
  for (let j = 0; j < gameState.length; j++) {
    for (let i = 0; i < gameState[j].length; i++) {
      if (gameState[j][i] === 0) { continue }
      const emphasize = gameStateDone[j][i] === gameState[j][i]
      drawPlayer(canvas, i, j, gameState[j][i], emphasize)
    }
  }
}

function renderInstructions (response) {
  const next = response.nextPlayer
  const result = response.result
  if (result === 0) {
    instructions.innerText = 'Next Player: ' + String.fromCharCode(next)
  } else if (result === 1) {
    instructions.innerText = 'You lose!'
    if (gameHasResult) {
      instructions.innerText = "Don't be silly you lost!"
    }
  } else if (result === 2) {
    instructions.innerText = 'A draw'
  }
}

function player(gameState) {
  let count = 0
  for (let j = 0; j < gameState.length; j++) {
    for (let i = 0; i < gameState[j].length; i++) {
      if (gameState[j][i] !== 0){
        count++
      }
    }
  }
  if (count % 2 === 0) {
    return 'X'
  } else {
    return '0'
  }
}

function onMouseDown (e) {
  const [x, y] = getMouseClickCoords(canvas, e)
  const [isValid, tx, ty] = getTranslatedCoords(canvas, x, y)
  if (!isValid) {
    return
  }
  if (typeof gameState !== 'undefined' && gameState != null && gameState.length != null &&
        gameState[ty][tx] === 0) {
    gameState[ty][tx] = player(gameState).charCodeAt(0)
    drawState(canvas, gameState, winningRow)
    canvas.toBlob(function(blob)
    {
      saveAs(blob, "x.png")
    })
  }
}

function gameStateResponse (response) {
  gameState = response.data.board
  renderInstructions(response.data)
  if (response.data.winningRow !== undefined) {
    winningRow = response.data.winningRow
    gameHasResult = true
  }
  drawState(canvas, gameState, winningRow)
}
