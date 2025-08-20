package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"log"
	"os"
	"strings"
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

// suppressStderr temporarily redirects stderr to suppress Fyne's internal error logging
// This uses os.File operations which are cross-platform compatible
func suppressStderr(fn func()) {
	// Create a null device to redirect stderr to (platform-specific)
	nullDevice := "/dev/null"
	if os.Getenv("OS") == "Windows_NT" || os.Getenv("GOOS") == "windows" {
		nullDevice = "NUL"
	}
	
	nullFile, err := os.OpenFile(nullDevice, os.O_WRONLY, 0)
	if err != nil {
		// If we can't open null device, try temp file approach
		tempFile, err := os.CreateTemp("", "stderr_suppress")
		if err != nil {
			// If we can't create temp file, just run the function normally
			fn()
			return
		}
		defer os.Remove(tempFile.Name())
		nullFile = tempFile
	}
	defer nullFile.Close()
	
	// Save original stderr
	originalStderr := os.Stderr
	
	// Redirect stderr to null device
	os.Stderr = nullFile
	
	// Run the function with stderr suppressed
	fn()
	
	// Restore original stderr
	os.Stderr = originalStderr
}

func showFolderPicker(parentWindow fyne.Window, inputWidget *widget.Entry) {
	var folderDialog *dialog.FileDialog
	
	// Suppress stderr during dialog creation to eliminate Fyne error messages
	suppressStderr(func() {
		// Create folder dialog with callback
		folderDialog = dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				// Handle the error (e.g., show an error dialog)
				log.Printf("Folder selection error: %v", err)
				dialog.ShowError(err, parentWindow)
				return
			}
			if uri == nil {
				// User cancelled the dialog
				log.Println("Folder selection cancelled")
				return
			}
			// Process the selected folder URI
			selectedPath := uri.Path()
			log.Println(selectedPath)
			fmt.Printf("Selected folder: %s\n", selectedPath)
			
			// Set the selected folder path in the input widget
			if inputWidget != nil {
				inputWidget.SetText(selectedPath)
			}
			
			// TODO: Save selectedPath to the config yaml file
			// Need a function here.

		}, parentWindow)
		
		// Configure dialog settings
		folderDialog.Resize(fyne.NewSize(800, 600))
		
		// Set a safe starting location - this helps minimize the favorites issue
		if startLocation := getSafeStartLocation(); startLocation != nil {
			folderDialog.SetLocation(startLocation)
		} else {
			// Fallback to a known safe location
			if userProfile := os.Getenv("USERPROFILE"); userProfile != "" {
				documentsPath := userProfile + "\\Documents"
				if documentsURI := storage.NewFileURI(documentsPath); documentsURI != nil {
					if listableURI, ok := documentsURI.(fyne.ListableURI); ok {
						folderDialog.SetLocation(listableURI)
					}
				}
			}
		}
	})

	// Show the dialog with stderr restored
	folderDialog.Show()
}

// CustomLogger implements a filtered logger that suppresses certain error messages
type CustomLogger struct {
	origLogger *log.Logger
	filterText []string
}

// NewCustomLogger creates a logger that suppresses specific error messages
func NewCustomLogger(filterMessages ...string) *CustomLogger {
	return &CustomLogger{
		origLogger: log.Default(),
		filterText: filterMessages,
	}
}

// Write implements io.Writer interface and filters out unwanted messages
func (cl *CustomLogger) Write(p []byte) (n int, err error) {
	msg := string(p)
	
	// Check if message contains any of the filter texts
	for _, filter := range cl.filterText {
		if strings.Contains(msg, filter) {
			// Message should be filtered - don't log it
			return len(p), nil
		}
	}
	
	// Message passed the filter, log it using the original logger
	return os.Stderr.Write(p)
}

// setupErrorFiltering configures the logger to suppress certain Fyne errors
func setupErrorFiltering() {
	// Create a custom logger that filters out the URI listable error
	customLogger := NewCustomLogger(
		"Fyne error:  Getting favorite locations",
		"Cause: uri is not listable",
		"At: C:/Users/ggivl/go/pkg/mod/fyne.io/fyne/v2", // Filter the stack trace line
		"dialog/file.go:", // Filter any file.go error location
	)
	
	// Set it as the output for the default logger
	log.SetOutput(customLogger)
}

func main() {
	// Setup error filtering to handle Fyne errors gracefully
	setupErrorFiltering()
	
	// KNOWN ISSUE: You may see a "Getting favorite locations - Cause: uri is not listable" error
	// when clicking the folder selection button. This is a harmless Fyne framework issue
	// on Windows systems. The error occurs because Fyne tries to access Windows favorite
	// locations that may not be accessible. This does NOT affect functionality - the
	// file dialog will still work correctly and you can select folders normally.
	// 
	// We've added a custom error filter to suppress these messages.
	
	// Create app with unique ID to fix "Preferences API requires a unique ID" error
	a := app.NewWithID("com.gomusiclibrary.wizardapp")
	w := a.NewWindow("Music Library Wizard")
	w.Resize(fyne.NewSize(600, 400))

	// Step 1 content
	input := widget.NewEntry() // Create input widget
	input.SetPlaceHolder("Selected folder will appear here")
	step1Content := container.NewVBox(
		widget.NewLabel("Welcome to the Wizard!"),
		widget.NewLabel("Select the Music Library Folder to process"),
		widget.NewButton("Select the Music Library Folder", func() {
			fmt.Println("Enter Select the Music Library Folder to process")
			showFolderPicker(w, input)
		}),
		input, // Add the input widget
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
