package main

import (
	"fmt"
	"image"
	"io/ioutil"
	"time"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	geojson "github.com/paulmach/go.geojson"
)

type Image struct {
	*image.RGBA
}

type BoundingBox struct {
	//52.30785 minlat, 9.53888 minlon, 52.38335 maxlat, 9.72136 maxlon
	MinLat float64
	MinLon float64
	MaxLat float64
	MaxLon float64
}

type Point struct {
	X, Y float64
}

type PNG struct {
	Width, Height float64
}

var FileName string

type RenderImage struct {
	Counter int
	Status  string
	Image   image.Image
}

func (r RenderImage) Update() {
	newImage = true
	mapImage = r.Image
}

var RI RenderImage

func initMap() geojson.FeatureCollection {
	RI.Counter = 12
	RI.Update()

	file, err := ioutil.ReadFile(FileName)
	if err != nil {
		fmt.Println("File indutten", err)
		panic(err)
	}
	fc, err := geojson.UnmarshalFeatureCollection(file)

	if err != nil {
		fmt.Print("should unmarshal feature collection without issue ")
		panic(err)
	}

	if fc.Type != "FeatureCollection" {
		fmt.Print("should have type of FeatureCollection, got %v", fc.Type)
		panic(err)
	}

	//remove unneeded items
	x := geojson.NewFeatureCollection()

	for i, _ := range fc.Features {
		f := fc.Features[i]
		switch {

		case f.Properties["natural"] == "water":
			x.Features = append(x.Features, f)
		case f.Properties["type"] == "waterway":
			x.Features = append(x.Features, f)
		case f.Properties["water"] != nil:
			x.Features = append(x.Features, f)
		case f.Properties["waterway"] != nil:
			x.Features = append(x.Features, f)
		case f.Properties["landuse"] == "forest": // ||  f.Properties["tags"].data["landuse"] == "forest"  {
			x.Features = append(x.Features, f)

		}

	}
	return *x
}

func makeMap() {

	fc := initMap()
	bb := BoundingBox{}
	if fc.BoundingBox != nil {
		// needs to be checked. Til now I didn't find a Boundingbox in data
		fmt.Print(fc.BoundingBox)
		bb.MinLat = fc.BoundingBox[0]
		bb.MinLon = fc.BoundingBox[1]
		bb.MaxLat = fc.BoundingBox[2]
		bb.MaxLon = fc.BoundingBox[3]
	} else {
		fmt.Println("Boundingbox not found")
		bb = calcBoundingBox(&fc, bb)
	}
	png := PNG{1024, 768}
	dc := gg.NewContext(int(png.Width), int(png.Height))
	dc.SetHexColor("fff")
	dc.Clear()

	for i, _ := range fc.Features {
		f := fc.Features[i]

		g := f.Geometry

		var typ string
		/*
			if f.Properties["name"]=="Waldkaterbach"{
				fmt.Println("Waldkaterbach")
			}
		*/
		switch {

		case f.Properties["natural"] == "water":
			typ = "water"
		case f.Properties["type"] == "waterway":
			typ = "water"
		case f.Properties["water"] != nil:
			typ = "water"
		case f.Properties["waterway"] != nil:
			typ = "water"
		case f.Properties["landuse"] == "forest": // ||  f.Properties["tags"].data["landuse"] == "forest"  {
			typ = "forest"
		default:
			typ = "unknown"

		}
		switch {
		case g.IsPolygon():
			for k, _ := range g.Polygon {
				//fmt.Println(g.Polygon[0][0])
				renderGeo(g, g.Polygon[k], dc, &bb, png, typ)
			}

		case g.IsMultiPolygon():
			for k, _ := range g.MultiPolygon {
				for p, _ := range g.MultiPolygon[k] {
					renderGeo(g, g.MultiPolygon[k][p], dc, &bb, png, typ)
				}
			}
		case g.IsLineString():
			renderGeo(g, g.LineString, dc, &bb, png, typ)

		case g.IsMultiLineString():
			for k, _ := range g.MultiLineString {
				renderGeo(g, g.MultiLineString[k], dc, &bb, png, typ)
			}
		}

		fmt.Printf("feature %d from %d rendered \n", i+1, len(fc.Features))
		m := image.Image(dc.Image())
		RI.Counter = 1
		RI.Image = m
		RI.Update()

	}

	dc.SetRGBA(0, 0.5, 0, 0.2)
	dc.SetLineWidth(1)

	dc.StrokePreserve()
	filename := fmt.Sprint("data/outtmp.png")
	dc.SavePNG(filename)

	m := image.Image(dc.Image())
	RI.Counter++
	RI.Image = m
	RI.Update()
	RI.Status = "newImage"
}

