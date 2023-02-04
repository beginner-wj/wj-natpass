package common

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type SubmitCallback func(bool)

func Toast(title string, tip string) {
	w := fyne.CurrentApp().NewWindow(title)
	tipText := widget.NewLabel(tip)
	subBtn := widget.NewButton("submit", func() {
		w.Close()
	})
	w.SetContent(container.NewVBox(tipText, subBtn))
	w.Resize(fyne.NewSize(100, 100))
	w.CenterOnScreen()
	w.Show()
}

func ToastSucc(tip string) {
	Toast("succ", tip)
}

func ToastSubmitFunc(title string, tip string, callback SubmitCallback) {
	w := fyne.CurrentApp().NewWindow(title)
	tipText := widget.NewLabel(tip)
	subBtn := widget.NewButton("submit", func() {
		w.Close()
		callback(true)
	})
	cancleBtn := widget.NewButton("cancel", func() {
		w.Close()
		callback(false)
	})
	v1 := container.NewBorder(layout.NewSpacer(), layout.NewSpacer(), cancleBtn, subBtn)
	w.SetContent(container.NewVBox(tipText, v1))
	w.Resize(fyne.NewSize(100, 100))
	w.CenterOnScreen()
	w.Show()
}

func ToastSelect(title string, tip string) {
	w := fyne.CurrentApp().NewWindow(title)
	tipText := widget.NewLabel(tip)
	subBtn := widget.NewButton("submit", func() {
		w.Close()
	})
	cancleBtn := widget.NewButton("cancel", func() {
		w.Close()
	})
	v1 := container.NewBorder(layout.NewSpacer(), layout.NewSpacer(), cancleBtn, subBtn)
	w.SetContent(container.NewVBox(tipText, v1))
	w.Resize(fyne.NewSize(100, 100))
	w.CenterOnScreen()
	w.Show()
}

func ToastError(tip string) {
	w := fyne.CurrentApp().NewWindow("error tips")

	tipText := widget.NewLabel(tip)

	subBtn := widget.NewButton("submit", func() {
		w.Close()
	})

	w.SetContent(container.NewVBox(tipText, subBtn))
	w.Resize(fyne.NewSize(100, 100))
	w.CenterOnScreen()
	w.Show()
}
