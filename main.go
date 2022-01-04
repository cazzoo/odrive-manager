package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"time"

	"cazzoo.me/godrive/godrive"
	"cazzoo.me/godrive/process"
	"github.com/andlabs/ui"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/skratchdot/open-golang/open"
)

var odriveAgentPath string
var odriveClientPath string
var odriveAgentHandler godrive.IOdriveAgentHandler
var odriveClientHandler godrive.IOdriveClientHandler
var menuItems map[string]*systray.MenuItem
var schedulerChan chan struct{}

func main() {
	onExit := func() {
		now := time.Now()
		ioutil.WriteFile(fmt.Sprintf(`on_exit_%d.txt`, now.UnixNano()), []byte(now.String()), 0644)
	}

	pathAgent, err := exec.LookPath("odriveagent")
	if err != nil {
		log.Fatal("Counldn't find [odriveagent] executable in PATH environment variable.\nPlease add it in order to continue.")
	} else {
		odriveAgentPath = pathAgent
	}

	pathClient, err := exec.LookPath("odrive")
	if err != nil {
		log.Fatal("Counldn't find [odrive] executable in PATH environment variable.\nPlease add it in order to continue.")
	} else {
		odriveClientPath = pathClient
	}

	odriveAgentHandler = godrive.OdriveAgentHandler(odriveAgentPath)
	odriveClientHandler = godrive.OdriveClientHandler(odriveClientPath)

	// Restarting the agent if already started
	if processes, err := process.FindProcess("odriveagent"); err == nil {
		fmt.Print("Process seemed to be started before, trying to restart it.")
		if err := process.KillProcesses(processes); err != nil {
			fmt.Printf("Error stoping agent.")
		} else {
			if err := odriveAgentHandler.Start(); err != nil {
				fmt.Printf("Unable to start [odriveagent]: %s", err)
			}
		}
	}

	menuItems = make(map[string]*systray.MenuItem)

	systray.Run(generateMenu, onExit)
}

func startScheduledChecks() {
	tickIndex := 0
	schedulerChan = make(chan struct{})
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				tickIndex++
				fmt.Printf("odrive agent status started: %t\n", odriveAgentHandler.HealthCheck())
				fmt.Printf("Ticking for the %d time\n", tickIndex)
			case <-schedulerChan:
				ticker.Stop()
				return
			}
		}
	}()
}

func generateMenu() {
	startScheduledChecks()
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("Odrive manager")
	systray.SetTooltip("Odrive manager")
	startOdrive := systray.AddMenuItem("Start odrive agent", "Start odrive agent")
	stopOdrive := systray.AddMenuItem("Stop odrive agent", "Stop odrive agent")
	stopChan := systray.AddMenuItem("Stop chan", "Stop chan")
	displayWindow := systray.AddMenuItem("Show window", "Show window")
	stopOdrive.Hide()
	if _, err := process.FindProcess("odriveagent"); err == nil {
		startOdrive.Hide()
		stopOdrive.Show()
	} else {
		startOdrive.Show()
		stopOdrive.Hide()
	}
	menuItems["startOdrive"] = startOdrive
	menuItems["stopOdrive"] = stopOdrive
	go func() {
		for {
			select {
			case <-startOdrive.ClickedCh:
				startAgent(odriveAgentHandler)
			case <-stopOdrive.ClickedCh:
				stopAgent(odriveAgentHandler)
			case <-displayWindow.ClickedCh:
				ui.Main(setupUI)
			case <-stopChan.ClickedCh:
				close(schedulerChan)
			}
		}
	}()
	onReady()
}

func startAgent(odriveAgentHandler godrive.IOdriveAgentHandler) {
	if err := odriveAgentHandler.Start(); err != nil {
		fmt.Printf("Unable to start [odriveagent]: %s", err)
	} else {
		menuItems["startOdrive"].Hide()
		menuItems["stopOdrive"].Show()
	}
}

func stopAgent(odriveAgentHandler godrive.IOdriveAgentHandler) {
	if err := odriveAgentHandler.Stop(); err != nil {
		fmt.Printf("Unable to stop [odriveagent]: %s", err)
	} else {
		menuItems["startOdrive"].Show()
		menuItems["stopOdrive"].Hide()
	}
}

