module farni.com/rendermap

go 1.13

replace farni.com/ui => ../rendermap/ui

require (
	farni.com/ui v0.0.0-00010101000000-000000000000
	fyne.io/fyne/v2 v2.0.3
	github.com/fogleman/gg v1.3.0
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/hajimehoshi/ebiten/v2 v2.1.0
	github.com/paulmach/go.geojson v1.4.0
	golang.org/x/tools v0.1.0
)
