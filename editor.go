package main

import (
	"io/ioutil"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type config_editor struct {
	EditWidget    *widget.Entry
	PreviewWidget *widget.RichText
	CurrenFile    fyne.URI // Path
	SaveMenuItem  *fyne.MenuItem
}

var editor config_editor

func main() {
	a := app.New()

	win := a.NewWindow("Mardown Editor")

	edit, preview := editor.Init_Editor()

	win.SetContent(container.NewHSplit(preview, edit))
	editor.createItemMenu(win)

	win.Resize(fyne.Size{Width: 1000, Height: 1000})
	win.CenterOnScreen()
	win.ShowAndRun()
}

func (c *config_editor) Init_Editor() (*widget.Entry, *widget.RichText) {
	edit := widget.NewMultiLineEntry()
	preview := widget.NewRichTextFromMarkdown("")

	c.EditWidget = edit
	c.PreviewWidget = preview

	edit.OnChanged = preview.ParseMarkdown

	return edit, preview
}

func (c *config_editor) openMenuFunc(win fyne.Window) func() {
	return func() {
		openDialog := dialog.NewFileOpen(func(read fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
			}

			if read == nil {
				return
			}

			data, err := ioutil.ReadAll(read)

			defer read.Close()

			if err != nil {
				dialog.ShowError(err, win)
			}

			c.EditWidget.SetText(string(data))

			c.CurrenFile = read.URI()
			c.SaveMenuItem.Disabled = false

		}, win)
		openDialog.Show()
	}
}

func (c *config_editor) saveFileMenu(win fyne.Window) func() {
	return func() {
		if c.CurrenFile != nil {
			writer, err := storage.Writer(c.CurrenFile)

			defer writer.Close()

			if err != nil {
				dialog.ShowError(err, win)
			}

			writer.Write([]byte(c.EditWidget.Text))

		}
	}
}

func (c *config_editor) saveAsFileMenu(win fyne.Window) func() {
	return func() {
		saveDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
			}

			if writer == nil {
				return
			}

			if len(c.EditWidget.Text) == 0 {
				return
			}

			// only save .md file
			if !strings.HasSuffix(strings.ToLower(writer.URI().String()), ".md") {
				dialog.ShowInformation("Error", "Save Only md File", win)
			}

			defer writer.Close()

			writer.Write([]byte(c.EditWidget.Text))

			win.SetTitle(win.Title() + " - " + writer.URI().Name())

			c.CurrenFile = writer.URI()
			c.SaveMenuItem.Disabled = false

		}, win)

		saveDialog.Show()
	}
}

func (c *config_editor) createItemMenu(win fyne.Window) {
	openFileMenu := fyne.NewMenuItem("Open File", c.openMenuFunc(win))
	saveFileMenu := fyne.NewMenuItem("Save File", c.saveFileMenu(win))

	saveAsFileMenu := fyne.NewMenuItem("Save as File", c.saveAsFileMenu(win))

	c.SaveMenuItem = saveFileMenu
	c.SaveMenuItem.Disabled = true

	fileMenu := fyne.NewMenu("Open Menu", openFileMenu, saveFileMenu, saveAsFileMenu)
	menu := fyne.NewMainMenu(fileMenu)

	win.SetMainMenu(menu)
}