func onReady() {
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuitOrig.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

	// We can manipulate the systray in other goroutines
	go func() {
		mChange := systray.AddMenuItem("Change Me", "Change Me")
		mChecked := systray.AddMenuItemCheckbox("Unchecked", "Check Me", true)
		mEnabled := systray.AddMenuItem("Enabled", "Enabled")
		// Sets the icon of a menu item. Only available on Mac.
		mEnabled.SetTemplateIcon(icon.Data, icon.Data)

		systray.AddMenuItem("Ignored", "Ignored")

		subMenuTop := systray.AddMenuItem("SubMenuTop", "SubMenu Test (top)")
		subMenuMiddle := subMenuTop.AddSubMenuItem("SubMenuMiddle", "SubMenu Test (middle)")
		subMenuBottom := subMenuMiddle.AddSubMenuItemCheckbox("SubMenuBottom - Toggle Panic!", "SubMenu Test (bottom) - Hide/Show Panic!", false)
		subMenuBottom2 := subMenuMiddle.AddSubMenuItem("SubMenuBottom - Panic!", "SubMenu Test (bottom)")

		mUrl := systray.AddMenuItem("Open UI", "my home")
		mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

		// Sets the icon of a menu item. Only available on Mac.
		mQuit.SetIcon(icon.Data)

		systray.AddSeparator()
		mToggle := systray.AddMenuItem("Toggle", "Toggle the Quit button")
		shown := true
		toggle := func() {
			if shown {
				subMenuBottom.Check()
				subMenuBottom2.Hide()
				mQuitOrig.Hide()
				mEnabled.Hide()
				shown = false
			} else {
				subMenuBottom.Uncheck()
				subMenuBottom2.Show()
				mQuitOrig.Show()
				mEnabled.Show()
				shown = true
			}
		}

		for {
			select {
			case <-mChange.ClickedCh:
				mChange.SetTitle("I've Changed")
			case <-mChecked.ClickedCh:
				if mChecked.Checked() {
					mChecked.Uncheck()
					mChecked.SetTitle("Unchecked")
				} else {
					mChecked.Check()
					mChecked.SetTitle("Checked")
				}
			case <-mEnabled.ClickedCh:
				mEnabled.SetTitle("Disabled")
				mEnabled.Disable()
			case <-mUrl.ClickedCh:
				open.Run("https://www.getlantern.org")
			case <-subMenuBottom2.ClickedCh:
				panic("panic button pressed")
			case <-subMenuBottom.ClickedCh:
				toggle()
			case <-mToggle.ClickedCh:
				toggle()
			case <-mQuit.ClickedCh:
				systray.Quit()
				fmt.Println("Quit2 now...")
				return
			}
		}
	}()
}

func setupUI() {
	mainWindow := ui.NewWindow("libui Updating UI", 640, 480, true)
	mainWindow.OnClosing(func(*ui.Window) bool {
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
	message.SetText("Hello World")
	inputForm.Append("What message do you want to show?", message, false)

	showMessageButton := ui.NewButton("Show message")
	clearMessageButton := ui.NewButton("Clear message")
	showOdriveStatusButton := ui.NewButton("Get Odrive status")

	vbInput.Append(inputForm, false)
	vbInput.Append(showMessageButton, false)
	vbInput.Append(clearMessageButton, false)
	vbInput.Append(showOdriveStatusButton, false)

	inputGroup.SetChild(vbInput)

	messageGroup := ui.NewGroup("Message")
	messageGroup.SetMargined(true)

	vbMessage := ui.NewVerticalBox()
	vbMessage.SetPadded(true)

	messageLabel := ui.NewLabel("")

	vbMessage.Append(messageLabel, false)

	messageGroup.SetChild(vbMessage)

	countGroup := ui.NewGroup("Counter")
	countGroup.SetMargined(true)

	vbCounter := ui.NewVerticalBox()
	vbCounter.SetPadded(true)

	countGroup.SetChild(vbCounter)

	odriveGroup := ui.NewGroup("odrive command result")
	odriveGroup.SetMargined(true)

	vbOdrive := ui.NewVerticalBox()
	vbOdrive.SetPadded(true)

	odriveLabel := ui.NewLabel("")
	vbOdrive.Append(odriveLabel, false)

	odriveGroup.SetChild(vbOdrive)

	vbContainer.Append(inputGroup, false)
	vbContainer.Append(messageGroup, false)
	vbContainer.Append(countGroup, false)
	vbContainer.Append(odriveGroup, false)

	mainWindow.SetChild(vbContainer)

	showMessageButton.OnClicked(func(*ui.Button) {
		// Update the UI directly as it is called from the main thread
		messageLabel.SetText(message.Text())
	})

	clearMessageButton.OnClicked(func(*ui.Button) {
		// Update the UI directly as it is called from the main thread
		messageLabel.SetText("")
	})

	showOdriveStatusButton.OnClicked(func(*ui.Button) {
		err, output := odriveClientHandler.Call(godrive.Status)
		if err == nil {
			messageLabel.SetText(string(output))
		} else {
			fmt.Printf("ERROR: %s", err)
		}
	})

	mainWindow.Show()
}
