// MIT License
//
// Copyright (c) 2021 Yuchen Niu and EASE lab
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package profile

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"

	log "github.com/sirupsen/logrus"
)

// PlotCVS plots every attribute as VM number increases
func PlotCVS(xStep int, filePath, inFile, xLable string) {
	var (
		records = readResultCSV(filePath, inFile)
		rows    = len(records)
		cols    = len(records[0])
	)

	for col := 0; col < cols; col++ {
		// create a new plot for a metric
		p, err := plot.New()
		if err != nil {
			log.Fatalf("Failed creating plot: %v", err)
		}

		p.X.Label.Text = xLable
		p.Y.Label.Text = records[0][col]
		p.Y.Min = 0

		// setup data
		pts := make(plotter.XYs, rows-1)
		vmNum := xStep
		for row := 1; row < rows; row++ {
			pts[row-1].X = float64(vmNum)
			valStr := records[row][col]
			if valStr != "" {
				value, err := strconv.ParseFloat(valStr, 64)
				if err != nil {
					log.Fatalf("Failed parsing string to float: %v", err)
				}
				pts[row-1].Y = value
			}
			vmNum += xStep
		}

		err = plotutil.AddLinePoints(p, pts)
		if err != nil {
			log.Fatalf("Failed plotting data: %v", err)
		}

		p.Y.Label.Text = strings.ReplaceAll(p.Y.Label.Text, "/", "-")
		fileName := filepath.Join(filePath, p.Y.Label.Text+".png")
		if err := p.Save(4*vg.Inch, 4*vg.Inch, fileName); err != nil {
			log.Fatalf("Failed saving plot: %v", err)
		}
	}

	log.Info("Plot counters finished.")
}

// readResultCSV retrieves data from csv file
func readResultCSV(filePath, inFile string) [][]string {
	fileName := filepath.Join(filePath, inFile)
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Failed opening file: %v", err)
	}
	defer f.Close()

	r := csv.NewReader(f)

	records, err := r.ReadAll()
	if err != nil {
		log.Fatalf("Failed reading file %s: %v", filePath, err)
	}

	return records
}
