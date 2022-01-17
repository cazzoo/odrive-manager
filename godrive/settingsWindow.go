package godrive

import (
	"fmt"

	"github.com/andlabs/ui"
)

type settingsWindow struct {
	clientHandler  IOdriveClientHandler
	settingsWindow func()
	settings       odriveSettings
}

type FormEntry struct {
	label     string
	component ui.Control
}

func SettingWindow(client IOdriveClientHandler) func() {
	w := &settingsWindow{}
	w.clientHandler = client
	w.settings = LoadSettings()
	w.settingsWindow = w.new
	return w.settingsWindow
}

func (w settingsWindow) new() {
	mainWindow := ui.NewWindow("Godrive settings window", 640, 480, true)

	mainWindow.OnClosing(func(*ui.Window) bool {
		mainWindow.Destroy()
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainWindow.Destroy()
		return true
	})

	vbContainer := ui.NewVerticalBox()
	vbContainer.SetPadded(true)

	inputGroup := ui.NewGroup("Input")
	inputGroup.SetMargined(true)

	vbInput := ui.NewVerticalBox()
	vbInput.SetPadded(true)

	inputForm := ui.NewForm()
	inputForm.SetPadded(true)

	message := ui.NewEntry()
	message.SetText("Default value")
	inputForm.Append("Select Odrive mount path", message, false)

	vbInput.Append(inputForm, false)
	entries := []FormEntry{}

	comboSplitSize := w.generateDropdown(splitSizeElements(), w.settings.Xlthreshold.String())
	entries = append(entries, FormEntry{"Split files that are larger than this threshold into chunks (default Large).", comboSplitSize})
	comboTrashCleanFreq := w.generateDropdown(trashCleanFrequencyElements(), w.settings.TrashClean.String())
	entries = append(entries, FormEntry{"Set rule for automatically emptying the odrive trash (default Never).", comboTrashCleanFreq})
	comboAutoUnsync := w.generateDropdown(unsyncPeriodElements(), w.settings.Autounsyncthreshold.String())
	entries = append(entries, FormEntry{"Set rule for automatically unsyncing files that have not been modified with a certain amount of time (default Never).", comboAutoUnsync})
	comboPlaceholderTreshold := w.generateDropdown(placeholderSizeElements(), w.settings.Placeholderthreshold.String())
	entries = append(entries, FormEntry{"Set rule for automatically downloading files under a specified size when syncing/expanding a folder (default Never).", comboPlaceholderTreshold})
	vbInput.Append(w.generateForm(entries), false)

	saveButton := ui.NewButton("Save")
	cancelButton := ui.NewButton("Cancel")

	vbInput.Append(saveButton, false)
	vbInput.Append(cancelButton, false)

	inputGroup.SetChild(vbInput)

	odriveGroup := ui.NewGroup("odrive command result")
	odriveGroup.SetMargined(true)

	vbOdrive := ui.NewVerticalBox()
	vbOdrive.SetPadded(true)

	odriveLabel := ui.NewLabel("")
	vbOdrive.Append(odriveLabel, false)

	odriveGroup.SetChild(vbOdrive)

	vbContainer.Append(inputGroup, false)
	vbContainer.Append(odriveGroup, false)

	mainWindow.SetChild(vbContainer)

	saveButton.OnClicked(func(*ui.Button) {
		settings := odriveSettings{
			MountPaths:           []string{},
			Autounsyncthreshold:  unsyncPeriod(comboAutoUnsync.Selected()),
			Placeholderthreshold: placeholderSize(comboPlaceholderTreshold.Selected()),
			TrashClean:           trashCleanFrequency(comboTrashCleanFreq.Selected()),
			Xlthreshold:          splitSize(comboSplitSize.Selected()),
		}

		SaveSettings(settings)
		fmt.Printf("Updating settings.\n")
		mainWindow.Destroy()
	})

	cancelButton.OnClicked(func(*ui.Button) {
		fmt.Printf("Canceling changes.\n")
		mainWindow.Destroy()
	})

	mainWindow.Show()
}

func (w settingsWindow) generateForm(elements []FormEntry) *ui.Form {
	inputForm := ui.NewForm()
	inputForm.SetPadded(true)

	for _, element := range elements {
		inputForm.Append(element.label, element.component, false)
	}
	return inputForm
}

func (w settingsWindow) generateDropdown(values []string, defaultValue string) *ui.Combobox {
	dropdown := ui.NewCombobox()
	defaultSelection := 0
	for i, value := range values {
		dropdown.Append(value)
		if value == defaultValue {
			defaultSelection = i
		}
	}
	dropdown.SetSelected(defaultSelection)
	return dropdown
}
