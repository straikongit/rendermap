package ui

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	//	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type RenderImage struct {
	Counter int
	Image   image.Image
}

func (r RenderImage) Update() {
	fmt.Println("Counter = %v", r.Counter)
	newImage = true
	mapImage = r.Image
}

func main() {
	fmt.Println("hi von ui")
}

type Game struct{}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// Write your game's logical update.
	return nil
}

var newImage bool
var mapImage image.Image
var i int

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	i++
	if newImage {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, 0)

		ebitMapImage := ebiten.NewImageFromImage(mapImage)
		screen.DrawImage(ebitMapImage, op)
	}
	newImage = false
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1000, 1000
}

func Start() {
	//ebiten
	game := &Game{}
	// Sepcify the window size as you like. Here, a doulbed size is specified.
	ebiten.SetWindowSize(1000, 1000)
	ebiten.SetWindowTitle("drawmap")
	ebiten.SetMaxTPS(1)
	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		fmt.Println(err)
	}
}
