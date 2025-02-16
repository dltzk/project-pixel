package main

import (
	"image"
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func GrayscaleButton(img *canvas.Image) fyne.CanvasObject {

	button := widget.NewButton("GrayScale", func() {

		if img.Image == nil {
			return
		}

		bounds := img.Image.Bounds()
		grayImg := image.NewRGBA(bounds)

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				currentPixelColor := img.Image.At(x, y)
				r, g, b, a := currentPixelColor.RGBA()
				r = r >> 8
				g = g >> 8
				b = b >> 8
				a = a >> 8

				grayScale := uint8(0.3*float64(r) + 0.59*float64(g) + 0.11*float64(b))

				grayImg.Set(x, y, color.RGBA{grayScale, grayScale, grayScale, uint8(a)})
			}

		}

		img.Image = grayImg
		img.Refresh()
	})

	return button
}

func NegativeButton(img *canvas.Image, window fyne.Window) fyne.CanvasObject {
	button := widget.NewButton("Negative", func() {
		if img.Image == nil {
			return
		}

		bounds := img.Image.Bounds()
		negativeImg := image.NewRGBA(bounds)

		getCeiling := widget.NewEntry()
		getCeiling.SetPlaceHolder("Число")
		getCeiling.SetText("0")

		content := container.NewVBox( 
			getCeiling,
			widget.NewLabel(""),
		)

		customDialog := dialog.NewCustomWithoutButtons("Negative", content, window)

		confirmButton := widget.NewButton("OK", func() {

			number := getCeiling.Text
			negativeCeiling, err := strconv.Atoi(number)

			if err != nil || negativeCeiling < 0 || negativeCeiling > 255 {
				dialog.ShowInformation("Ошибка", "Введите корректное число", window)
				return
			}

			customDialog.Hide()

			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					currentPixelColor := img.Image.At(x, y)
					r, g, b, a := currentPixelColor.RGBA()
					r = r >> 8
					g = g >> 8
					b = b >> 8
					a = a >> 8

					var newR, newG, newB uint32 = r, g, b

					if r >= uint32(negativeCeiling) {
						newR = (255 - r)
					}
					if g >= uint32(negativeCeiling) {
						newG = (255 - g)
					}
					if b >= uint32(negativeCeiling) {
						newB = (255 - b)
					}
	
					negativeImg.Set(x, y, color.RGBA{uint8(newR), uint8(newG), uint8(newB), uint8(a)})
				}
	
			}

			img.Image = negativeImg
			img.Refresh()
		})

		dissmisButton := widget.NewButton("Cancel", func() {
			customDialog.Hide()
		})
		
		fixedSizeButton := container.NewGridWrap(
			fyne.NewSize(100, 35),
			confirmButton,
			dissmisButton,
		)

		centeredButton := container.NewCenter(fixedSizeButton)
		content.Add(centeredButton)

		customDialog.Resize(fyne.NewSize(300,100))

		customDialog.Show()

	})

	return button
}

func AdjustBrightnessButton(img *canvas.Image, window fyne.Window) fyne.CanvasObject {
	button := widget.NewButton("Brightness", func() {
		if img.Image == nil {
			return
		}

		bounds := img.Image.Bounds()
		changedImg := image.NewRGBA(bounds)

		howBrightSlider := widget.NewSlider(0, 255)
		howBrightSlider.Value = 0
		howBrightSlider.Step = 1

		valueLabel := widget.NewLabel("Value: " + strconv.Itoa(int(howBrightSlider.Value)))

		howBrightSlider.OnChangeEnded = func(value float64) {
			valueLabel.SetText("Value: " + strconv.Itoa(int(value)))
		}

		content := container.NewVBox(
			valueLabel,
			howBrightSlider,
			widget.NewLabel(""),
		)

		customDialog := dialog.NewCustomWithoutButtons("Brightness", content, window)

		confirmButton := widget.NewButton("OK", func() {

			number := howBrightSlider.Value

			customDialog.Hide()

			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					currentPixelColor := img.Image.At(x, y)
					r, g, b, a := currentPixelColor.RGBA()
					r = r >> 8
					g = g >> 8
					b = b >> 8
					a = a >> 8

					newR := r + uint32(number)
					newG := g + uint32(number)
					newB := b + uint32(number)

					if newR > 255 {
						newR = 255
					}
					if newG > 255 {
						newG = 255
					}
					if newB > 255 {
						newB = 255
					}
	
					changedImg.Set(x, y, color.RGBA{uint8(newR), uint8(newG), uint8(newB), uint8(a)})
				}
	
			}

			img.Image = changedImg
			img.Refresh()
		})

		dissmisButton := widget.NewButton("Cancel", func() {
			customDialog.Hide()
		})
		
		fixedSizeButton := container.NewGridWrap(
			fyne.NewSize(100, 35),
			confirmButton,
			dissmisButton,
		)

		centeredButton := container.NewCenter(fixedSizeButton)
		content.Add(centeredButton)

		customDialog.Resize(fyne.NewSize(300,100))

		customDialog.Show()

	})

	return button
}

func BinarizationButton(img *canvas.Image, window fyne.Window) fyne.CanvasObject {
	button := widget.NewButton("Binarization", func() {
		if img.Image == nil {
			return
		}

		bounds := img.Image.Bounds()
		binarizedImg := image.NewRGBA(bounds)

		getCeiling := widget.NewEntry()
		getCeiling.SetPlaceHolder("Число")
		getCeiling.SetText("0")

		content := container.NewVBox( 
			getCeiling,
			widget.NewLabel(""),
		)

		customDialog := dialog.NewCustomWithoutButtons("Binarization", content, window)

		confirmButton := widget.NewButton("OK", func() {

			number := getCeiling.Text
			negativeCeiling, err := strconv.Atoi(number)

			if err != nil || negativeCeiling < 0 || negativeCeiling > 255 {
				dialog.ShowInformation("Ошибка", "Введите корректное число", window)
				return
			}

			customDialog.Hide()

			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					currentPixelColor := img.Image.At(x, y)
					r, g, b, a := currentPixelColor.RGBA()
					r = r >> 8
					g = g >> 8
					b = b >> 8
					a = a >> 8

					grayScale := uint8(0.3*float64(r) + 0.59*float64(g) + 0.11*float64(b))

					if grayScale < uint8(negativeCeiling) {
						binarizedImg.Set(x, y, color.RGBA{uint8(0), uint8(0), uint8(0), uint8(a)})
					} else {
						binarizedImg.Set(x, y, color.RGBA{uint8(255), uint8(255), uint8(255), uint8(a)})
					}
	

				}
	
			}

			img.Image = binarizedImg
			img.Refresh()
		})

		dissmisButton := widget.NewButton("Cancel", func() {
			customDialog.Hide()
		})
		
		fixedSizeButton := container.NewGridWrap(
			fyne.NewSize(100, 35),
			confirmButton,
			dissmisButton,
		)

		centeredButton := container.NewCenter(fixedSizeButton)
		content.Add(centeredButton)

		customDialog.Resize(fyne.NewSize(300,100))

		customDialog.Show()

	})

	return button
}