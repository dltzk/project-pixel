package main

import (
	"image"
	"image/color"
	"math"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func checkForLimit(value float64) float64 {
	if value > 255 {
		return 255
	} else if value < 0 {
		return 0
	}
	return value
}

func originalButton(img *canvas.Image, origImg *canvas.Image) fyne.CanvasObject {
	button := widget.NewButton("Original", func() {
		if img.Image == nil || origImg == nil {
			return
		}
		img.Image = origImg.Image
		img.Refresh()
	})

	return button
}

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

		customDialog.Resize(fyne.NewSize(300, 100))

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

		howBrightSlider.OnChanged = func(value float64) {
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

		customDialog.Resize(fyne.NewSize(300, 100))

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

		customDialog.Resize(fyne.NewSize(300, 100))

		customDialog.Show()

	})

	return button
}

func increaseContrastButton(img *canvas.Image, window fyne.Window) fyne.CanvasObject {
	button := widget.NewButton("Contrast+", func() {
		if img.Image == nil {
			return
		}

		bounds := img.Image.Bounds()
		changedImg := image.NewRGBA(bounds)

		getQ2 := widget.NewEntry()
		getQ1 := widget.NewEntry()
		getQ2.SetPlaceHolder("Q2")
		getQ1.SetPlaceHolder("Q1")

		content := container.NewVBox(
			getQ2,
			getQ1,
			widget.NewLabel(""),
		)

		customDialog := dialog.NewCustomWithoutButtons("Increase contrast", content, window)

		confirmButton := widget.NewButton("OK", func() {

			Q2 := getQ2.Text
			Q1 := getQ1.Text
			newQ2, err := strconv.Atoi(Q2)
			if err != nil {
				dialog.ShowInformation("Ошибка", "Введите корректное число", window)
				return
			}
			newQ1, err := strconv.Atoi(Q1)
			if err != nil {
				dialog.ShowInformation("Ошибка", "Введите корректное число", window)
				return
			}

			if newQ1 < 0 || newQ1 > 255 || newQ2 < 0 || newQ2 > 255 || newQ1 > newQ2 || newQ1 == newQ2 {
				dialog.ShowInformation("Ошибка", "Q1 не может быть больше Q2 или быть равен ему", window)
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

					coefficient := 255. / float64(newQ2-newQ1)

					newR := checkForLimit(float64(int(r)-int(newQ1)) * coefficient)
					newG := checkForLimit(float64(int(g)-int(newQ1)) * coefficient)
					newB := checkForLimit(float64(int(b)-int(newQ1)) * coefficient)

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

		customDialog.Resize(fyne.NewSize(300, 100))

		customDialog.Show()

	})

	return button
}

func decreaseContrastButton(img *canvas.Image, window fyne.Window) fyne.CanvasObject {
	button := widget.NewButton("Contrast-", func() {
		if img.Image == nil {
			return
		}

		bounds := img.Image.Bounds()
		changedImg := image.NewRGBA(bounds)

		getQ2 := widget.NewEntry()
		getQ1 := widget.NewEntry()
		getQ2.SetPlaceHolder("Q2")
		getQ1.SetPlaceHolder("Q1")

		content := container.NewVBox(
			getQ2,
			getQ1,
			widget.NewLabel(""),
		)

		customDialog := dialog.NewCustomWithoutButtons("Decrease contrast", content, window)

		confirmButton := widget.NewButton("OK", func() {

			Q2 := getQ2.Text
			Q1 := getQ1.Text
			newQ2, err := strconv.Atoi(Q2)
			if err != nil {
				dialog.ShowInformation("Ошибка", "Введите корректное число", window)
				return
			}
			newQ1, err := strconv.Atoi(Q1)
			if err != nil {
				dialog.ShowInformation("Ошибка", "Введите корректное число", window)
				return
			}

			if newQ1 < 0 || newQ1 > 255 || newQ2 < 0 || newQ2 > 255 || newQ1 > newQ2 || newQ1 == newQ2 {
				dialog.ShowInformation("Ошибка", "Q1 не может быть больше Q2 или быть равен ему", window)
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

					var newR, newG, newB uint32

					// rightNumber := uint32( (newQ2 - newQ1) / 255 )
					newR = uint32(newQ1) + (r*uint32(newQ2-newQ1))/255
					newG = uint32(newQ1) + (g*uint32(newQ2-newQ1))/255
					newB = uint32(newQ1) + (b*uint32(newQ2-newQ1))/255

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

		customDialog.Resize(fyne.NewSize(300, 100))

		customDialog.Show()

	})

	return button
}

func createHistogramButton(img *canvas.Image, window fyne.Window) fyne.CanvasObject {
	button := widget.NewButton("Histogram", func() {
		if img.Image == nil {
			return
		}

		// Создаём массив для хранения количества пикселей каждой яркости (0-255)
		histogram := make([]int, 256)

		bounds := img.Image.Bounds()
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				currentPixelColor := img.Image.At(x, y)
				r, g, b, _ := currentPixelColor.RGBA()
				r = r >> 8
				g = g >> 8
				b = b >> 8

				var color int
				if r != g && g != b {
					color = int(0.3*float64(r) + 0.59*float64(g) + 0.11*float64(b))
				} else {
					color = int(r)
				}
				histogram[color]++
			}
		}

		// Создаём изображение для гистограммы
		histogramImg := image.NewRGBA(image.Rect(0, 0, 512, 400))

		// Заливаем фон чёрным цветом
		for y := 0; y < 521; y++ {
			for x := 0; x < 521; x++ {
				histogramImg.Set(x, y, color.RGBA{0, 0, 0, 255}) // Чёрный фон
			}
		}

		// Находим максимальное значение в гистограмме для масштабирования
		maxCount := 0
		for _, count := range histogram {
			if count > maxCount {
				maxCount = count
			}
		}

		// Рисуем гистограмму
		barWidth := 1 // Ширина каждого столбца
		spacing := 1  // Отступ между столбцами
		for i, count := range histogram {
			if maxCount == 0 {
				continue
			}
			barHeight := (count * 380) / maxCount // Масштабируем высоту столбца (380 вместо 400 для отступа сверху)
			for x := i * (barWidth + spacing); x < i*(barWidth+spacing)+barWidth; x++ {
				for y := 400 - barHeight; y < 400; y++ {
					histogramImg.Set(x, y, color.RGBA{255, 255, 255, 255}) // Белый цвет для столбцов
				}
			}
		}

		// Преобразуем изображение гистограммы в canvas.Image
		canvasHistogram := canvas.NewImageFromImage(histogramImg)
		canvasHistogram.FillMode = canvas.ImageFillOriginal

		// Добавляем гистограмму в контейнер
		content := container.NewVBox(
			canvasHistogram,
			widget.NewLabel(""),
		)

		// Создаём диалог
		customDialog := dialog.NewCustomWithoutButtons("Гистограмма", content, window)

		// Кнопка "OK"
		confirmButton := widget.NewButton("OK", func() {
			customDialog.Hide()
		})

		fixedSizeButton := container.NewGridWrap(
			fyne.NewSize(100, 35),
			confirmButton,
		)

		centeredButton := container.NewCenter(fixedSizeButton)
		content.Add(centeredButton)

		customDialog.Resize(fyne.NewSize(420, 450)) // Немного увеличиваем размер окна
		customDialog.Show()
	})

	return button
}

func gammaButton(img *canvas.Image, window fyne.Window) fyne.CanvasObject {
	button := widget.NewButton("Gamma conversion", func() {
		if img.Image == nil {
			return
		}

		bounds := img.Image.Bounds()
		gammamizedImg := image.NewRGBA(bounds)

		gammaValue := widget.NewEntry()
		gammaValue.SetPlaceHolder("Число гамма")
		gammaValue.SetText("1")

		content := container.NewVBox(
			gammaValue,
			widget.NewLabel(""),
		)

		customDialog := dialog.NewCustomWithoutButtons("Gamma conversion", content, window)

		confirmButton := widget.NewButton("OK", func() {

			number := gammaValue.Text

			var negativeCeiling float64
			var err error
			var valueError bool

			if len(number) > 1 && string(number[1]) == "/" {
				var firstValue, secondValue int
				firstValue, err = strconv.Atoi(number[:1])
				if err != nil {
					valueError = true
				}
				secondValue, _ = strconv.Atoi(number[2:])
				if err != nil || secondValue == 0 {
					valueError = true
				}
				if !valueError {
					negativeCeiling = float64(firstValue) / float64(secondValue)
				}
			} else {
				var value int
				value, err = strconv.Atoi(number)
				if err != nil {
					valueError = true
				}
				if !valueError {
					negativeCeiling = float64(value)
				}
			}

			if negativeCeiling <= 0 || negativeCeiling > 255 {
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

					value := negativeCeiling

					newR := checkForLimit(255. * math.Pow(float64(r)/255., float64(value)))
					newG := checkForLimit(255. * math.Pow(float64(g)/255., float64(value)))
					newB := checkForLimit(255. * math.Pow(float64(b)/255., float64(value)))

					gammamizedImg.Set(x, y, color.RGBA{uint8(newR), uint8(newG), uint8(newB), uint8(a)})

				}

			}

			img.Image = gammamizedImg
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

		customDialog.Resize(fyne.NewSize(300, 100))

		customDialog.Show()

	})

	return button
}

func quantizationButton(img *canvas.Image, window fyne.Window) fyne.CanvasObject {
	button := widget.NewButton("Quantization", func() {
		if img.Image == nil {
			return
		}

		bounds := img.Image.Bounds()
		gammamizedImg := image.NewRGBA(bounds)

		numberOfQuantsSlider := widget.NewSlider(1, 255)
		numberOfQuantsSlider.Value = 1
		numberOfQuantsSlider.Step = 1

		valueLabel := widget.NewLabel("Quants value: " + strconv.Itoa(int(numberOfQuantsSlider.Value)))

		numberOfQuantsSlider.OnChanged = func(value float64) {
			valueLabel.SetText("Quants value: " + strconv.Itoa(int(value)))
		}

		content := container.NewVBox(
			numberOfQuantsSlider,
			valueLabel,
			widget.NewLabel(""),
		)

		customDialog := dialog.NewCustomWithoutButtons("Qunatization", content, window)

		confirmButton := widget.NewButton("OK", func() {

			var quantsArray [256]uint8

			quants := int(numberOfQuantsSlider.Value)

			if quants <= 0 || quants > 255 {
				dialog.ShowInformation("Ошибка", "Введите корректное число", window)
				return
			}

			customDialog.Hide()

			quantsSize := int(math.Ceil(256. / float64(quants)))

			for i := range quants {
				value := uint8(checkForLimit(float64(quantsSize - 1 + quantsSize*i)))

				start_point := quantsSize * i
				end_point := min(quantsSize+quantsSize*i, 256)

				for i := start_point; i < end_point; i++ {
					quantsArray[i] = value
				}
			}

			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					currentPixelColor := img.Image.At(x, y)
					r, g, b, a := currentPixelColor.RGBA()
					r = r >> 8
					g = g >> 8
					b = b >> 8
					a = a >> 8

					gammamizedImg.Set(x, y, color.RGBA{quantsArray[r], quantsArray[g], quantsArray[b], uint8(a)})

				}

			}

			img.Image = gammamizedImg
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

		customDialog.Resize(fyne.NewSize(300, 100))

		customDialog.Show()

	})

	return button
}

func solarizationButton(img *canvas.Image, window fyne.Window) fyne.CanvasObject {
	button := widget.NewButton("Solarization", func() {
		if img.Image == nil {
			return
		}

		bounds := img.Image.Bounds()
		solarizedImg := image.NewRGBA(bounds)

		solarizedSlider := widget.NewSlider(0, 0.05)
		solarizedSlider.Value = 4. / 255.
		solarizedSlider.Step = .00001

		value := strconv.FormatFloat(solarizedSlider.Value, 'f', -1, 64)
		valueLabel := widget.NewLabel("Value: " + value)

		solarizedSlider.OnChanged = func(value float64) {
			valueLabel.SetText("Value: " + strconv.FormatFloat(value, 'f', -1, 64))
		}

		content := container.NewVBox(
			valueLabel,
			solarizedSlider,
			widget.NewLabel(""),
		)

		customDialog := dialog.NewCustomWithoutButtons("Solarization", content, window)

		confirmButton := widget.NewButton("OK", func() {

			number := solarizedSlider.Value

			customDialog.Hide()

			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					currentPixelColor := img.Image.At(x, y)
					r, g, b, a := currentPixelColor.RGBA()
					r = r >> 8
					g = g >> 8
					b = b >> 8
					a = a >> 8

					newR := checkForLimit(float64((number) * float64(r*(255-r))))
					newG := checkForLimit(float64((number) * float64(g*(255-g))))
					newB := checkForLimit(float64((number) * float64(b*(255-b))))

					solarizedImg.Set(x, y, color.RGBA{uint8(newR), uint8(newG), uint8(newB), uint8(a)})

				}

			}

			img.Image = solarizedImg
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

		customDialog.Resize(fyne.NewSize(300, 100))

		customDialog.Show()

	})
	return button
}

func lowFreqFilterButton(img *canvas.Image, window fyne.Window) fyne.CanvasObject {
	button := widget.NewButton("Low Freq Filter", func() {
		if img.Image == nil {
			return
		}

		bounds := img.Image.Bounds()
		gammamizedImg := image.NewRGBA(bounds)

		H1Button := widget.NewButton("H1", func() {

			for y := 2; y < bounds.Max.Y-2; y++ {
				for x := 2; x < bounds.Max.X-2; x++ {
					_, _, _, newA := img.Image.At(x, y).RGBA()
					newA = newA >> 8

					firstPixelR, _, _, _ := img.Image.At(x, y).RGBA()
					firstPixelR = firstPixelR >> 8
					secondPixelR, _, _, _ := img.Image.At(x+1, y).RGBA()
					secondPixelR = secondPixelR >> 8
					thirdPixelR, _, _, _ := img.Image.At(x-1, y).RGBA()
					thirdPixelR = thirdPixelR >> 8

					fourthPixelR, _, _, _ := img.Image.At(x, y+1).RGBA()
					fourthPixelR = fourthPixelR >> 8
					fifthPixelR, _, _, _ := img.Image.At(x+1, y+1).RGBA()
					fifthPixelR = fifthPixelR >> 8
					sixthPixelR, _, _, _ := img.Image.At(x-1, y+1).RGBA()
					sixthPixelR = sixthPixelR >> 8

					seventhPixelR, _, _, _ := img.Image.At(x, y-1).RGBA()
					seventhPixelR = seventhPixelR >> 8
					eighthPixelR, _, _, _ := img.Image.At(x+1, y-1).RGBA()
					eighthPixelR = eighthPixelR >> 8
					ninthPixelR, _, _, _ := img.Image.At(x-1, y-1).RGBA()
					ninthPixelR = ninthPixelR >> 8

					valueR := firstPixelR + secondPixelR + thirdPixelR + fourthPixelR + fifthPixelR + sixthPixelR + seventhPixelR + eighthPixelR + ninthPixelR

					_, firstPixelG, _, _ := img.Image.At(x, y).RGBA()
					firstPixelG = firstPixelG >> 8
					_, secondPixelG, _, _ := img.Image.At(x+1, y).RGBA()
					secondPixelG = secondPixelG >> 8
					_, thirdPixelG, _, _ := img.Image.At(x-1, y).RGBA()
					thirdPixelG = thirdPixelG >> 8

					_, fourthPixelG, _, _ := img.Image.At(x, y+1).RGBA()
					fourthPixelG = fourthPixelG >> 8
					_, fifthPixelG, _, _ := img.Image.At(x+1, y+1).RGBA()
					fifthPixelG = fifthPixelG >> 8
					_, sixthPixelG, _, _ := img.Image.At(x-1, y+1).RGBA()
					sixthPixelG = sixthPixelG >> 8

					_, seventhPixelG, _, _ := img.Image.At(x, y-1).RGBA()
					seventhPixelG = seventhPixelG >> 8
					_, eighthPixelG, _, _ := img.Image.At(x+1, y-1).RGBA()
					eighthPixelG = eighthPixelG >> 8
					_, ninthPixelG, _, _ := img.Image.At(x-1, y-1).RGBA()
					ninthPixelG = ninthPixelG >> 8

					valueG := firstPixelG + secondPixelG + thirdPixelG + fourthPixelG + fifthPixelG + sixthPixelG + seventhPixelG + eighthPixelG + ninthPixelG

					_, _, firstPixelB, _ := img.Image.At(x, y).RGBA()
					firstPixelB = firstPixelB >> 8
					_, _, secondPixelB, _ := img.Image.At(x+1, y).RGBA()
					secondPixelB = secondPixelB >> 8
					_, _, thirdPixelB, _ := img.Image.At(x-1, y).RGBA()
					thirdPixelB = thirdPixelB >> 8

					_, _, fourthPixelB, _ := img.Image.At(x, y+1).RGBA()
					fourthPixelB = fourthPixelB >> 8
					_, _, fifthPixelB, _ := img.Image.At(x+1, y+1).RGBA()
					fifthPixelB = fifthPixelB >> 8
					_, _, sixthPixelB, _ := img.Image.At(x-1, y+1).RGBA()
					sixthPixelB = sixthPixelB >> 8

					_, _, seventhPixelB, _ := img.Image.At(x, y-1).RGBA()
					seventhPixelB = seventhPixelB >> 8
					_, _, eighthPixelB, _ := img.Image.At(x+1, y-1).RGBA()
					eighthPixelB = eighthPixelB >> 8
					_, _, ninthPixelB, _ := img.Image.At(x-1, y-1).RGBA()
					ninthPixelB = ninthPixelB >> 8

					valueB := firstPixelB + secondPixelB + thirdPixelB + fourthPixelB + fifthPixelB + sixthPixelB + seventhPixelB + eighthPixelB + ninthPixelB

					gammamizedImg.Set(x, y, color.RGBA{uint8(valueR / 9), uint8(valueG / 9), uint8(valueB / 9), uint8(newA)})
				}
			}
			img.Image = gammamizedImg
			img.Refresh()
		})

		H2Button := widget.NewButton("H2", func() {

			for y := 2; y < bounds.Max.Y-2; y++ {
				for x := 2; x < bounds.Max.X-2; x++ {
					_, _, _, newA := img.Image.At(x, y).RGBA()
					newA = newA >> 8

					firstPixelR, _, _, _ := img.Image.At(x, y).RGBA()
					firstPixelR = firstPixelR >> 8
					firstPixelR = firstPixelR * 2
					secondPixelR, _, _, _ := img.Image.At(x+1, y).RGBA()
					secondPixelR = secondPixelR >> 8
					thirdPixelR, _, _, _ := img.Image.At(x-1, y).RGBA()
					thirdPixelR = thirdPixelR >> 8

					fourthPixelR, _, _, _ := img.Image.At(x, y+1).RGBA()
					fourthPixelR = fourthPixelR >> 8
					fifthPixelR, _, _, _ := img.Image.At(x+1, y+1).RGBA()
					fifthPixelR = fifthPixelR >> 8
					sixthPixelR, _, _, _ := img.Image.At(x-1, y+1).RGBA()
					sixthPixelR = sixthPixelR >> 8

					seventhPixelR, _, _, _ := img.Image.At(x, y-1).RGBA()
					seventhPixelR = seventhPixelR >> 8
					eighthPixelR, _, _, _ := img.Image.At(x+1, y-1).RGBA()
					eighthPixelR = eighthPixelR >> 8
					ninthPixelR, _, _, _ := img.Image.At(x-1, y-1).RGBA()
					ninthPixelR = ninthPixelR >> 8

					valueR := firstPixelR + secondPixelR + thirdPixelR + fourthPixelR + fifthPixelR + sixthPixelR + seventhPixelR + eighthPixelR + ninthPixelR

					_, firstPixelG, _, _ := img.Image.At(x, y).RGBA()
					firstPixelG = firstPixelG >> 8
					firstPixelG = firstPixelG * 2
					_, secondPixelG, _, _ := img.Image.At(x+1, y).RGBA()
					secondPixelG = secondPixelG >> 8
					_, thirdPixelG, _, _ := img.Image.At(x-1, y).RGBA()
					thirdPixelG = thirdPixelG >> 8

					_, fourthPixelG, _, _ := img.Image.At(x, y+1).RGBA()
					fourthPixelG = fourthPixelG >> 8
					_, fifthPixelG, _, _ := img.Image.At(x+1, y+1).RGBA()
					fifthPixelG = fifthPixelG >> 8
					_, sixthPixelG, _, _ := img.Image.At(x-1, y+1).RGBA()
					sixthPixelG = sixthPixelG >> 8

					_, seventhPixelG, _, _ := img.Image.At(x, y-1).RGBA()
					seventhPixelG = seventhPixelG >> 8
					_, eighthPixelG, _, _ := img.Image.At(x+1, y-1).RGBA()
					eighthPixelG = eighthPixelG >> 8
					_, ninthPixelG, _, _ := img.Image.At(x-1, y-1).RGBA()
					ninthPixelG = ninthPixelG >> 8

					valueG := firstPixelG + secondPixelG + thirdPixelG + fourthPixelG + fifthPixelG + sixthPixelG + seventhPixelG + eighthPixelG + ninthPixelG

					_, _, firstPixelB, _ := img.Image.At(x, y).RGBA()
					firstPixelB = firstPixelB >> 8
					firstPixelB = firstPixelB * 2
					_, _, secondPixelB, _ := img.Image.At(x+1, y).RGBA()
					secondPixelB = secondPixelB >> 8
					_, _, thirdPixelB, _ := img.Image.At(x-1, y).RGBA()
					thirdPixelB = thirdPixelB >> 8

					_, _, fourthPixelB, _ := img.Image.At(x, y+1).RGBA()
					fourthPixelB = fourthPixelB >> 8
					_, _, fifthPixelB, _ := img.Image.At(x+1, y+1).RGBA()
					fifthPixelB = fifthPixelB >> 8
					_, _, sixthPixelB, _ := img.Image.At(x-1, y+1).RGBA()
					sixthPixelB = sixthPixelB >> 8

					_, _, seventhPixelB, _ := img.Image.At(x, y-1).RGBA()
					seventhPixelB = seventhPixelB >> 8
					_, _, eighthPixelB, _ := img.Image.At(x+1, y-1).RGBA()
					eighthPixelB = eighthPixelB >> 8
					_, _, ninthPixelB, _ := img.Image.At(x-1, y-1).RGBA()
					ninthPixelB = ninthPixelB >> 8

					valueB := firstPixelB + secondPixelB + thirdPixelB + fourthPixelB + fifthPixelB + sixthPixelB + seventhPixelB + eighthPixelB + ninthPixelB

					gammamizedImg.Set(x, y, color.RGBA{uint8(valueR / 10), uint8(valueG / 10), uint8(valueB / 10), uint8(newA)})
				}
			}
			img.Image = gammamizedImg
			img.Refresh()
		})

		H3Button := widget.NewButton("H3", func() {

			for y := 2; y < bounds.Max.Y-2; y++ {
				for x := 2; x < bounds.Max.X-2; x++ {
					_, _, _, newA := img.Image.At(x, y).RGBA()
					newA = newA >> 8

					firstPixelR, _, _, _ := img.Image.At(x, y).RGBA()
					firstPixelR = firstPixelR >> 8
					firstPixelR = firstPixelR * 4
					secondPixelR, _, _, _ := img.Image.At(x+1, y).RGBA()
					secondPixelR = secondPixelR >> 8
					secondPixelR = secondPixelR * 2
					thirdPixelR, _, _, _ := img.Image.At(x-1, y).RGBA()
					thirdPixelR = thirdPixelR >> 8
					thirdPixelR = thirdPixelR * 2

					fourthPixelR, _, _, _ := img.Image.At(x, y+1).RGBA()
					fourthPixelR = fourthPixelR >> 8
					fourthPixelR = fourthPixelR * 2
					fifthPixelR, _, _, _ := img.Image.At(x+1, y+1).RGBA()
					fifthPixelR = fifthPixelR >> 8
					sixthPixelR, _, _, _ := img.Image.At(x-1, y+1).RGBA()
					sixthPixelR = sixthPixelR >> 8

					seventhPixelR, _, _, _ := img.Image.At(x, y-1).RGBA()
					seventhPixelR = seventhPixelR >> 8
					seventhPixelR = seventhPixelR * 2
					eighthPixelR, _, _, _ := img.Image.At(x+1, y-1).RGBA()
					eighthPixelR = eighthPixelR >> 8
					ninthPixelR, _, _, _ := img.Image.At(x-1, y-1).RGBA()
					ninthPixelR = ninthPixelR >> 8

					valueR := firstPixelR + secondPixelR + thirdPixelR + fourthPixelR + fifthPixelR + sixthPixelR + seventhPixelR + eighthPixelR + ninthPixelR

					_, firstPixelG, _, _ := img.Image.At(x, y).RGBA()
					firstPixelG = firstPixelG >> 8
					firstPixelG = firstPixelG * 4
					_, secondPixelG, _, _ := img.Image.At(x+1, y).RGBA()
					secondPixelG = secondPixelG >> 8
					secondPixelG = secondPixelG * 2
					_, thirdPixelG, _, _ := img.Image.At(x-1, y).RGBA()
					thirdPixelG = thirdPixelG >> 8
					thirdPixelG = thirdPixelG * 2

					_, fourthPixelG, _, _ := img.Image.At(x, y+1).RGBA()
					fourthPixelG = fourthPixelG >> 8
					fourthPixelG = fourthPixelG * 2
					_, fifthPixelG, _, _ := img.Image.At(x+1, y+1).RGBA()
					fifthPixelG = fifthPixelG >> 8
					_, sixthPixelG, _, _ := img.Image.At(x-1, y+1).RGBA()
					sixthPixelG = sixthPixelG >> 8

					_, seventhPixelG, _, _ := img.Image.At(x, y-1).RGBA()
					seventhPixelG = seventhPixelG >> 8
					seventhPixelG = seventhPixelG * 2
					_, eighthPixelG, _, _ := img.Image.At(x+1, y-1).RGBA()
					eighthPixelG = eighthPixelG >> 8
					_, ninthPixelG, _, _ := img.Image.At(x-1, y-1).RGBA()
					ninthPixelG = ninthPixelG >> 8

					valueG := firstPixelG + secondPixelG + thirdPixelG + fourthPixelG + fifthPixelG + sixthPixelG + seventhPixelG + eighthPixelG + ninthPixelG

					_, _, firstPixelB, _ := img.Image.At(x, y).RGBA()
					firstPixelB = firstPixelB >> 8
					firstPixelB = firstPixelB * 4
					_, _, secondPixelB, _ := img.Image.At(x+1, y).RGBA()
					secondPixelB = secondPixelB >> 8
					secondPixelB = secondPixelB * 2
					_, _, thirdPixelB, _ := img.Image.At(x-1, y).RGBA()
					thirdPixelB = thirdPixelB >> 8
					thirdPixelB = thirdPixelB * 2

					_, _, fourthPixelB, _ := img.Image.At(x, y+1).RGBA()
					fourthPixelB = fourthPixelB >> 8
					fourthPixelB = fourthPixelB * 2
					_, _, fifthPixelB, _ := img.Image.At(x+1, y+1).RGBA()
					fifthPixelB = fifthPixelB >> 8
					_, _, sixthPixelB, _ := img.Image.At(x-1, y+1).RGBA()
					sixthPixelB = sixthPixelB >> 8

					_, _, seventhPixelB, _ := img.Image.At(x, y-1).RGBA()
					seventhPixelB = seventhPixelB >> 8
					seventhPixelB = seventhPixelB * 2
					_, _, eighthPixelB, _ := img.Image.At(x+1, y-1).RGBA()
					eighthPixelB = eighthPixelB >> 8
					_, _, ninthPixelB, _ := img.Image.At(x-1, y-1).RGBA()
					ninthPixelB = ninthPixelB >> 8

					valueB := firstPixelB + secondPixelB + thirdPixelB + fourthPixelB + fifthPixelB + sixthPixelB + seventhPixelB + eighthPixelB + ninthPixelB

					gammamizedImg.Set(x, y, color.RGBA{uint8(valueR / 16), uint8(valueG / 16), uint8(valueB / 16), uint8(newA)})
				}
			}
			img.Image = gammamizedImg
			img.Refresh()
		})

		content := container.NewVBox(
			H1Button,
			H2Button,
			H3Button,
			widget.NewLabel(""),
		)

		customDialog := dialog.NewCustomWithoutButtons("Qunatization", content, window)

		dissmisButton := widget.NewButton("Cancel", func() {
			customDialog.Hide()
		})

		fixedSizeButton := container.NewGridWrap(
			fyne.NewSize(100, 35),
			dissmisButton,
		)

		centeredButton := container.NewCenter(fixedSizeButton)
		content.Add(centeredButton)

		customDialog.Resize(fyne.NewSize(300, 100))

		customDialog.Show()

	})

	return button
}
