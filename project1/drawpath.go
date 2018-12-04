package main

//sample code at: https://github.com/campoy/justforfunc/blob/master/34-gonum-plot/main.go

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
)

func drawPath(trip *Trip, fn string) {
	xys := makePlotter(*trip.path)
	err := plotData(xys, fn)
	if err != nil {
		log.Fatalf("Could not plot data: %v", err)
	}
}

func makePlotter(nn []Node) plotter.XYs {
	var xys plotter.XYs
	for _, v := range nn {
		xys = append(xys, struct{ X, Y float64 }{v.X, v.Y})
	}
	return xys
}

func plotData(xys plotter.XYs, fn string) error {
	filename := fn + ".png"
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create %s: %v", filename, err)
	}

	//create new plot
	p, err := plot.New()
	if err != nil {

	}

	//create scatter with points read
	s, err := plotter.NewScatter(xys)
	if err != nil {
		return fmt.Errorf("could not create scatter plot: %v", err)
	}

	s.GlyphStyle.Shape = draw.CrossGlyph{}
	s.Color = color.RGBA{R: 255, A: 255}
	p.Add(s)

	//create lines
	l, err := plotter.NewLine(xys)
	p.Add(l)

	//write to file
	wt, err := p.WriterTo(256, 256, "png")
	if err != nil {
		return fmt.Errorf("could not create writer: %v", err)
	}

	_, err = wt.WriteTo(f)
	if err != nil {
		return fmt.Errorf("could not write to %s: %v", filename, err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("could not close %s, %v", filename, err)
	}

	return nil
}