func renderGeo(g *geojson.Geometry, geo [][]float64, dc *gg.Context, bb *BoundingBox, png PNG, typ string) {
	switch typ {
	case "water":

		dc.SetRGBA(0, 0, 0.7, 0.5)
	case "forest":
		dc.SetRGBA(0, 0.8, 0, 0.6)

	default:
		return
	}
	//
	p := Point{}

	for j, _ := range geo { //[pol] {

		max := (bb.MaxLon - bb.MinLon)
		lon := (bb.MaxLon - geo[j][1])
		p.X = lon / max * png.Height

		max = (bb.MaxLat - bb.MinLat)
		lat := (max - (bb.MaxLat - geo[j][0]))
		p.Y = lat / max * png.Width

		dc.LineTo(p.Y, p.X)
	}

	switch {
	case g.IsPolygon():

		dc.Fill()
	case g.IsMultiPolygon():
		dc.Fill()
	case g.IsLineString():
		dc.SetLineWidth(2)
		dc.Stroke()
	case g.IsMultiLineString():
		dc.SetFillRule(gg.FillRuleWinding)
		dc.SetLineWidth(2)
		dc.Stroke()
	}
	m := image.Image(dc.Image())
	RI.Image = m
	time.Sleep(10 * time.Millisecond)
}

func calcBoundingBox(fc *geojson.FeatureCollection, bb BoundingBox) BoundingBox {
	for i, _ := range fc.Features {
		f := fc.Features[i]
		g := f.Geometry
		if g.Type == "MultiPolygon" {
			fmt.Println(g.MultiPolygon[0][0][0])
			for k, _ := range g.MultiPolygon {
				for p, _ := range g.MultiPolygon[k] {
					setMinMaxLonLat(g.MultiPolygon[k][p], &bb)
				}
			}
		}

		if g.Type == "Polygon" {
			fmt.Println(g.Polygon[0][0])
			for p, _ := range g.Polygon {
				setMinMaxLonLat(g.Polygon[p], &bb)
			}
		}
		switch g.Type {
		case "LineString":
			setMinMaxLonLat(g.LineString, &bb)

		case "MultiLineString":
			for k, _ := range g.MultiLineString {
				setMinMaxLonLat(g.MultiLineString[k], &bb)
			}
		}

	}
	fmt.Println("BB")
	fmt.Println(bb)
	return bb
}

func setMinMaxLonLat(coordinates [][]float64, bb *BoundingBox) {
	for j, _ := range coordinates {
		if coordinates[j][1] > bb.MaxLon {
			bb.MaxLon = coordinates[j][1]
		}
		if coordinates[j][1] < bb.MinLon || bb.MinLon == 0 {
			bb.MinLon = coordinates[j][1]
		}
		if coordinates[j][0] > bb.MaxLat {
			bb.MaxLat = coordinates[j][0]
		}
		if coordinates[j][0] < bb.MinLat || bb.MinLat == 0 {
			bb.MinLat = coordinates[j][0]
		}
	}

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
var ebitMapImage *ebiten.Image
var i int
var err error

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	i++
	switch RI.Status {
	case "init":
		RI.Status = "waiting"
		go makeMap()
		RI.Status = "newImage"
	case "newImage":
		RI.Status = "waiting"

		if err != nil {
			fmt.Println(err)
		}
	}
	if RI.Image != nil {
		screen.DrawImage(ebiten.NewImageFromImage(RI.Image), nil)
	}
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1024, 768
}

func main() {
	//FileName = "data/velba.geojson"
	FileName = "../tmp/suentel.geojson"
	game := &Game{}
	// Sepcify the window size as you like. Here, a doulbed size is specified.
	ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowTitle("drawmap")
	ebiten.SetMaxTPS(1)
	// Call ebiten.RunGame to start your game loop.
	RI.Status = "init"
	if err := ebiten.RunGame(game); err != nil {
		fmt.Println(err)
	}

}
