// Copyright 2019 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Binary textinputdemo shows the functionality of a text input field.
package main

import (
	"context"
	"strconv"
	"time"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/align"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/keyboard"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/widgets/button"
	"github.com/mum4k/termdash/widgets/segmentdisplay"
	"github.com/mum4k/termdash/widgets/textinput"
	"./visualizer"
)

func rotate(inputs []rune, step int) []rune {
	return append(inputs[step:], inputs[:step]...)
}

// textState creates a rotated state for the text we are displaying.
func textState(text string, capacity, step int) []rune {
	if capacity == 0 {
		return nil
	}

	var state []rune
	for i := 0; i < capacity; i++ {
		state = append(state, ' ')
	}
	state = append(state, []rune(text)...)
	step = step % len(state)
	return rotate(state, step)
}

// rollText rolls a text across the segment display.
// Exists when the context expires.
func rollText(ctx context.Context, sd *segmentdisplay.SegmentDisplay) {
	colors := []cell.Color{
		cell.ColorNumber(33),
		cell.ColorRed,
		cell.ColorYellow,
		cell.ColorNumber(33),
		cell.ColorGreen,
		cell.ColorRed,
		cell.ColorGreen,
		cell.ColorRed,
	}

	text := "SORT VISUALIZER"
	step := 0
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			state := textState(text, sd.Capacity(), step)
			var chunks []*segmentdisplay.TextChunk
			for i := 0; i < sd.Capacity(); i++ {
				if i >= len(state) {
					break
				}

				color := colors[i%len(colors)]
				chunks = append(chunks, segmentdisplay.NewChunk(
					string(state[i]),
					segmentdisplay.WriteCellOpts(cell.FgColor(color)),
				))
			}
			if len(chunks) == 0 {
				continue
			}
			if err := sd.Write(chunks); err != nil {
				panic(err)
			}
			step++
		}
	}
}

func main() {
	t, err := tcell.New()
	if err != nil {
		panic(err)
	}
	defer t.Close()

	ctx, cancel := context.WithCancel(context.Background())
	rollingSD, err := segmentdisplay.New(
		segmentdisplay.MaximizeSegmentHeight(),
	)
	if err != nil {
		panic(err)
	}
	go rollText(ctx, rollingSD)

	input, err := createInput("Semilla:")
	if err != nil {
		panic(err)
	}
	input2, err := createInput("Periodo")
	if err != nil {
		panic(err)
	}
	input3, err := createInput("Tamaño:")
	if err != nil {
		panic(err)
	}

	EmpezarB, err := button.New("Empezar", func() error {
		var intLen,intPeriod,intSeed int
		primos := map[int]int{
			11: 11, 13: 13, 17: 17, 19: 19,
			23: 23, 29: 29, 31: 31, 37: 37,
			41: 41, 43: 43, 47: 47, 53: 53,
			59: 59, 61: 61, 67: 67, 71: 71,
			73: 73, 79: 79, 83: 83, 89: 89,
			97: 97,101: 101,
		}
		//Valida tamanno mayor a 0
		len := input3.ReadAndClear()
		intLen, _ = strconv.Atoi(len)
		if intLen < 0 {
			return nil
		}
		
		//Valida periodo mayor 2048
		period := input2.ReadAndClear()
		intPeriod, _ = strconv.Atoi(period)
		if intPeriod < 2048 {
			return nil
		}
		
		//Valida semilla que sea primo
		seed := input.ReadAndClear()
		intSeed, _ = strconv.Atoi(seed)
		if primos[intSeed] != intSeed{
			return nil
		}
		
		visualizer.Start(intLen, intSeed, intPeriod)
		//cancel()
		return nil
	},
		button.GlobalKey(keyboard.KeyEnter),
		button.FillColor(cell.ColorNumber(220)),
	)
	if err != nil {
		panic(err)
	}
	quitB, err := button.New("Quitar", func() error {
		cancel()
		return nil
	},
		button.WidthFor("Empezar"),
		button.FillColor(cell.ColorNumber(196)),
	)
	if err != nil {
		panic(err)
	}

	builder := grid.New()
	builder.Add(
		grid.RowHeightPerc(40,
			grid.Widget(
				rollingSD,
			),
		),
	)
	builder.Add(
		grid.RowHeightPerc(20,
			grid.ColWidthPerc(20),
			grid.ColWidthPerc(20,
				grid.Widget(
					input,
					container.AlignHorizontal(align.HorizontalCenter),
					container.AlignVertical(align.VerticalBottom),
					container.MarginBottom(1),
				),
			),
			grid.ColWidthPerc(20,
				grid.Widget(
					input2,
					container.AlignHorizontal(align.HorizontalCenter),
					container.AlignVertical(align.VerticalBottom),
					container.MarginBottom(1),
				),
			),
			grid.ColWidthPerc(20,
				grid.Widget(
					input3,
					container.AlignHorizontal(align.HorizontalCenter),
					container.AlignVertical(align.VerticalBottom),
					container.MarginBottom(1),
				),
			),
			grid.ColWidthPerc(20),
		),
	)

	builder.Add(
		grid.RowHeightPerc(40,
			grid.ColWidthPerc(20),
			grid.ColWidthPerc(20,
				grid.Widget(
					EmpezarB,
					container.AlignVertical(align.VerticalTop),
					container.AlignHorizontal(align.HorizontalRight),
				),
			),
			grid.ColWidthPerc(20,
				grid.Widget(
					quitB,
					container.AlignVertical(align.VerticalTop),
					container.AlignHorizontal(align.HorizontalLeft),
				),
			),
			grid.ColWidthPerc(20),
		),
	)

	gridOpts, err := builder.Build()
	if err != nil {
		panic(err)
	}
	c, err := container.New(t, gridOpts...)
	if err != nil {
		panic(err)
	}

	if err := termdash.Run(ctx, t, c, termdash.RedrawInterval(500*time.Millisecond)); err != nil {
		panic(err)
	}
}
func createInput(message string) (*textinput.TextInput, error) {
	input, err := textinput.New(
		textinput.Label(message, cell.FgColor(cell.ColorNumber(33))),
		textinput.MaxWidthCells(20),
		textinput.Border(linestyle.Light),
		textinput.PlaceHolder("Digite un numero"),
	)
	return input, err
}