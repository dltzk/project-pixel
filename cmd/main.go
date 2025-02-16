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
			img.Refresh()
		}
	})

	grayScaleButton := GrayscaleButton(img)
	negativeButton := NegativeButton(img, DragAndDropwindow)
	adjustBrightnessButton := AdjustBrightnessButton(img, DragAndDropwindow)
	binarizedButton := BinarizationButton(img, DragAndDropwindow)

	leftButtons := container.NewVBox(
		grayScaleButton, 
		negativeButton, 
		adjustBrightnessButton,
		binarizedButton,
	)

	content := container.NewHSplit(
		leftButtons, 
		imgContainer,
	)

	content.SetOffset(0.2)

	DragAndDropwindow.SetContent(content)
	DragAndDropwindow.Resize(fyne.NewSize(900, 600))
	DragAndDropwindow.ShowAndRun()
}
