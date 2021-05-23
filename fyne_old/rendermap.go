package main

import (
	"fmt"
	"image"
	"io/ioutil"

	"github.com/fogleman/gg"
	geojson "github.com/paulmach/go.geojson"

	//	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	//	"fyne.io/fyne/v2/canvas"
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

func initMap() geojson.FeatureCollection {

	//FileName = "export_velba.geojson"
	//FileName = "export_velba_total.geojson"
	//file, err := ioutil.ReadFile("export_springe.geojson")
	//FileName = "export_harz.geojson"
	//FileName = "export_süntel_total.geojson"
	FileName = "bbb.geojson"
	//FileName = "export_toledo.geojson"
	//FileName = "test.json"
	//FileName = "export_steinhude.geojson"
	//FileName = "export_deister_total.geojson"
	//file, err := ioutil.ReadFile("export_deister.geojson")
	//file, err := ioutil.ReadFile("export.geojson")
	//file, err := ioutil.ReadFile("export_benther.geojson")
	file, err := ioutil.ReadFile(FileName)
	if err != nil {
		fmt.Println("File indutten", err)
		panic(err)
	}
	fc, err := geojson.UnmarshalFeatureCollection(file)

	//fc, err := geojson.UnmarshalFeatureCollection([]byte(file))
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
		//tags:=mapTags(f.Properties["tags"])

		//for _,tag := range f.Properties["tags"]{

		//	g := f.Geometry
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

//func makeMap(fc *geojson.FeatureCollection, bb BoundingBox) {
func makeMap(w fyne.Window) {

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
	//fmt.Println(fc.Features[2])
	png := PNG{1000, 1000}
	dc := gg.NewContext(int(png.Width), int(png.Height))
	dc.SetHexColor("fff")
	dc.Clear()

	m := image.Image(dc.Image())
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
		/*if f.Properties["natural"] == "water" || f.Properties["type"] == "waterway" || f.Properties["water"] != nil || f.Properties["waterway"] != nil {
			typ = "water"
		} else if f.Properties["landuse"] == "forest" {
			typ = "forest"
		} else {
			typ = "unknown"
		}
		*/
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

		//fmt.Println(g)
		fmt.Printf("feature %d from %d rendered \n", i+1, len(fc.Features))
		//if g.Type == "Polygon" {
		//dc.Push()

	}

	m = image.Image(dc.Image())
	w.SetContent(canvas.NewImageFromImage(m))
	//time.Sleep(100 * time.Millisecond)
	//	dc.SetFillRule(gg.FillRuleEvenOdd)
	//	dc.SetFillRule(gg.FillRuleWinding)
	//dc.SetRGBA(0, 1, 0, 0.5)
	//dc.SetLineWidth(14)
	dc.SetRGBA(0, 0.5, 0, 0.2)
	dc.SetLineWidth(1)

	dc.StrokePreserve()
	filename := fmt.Sprint("out/out.", FileName, ".png")
	dc.SavePNG(filename)

	//myCanvas := w.Show()
	//myCanvas.Image = dc.Image
	//rect := image.Rect(0, 0, 255, 255)
	//myImage := image.NewRGBA(rect)

	//im := Image{myImage}
	//w.SetContent(canvas.NewImageFromImage(im))
	//w.SetContent(canvas.NewImageFromFile(filename))
	m = image.Image(dc.Image())
	w.SetContent(canvas.NewImageFromImage(m))
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
		//lon := max - (bb.MaxLon - g.Polygon[pol][j][1])
		lon := (bb.MaxLon - geo[j][1])
		p.X = lon / max * png.Width

		max = (bb.MaxLat - bb.MinLat)
		lat := (max - (bb.MaxLat - geo[j][0]))
		p.Y = lat / max * png.Height
		//fmt.Println(pos)
		//		fmt.Println(p)

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
		//dc.Stroke()
		//dc.SavePNG(filename)
		//dc.Pop()
	}
}

func calcBoundingBox(fc *geojson.FeatureCollection, bb BoundingBox) BoundingBox {
	//fmt.Println(fc.Features[2])
	for i, _ := range fc.Features {
		f := fc.Features[i]
		g := f.Geometry
		//	fmt.Println(g)
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

			bb.MinLat = bb.MinLat + 0
			bb.MinLat = bb.MinLat + 0
			bb.MinLon = bb.MinLon + 0

		}

	}
	/*
		bb.MinLat = 2
		bb.MinLon = 11
		bb.MaxLat = 3
		bb.MaxLon = 12
	*/
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

/*
func RenderMap(Width float64, Height float64, Mapname string) string {
	s := Mapname
	return fmt.Sprintf(s)
}
*/

func main() {

	a := app.New()
	w := a.NewWindow("map")
	w.Resize(fyne.NewSize(1000, 1000))
	//fc1 = new geojson.FeatureCollection
	/*
		go func() {
			makeMap(&fc, bb)
		}()

	*/
	go makeMap(w)
	w.ShowAndRun()
	/*fmt.Println(fc.Features[1])

	f := fc.Features[1]
	fmt.Println(f.Properties["landuse"])
	fmt.Println(f.Properties["name"])

	g := f.Geometry

	//fmt.Println(g)
	fmt.Println(g.Type)
	if g.Type == "Polygon" {
		fmt.Println(g.Polygon[0][0])
	}

	fmt.Println("BB")
	fmt.Println(bb)
	*/
}
