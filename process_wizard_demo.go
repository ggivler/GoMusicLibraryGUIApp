package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Wizard Dialog Example")
	w.Resize(fyne.NewSize(400, 300))

	// Step 1 content
	step1Content := container.NewVBox(
		widget.NewLabel("Welcome to the Wizard!"),
		widget.NewLabel("This is the first step."),
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
