package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

// customLogWriter filters out specific Fyne error messages
type customLogWriter struct {
	original io.Writer
}

func (w *customLogWriter) Write(p []byte) (n int, err error) {
	message := string(p)
	// Filter out the specific "Getting favorite locations" error and related messages
	if strings.Contains(message, "Getting favorite locations") ||
		strings.Contains(message, "uri is not listable") ||
		strings.Contains(message, "dialog/file.go:367") {
		// Silently discard these messages
		return len(p), nil
	}
	// Pass through all other messages
	return w.original.Write(p)
}

// getSafeStartLocation returns a safe starting location for file dialogs
func getSafeStartLocation() fyne.ListableURI {
	// Try Documents directory first (common and usually accessible)
	documentsDir := os.Getenv("USERPROFILE") + "\\Documents"
	if _, err := os.Stat(documentsDir); err == nil {
		if uri := storage.NewFileURI(documentsDir); uri != nil {
			if listableURI, ok := uri.(fyne.ListableURI); ok {
				return listableURI
			}
		}
	}
	
	// Fallback to home directory
	homeDir, err := os.UserHomeDir()
	if err == nil {
		if uri := storage.NewFileURI(homeDir); uri != nil {
			if listableURI, ok := uri.(fyne.ListableURI); ok {
				return listableURI
			}
		}
	}
	
	// Final fallback to current directory
	currentDir, _ := os.Getwd()
	if uri := storage.NewFileURI(currentDir); uri != nil {
		if listableURI, ok := uri.(fyne.ListableURI); ok {
			return listableURI
		}
	}
	
	return nil
}

func showFolderPicker(parentWindow fyne.Window) {
	folderDialog := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
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
	
	// Set a safe starting location
	if startLocation := getSafeStartLocation(); startLocation != nil {
		folderDialog.SetLocation(startLocation)
	}
	
	folderDialog.Show()
}

func main() {
	// Setup custom log writer to filter out Fyne favorite location errors
	customWriter := &customLogWriter{
		original: os.Stderr,
	}
	log.SetOutput(customWriter)
	
	// Set environment variable to potentially reduce favorite location issues
	os.Setenv("FYNE_THEME", "light")
	
	myApp := app.NewWithID("GoMusicLibraryGUI")
	w := myApp.NewWindow("GoMusicLibrary")

	newItem := fyne.NewMenuItem("New", func() {
		fmt.Println("New file selected")
		// Add your logic for creating a new file here
		showFolderPicker(w)

	})
	openItem := fyne.NewMenuItem("Open", func() {
		fmt.Println("Open file selected")
		// Add your logic for opening a file here
		fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				// Handle the opened file
				defer reader.Close()
				fmt.Println("Opened file:", reader.URI().Name())
			}
		}, w)
		
		// Set a safe starting location
		if startLocation := getSafeStartLocation(); startLocation != nil {
			fileDialog.SetLocation(startLocation)
		}
		
		fileDialog.Show()
	})
	saveItem := fyne.NewMenuItem("Save", func() {
		fmt.Println("Save file selected")
		// Add your logic for saving a file here
		saveDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err == nil && writer != nil {
				// Handle the file to be saved
				defer writer.Close()
				fmt.Println("Saving to:", writer.URI().Name())
			}
		}, w)
		
		// Set a safe starting location
		if startLocation := getSafeStartLocation(); startLocation != nil {
			saveDialog.SetLocation(startLocation)
		}
		
		saveDialog.Show()
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
	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}
