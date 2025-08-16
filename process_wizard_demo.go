package main

import (
	"fmt"
	"log"
	"os"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

// getSafeStartLocation returns a safe starting location for file dialogs
func getSafeStartLocation() fyne.ListableURI {
	// Helper function to safely create and test a listable URI
	createListableURI := func(path string) fyne.ListableURI {
		if _, err := os.Stat(path); err != nil {
			return nil
		}
		uri := storage.NewFileURI(path)
		if uri == nil {
			return nil
		}
		if listableURI, ok := uri.(fyne.ListableURI); ok {
			// Test if we can actually list the directory
			if _, err := listableURI.List(); err == nil {
				return listableURI
			}
		}
		return nil
	}

	// Try C:\ root directory first (always accessible on Windows)
	if rootURI := createListableURI("C:\\"); rootURI != nil {
		return rootURI
	}

	// Try Documents directory
	if userProfile := os.Getenv("USERPROFILE"); userProfile != "" {
		documentsDir := userProfile + "\\Documents"
		if docsURI := createListableURI(documentsDir); docsURI != nil {
			return docsURI
		}
	}

	// Fallback to home directory
	if homeDir, err := os.UserHomeDir(); err == nil {
		if homeURI := createListableURI(homeDir); homeURI != nil {
			return homeURI
		}
	}

	// Final fallback to current directory
	if currentDir, err := os.Getwd(); err == nil {
		if currentURI := createListableURI(currentDir); currentURI != nil {
			return currentURI
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
		// Save uri.Path() to the config yaml file
		// Need a function here.

	}, parentWindow)

	// Set a safe starting location only if we can get one
	if startLocation := getSafeStartLocation(); startLocation != nil {
		folderDialog.SetLocation(startLocation)
	}

	// Configure dialog to avoid problematic favorite locations
	// This helps prevent the "uri is not listable" error
	folderDialog.Resize(fyne.NewSize(800, 600))

	folderDialog.Show()
}

func main() {
	// Note: The "Getting favorite locations - Cause: uri is not listable" warning
	// is a known Fyne issue on Windows where the file dialog tries to access
	// system favorite locations that may not be accessible. This doesn't affect
	// functionality and can be safely ignored.

	// Create app with unique ID to fix "Preferences API requires a unique ID" error
	a := app.NewWithID("com.gomusiclibrary.wizardapp")
	w := a.NewWindow("Wizard Dialog Example")
	w.Resize(fyne.NewSize(600, 400))

	// Step 1 content
	step1Content := container.NewVBox(
		widget.NewLabel("Welcome to the Wizard!"),
		widget.NewLabel("Select the Music Library Folder to process"),
		widget.NewButton("Select the Music Library Folder", func() {
			fmt.Println("Enter Select the Music Library Folder to process")
			showFolderPicker(w)
		}),
		widget.NewEntry(), // Example input
	)

	// Step 2 content
	step2Content := container.NewVBox(
		widget.NewLabel("This is the second step."),
		widget.NewCheck("Option 1", func(b bool) {}),
		widget.NewCheck("Option 2", func(b bool) {}),
	)

	// Step 3 content
	step3Content := container.NewVBox(
		widget.NewLabel("You've reached the final step."),
		widget.NewLabel("Click 'Finish' to complete."),
	)

	// Keep track of current step
	currentStep := 0
	steps := []fyne.CanvasObject{step1Content, step2Content, step3Content}

	// Create a container to hold the current step's content
	contentContainer := container.NewStack()
	contentContainer.Add(steps[currentStep])

	// Navigation buttons
	backButton := widget.NewButton("Back", nil)
	nextButton := widget.NewButton("Next", nil)
	finishButton := widget.NewButton("Finish", nil)
	finishButton.Hide() // Hidden initially

	// Update button visibility and content based on step
	updateUI := func() {
		backButton.Enable()
		nextButton.Enable()
		finishButton.Hide()

		if currentStep == 0 {
			backButton.Disable()
		}
		if currentStep == len(steps)-1 {
			nextButton.Hide()
			finishButton.Show()
		} else {
			nextButton.Show()
		}

		contentContainer.Objects = []fyne.CanvasObject{steps[currentStep]}
		contentContainer.Refresh()
	}

	backButton.OnTapped = func() {
		if currentStep > 0 {
			currentStep--
			updateUI()
		}
	}

	nextButton.OnTapped = func() {
		if currentStep < len(steps)-1 {
			currentStep++
			updateUI()
		}
	}

	finishButton.OnTapped = func() {
		// Handle wizard completion logic here
		dialog.ShowInformation("Wizard Complete", "You have finished the wizard!", w)
		// Optionally close the dialog or perform other actions
	}

	// Layout for the dialog
	dialogContent := container.NewBorder(
		nil,
		container.NewHBox(backButton, layout.NewSpacer(), nextButton, finishButton),
		nil,
		nil,
		contentContainer,
	)

	// Show the custom dialog
	customDialog := dialog.NewCustom("Wizard", "Cancel", dialogContent, w)
	customDialog.SetOnClosed(func() {
		// Handle dialog closure (e.g., if user clicks "Cancel")
	})

	updateUI() // Initial UI setup
	customDialog.Show()

	w.ShowAndRun()
}
