package main

import (
	"image"
	"image/color"
	"math"
	"sort"
	"strconv"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func checkForLimit(value float64) float64 {
	if value > 255. {
		return 255.
	} else if value < 0. {
		return 0.
	}
	return value
}

func pascalRow(n int) []float64 {
	res := make([]float64, n)
	elem := 1.
	res[0] = 1.
	res[n-1] = 1.
	end := float64(n)/2. + 1

	for k := 1.; k < end; k++ {
		elem *= (float64(n-1) + 1 - k) / k

		pos := int(k)

		res[pos] = elem
		res[n-pos-1] = elem
	}

	return res
}

func GaussKernelByPascalRow(row []float64) [][]float64 {
	res := make([][]float64, 0, len(row))

	sum := 0.

	end := len(row)/2 + 1

	for range row {
		res = append(res, make([]float64, len(row)))
	}

	n := len(row) - 1

	setVal := func(i, j int, val float64) {
		res[i][j] = val
		res[n-i][j] = val
		res[i][n-j] = val
		res[n-i][n-j] = val
	}

	for i := range end {
		for j := range end {
			val := float64(row[i] * row[j])

			setVal(i, j, val)

			sum += val

			var t1, t2 bool

			if i != n-i {
				t1 = true
				sum += val
			}

			if j != n-j {
				t2 = true
				sum += val
			}

			if t1 && t2 {
				sum += val
			}
		}
	}

	for i := range end {
		for j := range end {
			val := res[i][j] / sum

			setVal(i, j, val)
		}
	}

	return res
}

func NewOriginalButton(img *canvas.Image, origImg *canvas.Image) fyne.CanvasObject {
	button := widget.NewButton("Original", func() {
		if img.Image == nil || origImg == nil {
			return
		}
		img.Image = origImg.Image
		img.Refresh()
	})

	return button
}

func NewGrayscaleButton(img *canvas.Image) fyne.CanvasObject {

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

func NewNegativeButton(img *canvas.Image, window fyne.Window) fyne.CanvasObject {
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

func NewAdjustBrightnessButton(img *canvas.Image, window fyne.Window) fyne.CanvasObject {
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

func NewBinarizationButton(img *canvas.Image, window fyne.Window) fyne.CanvasObject {
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

func NewIncreaseContrastButton(img *canvas.Image, window fyne.Window) fyne.CanvasObject {
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

func NewDecreaseContrastButton(img *canvas.Image, window fyne.Window) fyne.CanvasObject {
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

func NewCreateHistogramButton(img *canvas.Image, window fyne.Window) fyne.CanvasObject {
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

func NewGammaButton(img *canvas.Image, window fyne.Window) fyne.CanvasObject {
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

func NewQuantizationButton(img *canvas.Image, window fyne.Window) fyne.CanvasObject {
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

		customDialog := dialog.NewCustomWithoutButtons("Quantization", content, window)

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

func NewSolarizationButton(img *canvas.Image, window fyne.Window) fyne.CanvasObject {
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

func NewLowFreqFilterButton(img *canvas.Image, window fyne.Window) fyne.CanvasObject {

	button := widget.NewButton("Low Freq Filter", func() {
		if img.Image == nil {
			return
		}

		H1Button := widget.NewButton("H1", func() {
			bounds := img.Image.Bounds()
			highFreqImg := image.NewRGBA(bounds)

			for y := 1; y < bounds.Max.Y-1; y++ {
				for x := 1; x < bounds.Max.X-1; x++ {
					var sumR, sumG, sumB uint32

					for ky := -1; ky <= 1; ky++ {
						for kx := -1; kx <= 1; kx++ {

							px := x + kx
							py := y + ky

							r, g, b, _ := img.Image.At(px, py).RGBA()

							sumR += (r >> 8)
							sumG += (g >> 8)
							sumB += (b >> 8)
						}
					}

					highFreqImg.Set(x, y, color.RGBA{
						R: uint8(sumR / 9),
						G: uint8(sumG / 9),
						B: uint8(sumB / 9),
						A: 255,
					})
				}
			}

			img.Image = highFreqImg
			img.Refresh()
		})

		H2Button := widget.NewButton("H2", func() {
			if img.Image == nil {
				return
			}

			bounds := img.Image.Bounds()
			highFreqImg := image.NewRGBA(bounds)

			kernel := [3][3]uint32{
				{1, 1, 1},
				{1, 2, 1},
				{1, 1, 1},
			}

			for y := 1; y < bounds.Max.Y-1; y++ {
				for x := 1; x < bounds.Max.X-1; x++ {
					var sumR, sumG, sumB uint32

					for ky := -1; ky <= 1; ky++ {
						for kx := -1; kx <= 1; kx++ {

							px := x + kx
							py := y + ky

							r, g, b, _ := img.Image.At(px, py).RGBA()
							factor := kernel[ky+1][kx+1]

							sumR += (r >> 8) * factor
							sumG += (g >> 8) * factor
							sumB += (b >> 8) * factor
						}
					}

					highFreqImg.Set(x, y, color.RGBA{
						R: uint8(sumR / 10),
						G: uint8(sumG / 10),
						B: uint8(sumB / 10),
						A: 255,
					})
				}
			}

			img.Image = highFreqImg
			img.Refresh()
		})

		H3Button := widget.NewButton("H3", func() {
			if img.Image == nil {
				return
			}

			bounds := img.Image.Bounds()
			highFreqImg := image.NewRGBA(bounds)

			kernel := [3][3]uint32{
				{1, 2, 1},
				{2, 4, 2},
				{1, 2, 1},
			}

			for y := 1; y < bounds.Max.Y-1; y++ {
				for x := 1; x < bounds.Max.X-1; x++ {
					var sumR, sumG, sumB uint32

					for ky := -1; ky <= 1; ky++ {
						for kx := -1; kx <= 1; kx++ {

							px := x + kx
							py := y + ky

							r, g, b, _ := img.Image.At(px, py).RGBA()
							factor := kernel[ky+1][kx+1]

							sumR += (r >> 8) * factor
							sumG += (g >> 8) * factor
							sumB += (b >> 8) * factor
						}
					}

					highFreqImg.Set(x, y, color.RGBA{
						R: uint8(sumR / 16),
						G: uint8(sumG / 16),
						B: uint8(sumB / 16),
						A: 255,
					})
				}
			}

			img.Image = highFreqImg
			img.Refresh()
		})

		content := container.NewVBox(
			H1Button,
			H2Button,
			H3Button,
			widget.NewLabel(""),
		)

		customDialog := dialog.NewCustomWithoutButtons("Low Freq", content, window)

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

func NewHighFreqFilterButton(img *canvas.Image, window fyne.Window) fyne.CanvasObject {

	button := widget.NewButton("High Freq Filter", func() {
		if img.Image == nil {
			return
		}

		H1Button := widget.NewButton("H1", func() {
			bounds := img.Image.Bounds()
			highFreqImg := image.NewRGBA(bounds)

			kernel := [3][3]float64{
				{-1, -1, -1},
				{-1, 9, -1},
				{-1, -1, -1},
			}

			for y := 1; y < bounds.Max.Y-1; y++ {
				for x := 1; x < bounds.Max.X-1; x++ {
					var sumR, sumG, sumB float64

					for ky := -1; ky <= 1; ky++ {
						for kx := -1; kx <= 1; kx++ {

							px := x + kx
							py := y + ky

							r, g, b, _ := img.Image.At(px, py).RGBA()
							factor := kernel[ky+1][kx+1]

							sumR += float64(r>>8) * factor
							sumG += float64(g>>8) * factor
							sumB += float64(b>>8) * factor
						}
					}

					fr := checkForLimit(sumR)
					fg := checkForLimit(sumG)
					fb := checkForLimit(sumB)

					highFreqImg.Set(x, y, color.RGBA{
						R: uint8(fr),
						G: uint8(fg),
						B: uint8(fb),
						A: 255,
					})
				}
			}

			img.Image = highFreqImg
			img.Refresh()
		})

		H2Button := widget.NewButton("H2", func() {
			if img.Image == nil {
				return
			}

			bounds := img.Image.Bounds()
			highFreqImg := image.NewRGBA(bounds)

			kernel := [3][3]float64{
				{0, -1, 0},
				{-1, 5, -1},
				{0, -1, 0},
			}

			for y := 1; y < bounds.Max.Y-1; y++ {
				for x := 1; x < bounds.Max.X-1; x++ {
					var sumR, sumG, sumB float64

					for ky := -1; ky <= 1; ky++ {
						for kx := -1; kx <= 1; kx++ {

							px := x + kx
							py := y + ky

							r, g, b, _ := img.Image.At(px, py).RGBA()
							factor := kernel[ky+1][kx+1]

							sumR += float64(r>>8) * factor
							sumG += float64(g>>8) * factor
							sumB += float64(b>>8) * factor
						}
					}

					fr := checkForLimit(sumR)
					fg := checkForLimit(sumG)
					fb := checkForLimit(sumB)

					highFreqImg.Set(x, y, color.RGBA{
						R: uint8(fr),
						G: uint8(fg),
						B: uint8(fb),
						A: 255,
					})
				}
			}

			img.Image = highFreqImg
			img.Refresh()
		})

		H3Button := widget.NewButton("H3", func() {
			if img.Image == nil {
				return
			}

			bounds := img.Image.Bounds()
			highFreqImg := image.NewRGBA(bounds)

			kernel := [3][3]float64{
				{1, -2, 1},
				{-2, 5, -2},
				{1, -2, 1},
			}

			for y := 1; y < bounds.Max.Y-1; y++ {
				for x := 1; x < bounds.Max.X-1; x++ {
					var sumR, sumG, sumB float64

					for ky := -1; ky <= 1; ky++ {
						for kx := -1; kx <= 1; kx++ {

							px := x + kx
							py := y + ky

							r, g, b, _ := img.Image.At(px, py).RGBA()
							factor := kernel[ky+1][kx+1]

							sumR += float64(r>>8) * factor
							sumG += float64(g>>8) * factor
							sumB += float64(b>>8) * factor
						}
					}

					fr := checkForLimit(sumR)
					fg := checkForLimit(sumG)
					fb := checkForLimit(sumB)

					highFreqImg.Set(x, y, color.RGBA{
						R: uint8(fr),
						G: uint8(fg),
						B: uint8(fb),
						A: 255,
					})
				}
			}

			img.Image = highFreqImg
			img.Refresh()
		})

		content := container.NewVBox(
			H1Button,
			H2Button,
			H3Button,
			widget.NewLabel(""),
		)

		customDialog := dialog.NewCustomWithoutButtons("High Freq", content, window)

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

func NewMedianFilterButton(img *canvas.Image, window fyne.Window) fyne.CanvasObject {
	button := widget.NewButton("Median filter", func() {
		if img.Image == nil {
			return
		}

		bounds := img.Image.Bounds()
		changedImg := image.NewRGBA(bounds)

		sizeOfWindow := widget.NewEntry()
		sizeOfWindow.SetPlaceHolder("Size of window")

		content := container.NewVBox(
			sizeOfWindow,
			widget.NewLabel(""),
		)

		customDialog := dialog.NewCustomWithoutButtons("Median filter", content, window)

		confirmButton := widget.NewButton("OK", func() {

			size := sizeOfWindow.Text
			windowSize, err := strconv.Atoi(size)
			if err != nil || windowSize%2 == 0 {
				dialog.ShowInformation("Ошибка", "Введите корректное нечетное число", window)
				return
			}

			customDialog.Hide()

			window := windowSize / 2

			for y := window; y < bounds.Max.Y-window; y++ {
				for x := window; x < bounds.Max.X-window; x++ {

					var r, g, b []int
					for ky := -window; ky <= window; ky++ {
						for kx := -window; kx <= window; kx++ {

							pr, pg, pb, _ := img.Image.At(x+kx, y+ky).RGBA()
							r = append(r, int(pr>>8))
							g = append(g, int(pg>>8))
							b = append(b, int(pb>>8))
						}
					}

					sort.Ints(r)
					sort.Ints(g)
					sort.Ints(b)

					medianPos := len(r) / 2

					changedImg.Set(x, y, color.RGBA{
						R: uint8(r[medianPos]),
						G: uint8(g[medianPos]),
						B: uint8(b[medianPos]),
						A: 255,
					})
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

func NewGaussBlurButton(img *canvas.Image, window fyne.Window) fyne.CanvasObject {
	button := widget.NewButton("Gauss blur", func() {
		if img.Image == nil {
			return
		}

		bounds := img.Image.Bounds()
		changedImg := image.NewRGBA(bounds)

		sizeOfWindow := widget.NewEntry()
		sizeOfWindow.SetPlaceHolder("Size of window")

		content := container.NewVBox(
			sizeOfWindow,
			widget.NewLabel(""),
		)

		customDialog := dialog.NewCustomWithoutButtons("Gauss blur", content, window)

		confirmButton := widget.NewButton("OK", func() {

			size := sizeOfWindow.Text
			windowSize, err := strconv.Atoi(size)
			if err != nil || windowSize%2 == 0 {
				dialog.ShowInformation("Ошибка", "Введите корректное нечетное число", window)
				return
			}

			customDialog.Hide()

			kernel := GaussKernelByPascalRow(pascalRow(windowSize))

			var multiplier float64 = float64(windowSize) / 3.

			for y := windowSize; y < bounds.Max.Y-windowSize; y++ {
				for x := windowSize; x < bounds.Max.X-windowSize; x++ {
					var sumR, sumG, sumB, sumWeight float64
					for ky := -math.Round(multiplier); ky <= math.Round(multiplier); ky++ {
						for kx := -math.Round(multiplier); kx <= math.Round(multiplier); kx++ {
							weight := kernel[int(ky)+int(math.Round(multiplier))][int(kx)+int(math.Round(multiplier))]
							r, g, b, _ := img.Image.At(x+int(kx), y+int(ky)).RGBA()
							sumR += float64(r>>8) * weight
							sumG += float64(g>>8) * weight
							sumB += float64(b>>8) * weight
							sumWeight += weight
						}
					}

					changedImg.Set(x, y, color.RGBA{
						R: uint8(sumR / sumWeight),
						G: uint8(sumG / sumWeight),
						B: uint8(sumB / sumWeight),
						A: 255,
					})
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

func NewEdgeEmpowerButton(img *canvas.Image, window fyne.Window) fyne.CanvasObject {

	button := widget.NewButton("Edge empower", func() {

		if img.Image == nil {
			return
		}

		bounds := img.Image.Bounds()
		highFreqImg := image.NewRGBA(bounds)

		kernel := [3][3]float64{
			{0, 1, 0},
			{1, -4, 1},
			{0, 1, 0},
		}

		for y := 1; y < bounds.Max.Y-1; y++ {
			for x := 1; x < bounds.Max.X-1; x++ {
				var sumR, sumG, sumB float64

				for ky := -1; ky <= 1; ky++ {
					for kx := -1; kx <= 1; kx++ {

						px := x + kx
						py := y + ky

						r, g, b, _ := img.Image.At(px, py).RGBA()
						factor := kernel[ky+1][kx+1]

						grayScale := (0.3*float64(r>>8) + 0.59*float64(g>>8) + 0.11*float64(b>>8))

						sumR += grayScale * factor
						sumG += grayScale * factor
						sumB += grayScale * factor
					}
				}

				fr := math.Abs(sumR)
				fg := math.Abs(sumG)
				fb := math.Abs(sumB)

				highFreqImg.Set(x, y, color.RGBA{
					R: uint8(fr),
					G: uint8(fg),
					B: uint8(fb),
					A: 255,
				})
			}
		}

		img.Image = highFreqImg
		img.Refresh()
	})

	return button
}

func NewShiftEdgeButton(img *canvas.Image, window fyne.Window) fyne.CanvasObject {

	button := widget.NewButton("Shift edge", func() {
		if img.Image == nil {
			return
		}

		H1Button := widget.NewButton("Vertical", func() {
			bounds := img.Image.Bounds()
			highFreqImg := image.NewRGBA(bounds)

			kernel := [3][3]float64{
				{0, 0, 0},
				{-1, 1, 0},
				{0, 0, 0},
			}

			for y := 1; y < bounds.Max.Y-1; y++ {
				for x := 1; x < bounds.Max.X-1; x++ {
					var sumR, sumG, sumB float64

					for ky := -1; ky <= 1; ky++ {
						for kx := -1; kx <= 1; kx++ {

							px := x + kx
							py := y + ky

							r, g, b, _ := img.Image.At(px, py).RGBA()
							factor := kernel[ky+1][kx+1]

							grayScale := (0.3*float64(r>>8) + 0.59*float64(g>>8) + 0.11*float64(b>>8))

							sumR += grayScale * factor
							sumG += grayScale * factor
							sumB += grayScale * factor
						}
					}

					fr := checkForLimit(sumR)
					fg := checkForLimit(sumG)
					fb := checkForLimit(sumB)

					highFreqImg.Set(x, y, color.RGBA{
						R: uint8(fr),
						G: uint8(fg),
						B: uint8(fb),
						A: 255,
					})
				}
			}

			img.Image = highFreqImg
			img.Refresh()
		})

		H2Button := widget.NewButton("Horizontal", func() {
			if img.Image == nil {
				return
			}

			bounds := img.Image.Bounds()
			highFreqImg := image.NewRGBA(bounds)

			kernel := [3][3]float64{
				{0, -1, 0},
				{0, 1, 0},
				{0, 0, 0},
			}

			for y := 1; y < bounds.Max.Y-1; y++ {
				for x := 1; x < bounds.Max.X-1; x++ {
					var sumR, sumG, sumB float64

					for ky := -1; ky <= 1; ky++ {
						for kx := -1; kx <= 1; kx++ {

							px := x + kx
							py := y + ky

							r, g, b, _ := img.Image.At(px, py).RGBA()
							factor := kernel[ky+1][kx+1]

							grayScale := (0.3*float64(r>>8) + 0.59*float64(g>>8) + 0.11*float64(b>>8))

							sumR += grayScale * factor
							sumG += grayScale * factor
							sumB += grayScale * factor
						}
					}

					fr := checkForLimit(sumR)
					fg := checkForLimit(sumG)
					fb := checkForLimit(sumB)

					highFreqImg.Set(x, y, color.RGBA{
						R: uint8(fr),
						G: uint8(fg),
						B: uint8(fb),
						A: 255,
					})
				}
			}

			img.Image = highFreqImg
			img.Refresh()
		})

		H3Button := widget.NewButton("Diagonal", func() {
			if img.Image == nil {
				return
			}

			bounds := img.Image.Bounds()
			highFreqImg := image.NewRGBA(bounds)

			kernel := [3][3]float64{
				{-1, 0, 0},
				{0, 1, 0},
				{0, 0, 0},
			}

			for y := 1; y < bounds.Max.Y-1; y++ {
				for x := 1; x < bounds.Max.X-1; x++ {
					var sumR, sumG, sumB float64

					for ky := -1; ky <= 1; ky++ {
						for kx := -1; kx <= 1; kx++ {

							px := x + kx
							py := y + ky

							r, g, b, _ := img.Image.At(px, py).RGBA()
							factor := kernel[ky+1][kx+1]

							grayScale := (0.3*float64(r>>8) + 0.59*float64(g>>8) + 0.11*float64(b>>8))

							sumR += grayScale * factor
							sumG += grayScale * factor
							sumB += grayScale * factor
						}
					}

					fr := math.Abs(sumR)
					fg := math.Abs(sumG)
					fb := math.Abs(sumB)

					highFreqImg.Set(x, y, color.RGBA{
						R: uint8(fr),
						G: uint8(fg),
						B: uint8(fb),
						A: 255,
					})
				}
			}

			img.Image = highFreqImg
			img.Refresh()
		})

		content := container.NewVBox(
			H1Button,
			H2Button,
			H3Button,
			widget.NewLabel(""),
		)

		customDialog := dialog.NewCustomWithoutButtons("Shift edge", content, window)

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

func NewEmbossingButton(img *canvas.Image, window fyne.Window) fyne.CanvasObject {

	button := widget.NewButton("Embossing", func() {
		if img.Image == nil {
			return
		}

		H1Button := widget.NewButton("In", func() {
			bounds := img.Image.Bounds()
			highFreqImg := image.NewRGBA(bounds)

			kernel := [3][3]float64{
				{0, 1, 0},
				{-1, 0, 1},
				{0, -1, 0},
			}

			for y := 1; y < bounds.Max.Y-1; y++ {
				for x := 1; x < bounds.Max.X-1; x++ {
					var sumR, sumG, sumB float64

					for ky := -1; ky <= 1; ky++ {
						for kx := -1; kx <= 1; kx++ {

							px := x + kx
							py := y + ky

							r, g, b, _ := img.Image.At(px, py).RGBA()
							factor := kernel[ky+1][kx+1]

							grayScale := (0.3*float64(r>>8) + 0.59*float64(g>>8) + 0.11*float64(b>>8))

							sumR += grayScale * factor
							sumG += grayScale * factor
							sumB += grayScale * factor
						}
					}

					fr := checkForLimit(sumR + 128)
					fg := checkForLimit(sumG + 128)
					fb := checkForLimit(sumB + 128)

					highFreqImg.Set(x, y, color.RGBA{
						R: uint8(fr),
						G: uint8(fg),
						B: uint8(fb),
						A: 255,
					})
				}
			}

			img.Image = highFreqImg
			img.Refresh()
		})

		H2Button := widget.NewButton("Out", func() {
			if img.Image == nil {
				return
			}

			bounds := img.Image.Bounds()
			highFreqImg := image.NewRGBA(bounds)

			kernel := [3][3]float64{
				{0, -1, 0},
				{1, 0, -1},
				{0, 1, 0},
			}

			for y := 1; y < bounds.Max.Y-1; y++ {
				for x := 1; x < bounds.Max.X-1; x++ {
					var sumR, sumG, sumB float64

					for ky := -1; ky <= 1; ky++ {
						for kx := -1; kx <= 1; kx++ {

							px := x + kx
							py := y + ky

							r, g, b, _ := img.Image.At(px, py).RGBA()
							factor := kernel[ky+1][kx+1]

							grayScale := (0.3*float64(r>>8) + 0.59*float64(g>>8) + 0.11*float64(b>>8))

							sumR += grayScale * factor
							sumG += grayScale * factor
							sumB += grayScale * factor
						}
					}

					fr := checkForLimit(sumR + 128)
					fg := checkForLimit(sumG + 128)
					fb := checkForLimit(sumB + 128)

					highFreqImg.Set(x, y, color.RGBA{
						R: uint8(fr),
						G: uint8(fg),
						B: uint8(fb),
						A: 255,
					})
				}
			}

			img.Image = highFreqImg
			img.Refresh()
		})

		content := container.NewVBox(
			H1Button,
			H2Button,
			widget.NewLabel(""),
		)

		customDialog := dialog.NewCustomWithoutButtons("Shift edge", content, window)

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

func NewKirschButton(img *canvas.Image, window fyne.Window) fyne.CanvasObject {

	button := widget.NewButton("Kirsch", func() {

		if img.Image == nil {
			return
		}

		bounds := img.Image.Bounds()
		highFreqImg := image.NewRGBA(bounds)

		kernels := [8][3][3]float64{
			{{5, 5, 5}, {-3, 0, -3}, {-3, -3, -3}}, // 0°
			{{-3, 5, 5}, {-3, 0, 5}, {-3, -3, -3}}, // 45°
			{{-3, -3, 5}, {-3, 0, 5}, {-3, -3, 5}}, // 90°
			{{-3, -3, -3}, {-3, 0, 5}, {-3, 5, 5}}, // 135°
			{{-3, -3, -3}, {-3, 0, -3}, {5, 5, 5}}, // 180°
			{{-3, -3, -3}, {5, 0, -3}, {5, 5, -3}}, // 225°
			{{5, -3, -3}, {5, 0, -3}, {5, -3, -3}}, // 270°
			{{5, 5, -3}, {5, 0, -3}, {-3, -3, -3}}, // 315°
		}

		for y := 1; y < bounds.Max.Y-1; y++ {
			for x := 1; x < bounds.Max.X-1; x++ {

				maxGradient := 0.0

				for _, kernel := range kernels {

					var sumR, sumG, sumB float64

					for ky := -1; ky <= 1; ky++ {
						for kx := -1; kx <= 1; kx++ {

							px := x + kx
							py := y + ky

							r, g, b, _ := img.Image.At(px, py).RGBA()
							factor := kernel[ky+1][kx+1]

							grayScale := (0.3*float64(r>>8) + 0.59*float64(g>>8) + 0.11*float64(b>>8))

							sumR += grayScale * factor
							sumG += grayScale * factor
							sumB += grayScale * factor
						}
					}

					gradient := math.Max(math.Abs(sumR), math.Max(math.Abs(sumG), math.Abs(sumB)))
					if gradient > maxGradient {
						maxGradient = gradient
					}
				}
				value := checkForLimit(maxGradient)

				highFreqImg.Set(x, y, color.RGBA{
					R: uint8(value),
					G: uint8(value),
					B: uint8(value),
					A: 255,
				})
			}
		}

		img.Image = highFreqImg
		img.Refresh()
	})

	return button
}

func NewPravitButton(img *canvas.Image, window fyne.Window) fyne.CanvasObject {

	button := widget.NewButton("Pravit", func() {

		if img.Image == nil {
			return
		}

		bounds := img.Image.Bounds()
		highFreqImg := image.NewRGBA(bounds)

		kernels := [2][3][3]float64{
			{{1, 0, -1}, {1, 0, -1}, {1, 0, -1}},
			{{-1, -1, -1}, {0, 0, 0}, {1, 1, 1}},
		}

		for y := 1; y < bounds.Max.Y-1; y++ {
			for x := 1; x < bounds.Max.X-1; x++ {

				maxGradient := 0.0

				for _, kernel := range kernels {

					var sumR, sumG, sumB float64

					for ky := -1; ky <= 1; ky++ {
						for kx := -1; kx <= 1; kx++ {

							px := x + kx
							py := y + ky

							r, g, b, _ := img.Image.At(px, py).RGBA()
							factor := kernel[ky+1][kx+1]

							grayScale := (0.3*float64(r>>8) + 0.59*float64(g>>8) + 0.11*float64(b>>8))

							sumR += grayScale * factor
							sumG += grayScale * factor
							sumB += grayScale * factor
						}
					}

					gradient := math.Max(math.Abs(sumR), math.Max(math.Abs(sumG), math.Abs(sumB)))
					if gradient > maxGradient {
						maxGradient = gradient
					}
				}
				value := checkForLimit(maxGradient)

				highFreqImg.Set(x, y, color.RGBA{
					R: uint8(value),
					G: uint8(value),
					B: uint8(value),
					A: 255,
				})
			}
		}

		img.Image = highFreqImg
		img.Refresh()
	})

	return button
}

func NewSobelButton(img *canvas.Image, window fyne.Window) fyne.CanvasObject {

	button := widget.NewButton("Sobel", func() {

		if img.Image == nil {
			return
		}

		bounds := img.Image.Bounds()
		highFreqImg := image.NewRGBA(bounds)

		kernels := [2][3][3]float64{
			{{-1, 0, 1}, {-2, 0, 2}, {-1, 0, 1}},
			{{1, 2, 1}, {0, 0, 0}, {-1, -2, -1}},
		}

		for y := 1; y < bounds.Max.Y-1; y++ {
			for x := 1; x < bounds.Max.X-1; x++ {

				gradient := 0.0

				for _, kernel := range kernels {

					var sumR, sumG, sumB float64

					for ky := -1; ky <= 1; ky++ {
						for kx := -1; kx <= 1; kx++ {

							px := x + kx
							py := y + ky

							r, g, b, _ := img.Image.At(px, py).RGBA()
							factor := kernel[ky+1][kx+1]

							grayScale := (0.3*float64(r>>8) + 0.59*float64(g>>8) + 0.11*float64(b>>8))

							sumR += grayScale * factor
							sumG += grayScale * factor
							sumB += grayScale * factor
						}
					}

					gradient += math.Pow(sumR, 2)

				}

				value := checkForLimit(math.Sqrt(gradient))

				highFreqImg.Set(x, y, color.RGBA{
					R: uint8(value),
					G: uint8(value),
					B: uint8(value),
					A: 255,
				})
			}
		}

		img.Image = highFreqImg
		img.Refresh()
	})

	return button
}

func NewRobertsButton(img *canvas.Image, window fyne.Window) fyne.CanvasObject {

	button := widget.NewButton("Roberts", func() {

		if img.Image == nil {
			return
		}

		bounds := img.Image.Bounds()
		changedImg := image.NewRGBA(bounds)

		for y := 0; y < bounds.Max.Y-1; y++ {
			for x := 0; x < bounds.Max.X-1; x++ {

				firstPixelR, firstPixelG, firstPixelB, _ := img.Image.At(x, y).RGBA()
				secondPixelR, secondPixelG, secondPixelB, _ := img.Image.At(x+1, y+1).RGBA()

				thirdPixelR, thirdPixelG, thirdPixelB, _ := img.Image.At(x+1, y).RGBA()
				fourthPixelR, fourthPixelG, fourthPixelB, _ := img.Image.At(x, y+1).RGBA()

				firstPixel := (0.3*float64(firstPixelR>>8) + 0.59*float64(firstPixelG>>8) + 0.11*float64(firstPixelB>>8))
				secondPixel := (0.3*float64(secondPixelR>>8) + 0.59*float64(secondPixelG>>8) + 0.11*float64(secondPixelB>>8))
				thirdPixel := (0.3*float64(thirdPixelR>>8) + 0.59*float64(thirdPixelG>>8) + 0.11*float64(thirdPixelB>>8))
				fourthPixel := (0.3*float64(fourthPixelR>>8) + 0.59*float64(fourthPixelG>>8) + 0.11*float64(fourthPixelB>>8))

				R_x_y := math.Sqrt(math.Pow(firstPixel-secondPixel, 2) + math.Pow(thirdPixel-fourthPixel, 2))

				value := checkForLimit(R_x_y)

				changedImg.Set(x, y, color.RGBA{
					R: uint8(value),
					G: uint8(value),
					B: uint8(value),
					A: 255,
				})
			}
		}

		img.Image = changedImg
		img.Refresh()
	})

	return button
}
