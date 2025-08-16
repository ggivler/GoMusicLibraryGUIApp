package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func showFolderPicker(parentWindow fyne.Window) {
	dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
		if err != nil {
			// Handle the error (e.g., show an error dialog)
			dialog.ShowError(err, parentWindow)
			return
		}
		if uri == nil {
			// User cancelled the dialog
			return
		}
		// Process the selected folder URI (e.g., display its path)
		fmt.Printf("Selected folder: %s\n", uri.Path())
	}, parentWindow)
}

func main() {
	myApp := app.New()
	w := myApp.NewWindow("GoMusicLibrary")

	newItem := fyne.NewMenuItem("New", func() {
		fmt.Println("New file selected")
		// Add your logic for creating a new file here
		showFolderPicker(w)

	})
	openItem := fyne.NewMenuItem("Open", func() {
		fmt.Println("Open file selected")
		// Add your logic for opening a file here
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				// Handle the opened file
				defer reader.Close()
				fmt.Println("Opened file:", reader.URI().Name())
			}
		}, w)
	})
	saveItem := fyne.NewMenuItem("Save", func() {
		fmt.Println("Save file selected")
		// Add your logic for saving a file here
		dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err == nil && writer != nil {
				// Handle the file to be saved
				defer writer.Close()
				fmt.Println("Saving to:", writer.URI().Name())
			}

		}, w)

	})
	exitItem := fyne.NewMenuItem("Exit", func() {
		fmt.Println("Exit selected")
		w.Close()
	})

	// Create the File menu
	fileMenu := fyne.NewMenu("File", newItem, openItem, saveItem, fyne.NewMenuItemSeparator(), exitItem)

	// Create the main menu
	mainMenu := fyne.NewMainMenu(fileMenu)

	// Set the main menu on the window
	w.SetMainMenu(mainMenu)

	w.SetContent(widget.NewLabel("Click File for Menu"))
	w.Resize(fyne.NewSize(300, 200))
	w.ShowAndRun()
}
