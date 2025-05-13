package main

import (
	"image"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func main() {

	app := app.New()

	img := canvas.NewImageFromImage(nil)
	img.FillMode = canvas.ImageFillContain
	img.ScaleMode = canvas.ImageScaleFastest

	origImg := canvas.NewImageFromImage(nil)

	imgContainer := container.NewStack(img)

	DragAndDropwindow := app.NewWindow("Photoshop")
	DragAndDropwindow.SetOnDropped(func(pos fyne.Position, uris []fyne.URI) {

		if len(uris) > 0 {
			file, err := os.Open(uris[0].Path())
			if err != nil {
				return
			}
			defer file.Close()

			imgSrc, _, err := image.Decode(file)
			if err != nil {
				return
			}
			img.Image = imgSrc
			origImg.Image = imgSrc
			img.Refresh()
		}
	})

	origImgButton := NewOriginalButton(img, origImg)
	grayScaleButton := NewGrayscaleButton(img)
	negativeButton := NewNegativeButton(img, DragAndDropwindow)
	adjustBrightnessButton := NewAdjustBrightnessButton(img, DragAndDropwindow)
	binarizedButton := NewBinarizationButton(img, DragAndDropwindow)
	increaseContrastButton := NewIncreaseContrastButton(img, DragAndDropwindow)
	decreaseContrastButton := NewDecreaseContrastButton(img, DragAndDropwindow)
	createHistogramButton := NewCreateHistogramButton(img, DragAndDropwindow)
	gammaButton := NewGammaButton(img, DragAndDropwindow)
	quantizationButton := NewQuantizationButton(img, DragAndDropwindow)
	solarizationButton := NewSolarizationButton(img, DragAndDropwindow)
	lowFreqFilterButton := NewLowFreqFilterButton(img, DragAndDropwindow)
	highFreqFilterButton := NewHighFreqFilterButton(img, DragAndDropwindow)
	medianFilterButton := NewMedianFilterButton(img, DragAndDropwindow)
	gassBlurButton := NewGaussBlurButton(img, DragAndDropwindow)
	edgeEmpowerButton := NewEdgeEmpowerButton(img, DragAndDropwindow)
	shiftEdgeButton := NewShiftEdgeButton(img, DragAndDropwindow)
	embossingButton := NewEmbossingButton(img, DragAndDropwindow)
	kirschButton := NewKirschButton(img, DragAndDropwindow)
	pravitButton := NewPravitButton(img, DragAndDropwindow)
	sobelButton := NewSobelButton(img, DragAndDropwindow)
	robertsButton := NewRobertsButton(img, DragAndDropwindow)

	boxWithButtons := container.NewVBox(
		origImgButton,
		grayScaleButton,
		negativeButton,
		adjustBrightnessButton,
		binarizedButton,
		increaseContrastButton,
		decreaseContrastButton,
		createHistogramButton,
		gammaButton,
		quantizationButton,
		solarizationButton,
		lowFreqFilterButton,
		highFreqFilterButton,
		medianFilterButton,
		gassBlurButton,
		edgeEmpowerButton,
		shiftEdgeButton,
		embossingButton,
		kirschButton,
		pravitButton,
		sobelButton,
		robertsButton,
	)

	scrollButtons := container.NewVScroll(boxWithButtons)

	content := container.NewHSplit(
		scrollButtons,
		imgContainer,
	)

	content.SetOffset(0.2)

	DragAndDropwindow.SetContent(content)
	DragAndDropwindow.Resize(fyne.NewSize(900, 600))
	DragAndDropwindow.ShowAndRun()
}
