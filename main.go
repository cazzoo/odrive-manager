package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"time"

	"cazzoo.me/godrive/godrive"
	"cazzoo.me/godrive/process"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/skratchdot/open-golang/open"
)

var odriveAgentPath string
var odriveHandler godrive.IOdriveAgentHandler
var menuItems map[string]*systray.MenuItem

func main() {
	onExit := func() {
		now := time.Now()
		ioutil.WriteFile(fmt.Sprintf(`on_exit_%d.txt`, now.UnixNano()), []byte(now.String()), 0644)
	}

	pathAgent, err := exec.LookPath("odriveagent")
	if err != nil {
		log.Fatal("Counldn't find [odriveagent] executable in path.")
	} else {
		odriveAgentPath = pathAgent
	}

	odriveHandler = godrive.OdriveAgentHandler(odriveAgentPath)

	// Restarting the agent if already started
	if processes, err := process.FindProcess("odriveagent"); err == nil {
		fmt.Print("Process seemed to be started before, trying to restart it.")
		if err := process.KillProcesses(processes); err != nil {
			fmt.Printf("Error stoping agent.")
		} else {
			if err := odriveHandler.Start(); err != nil {
				fmt.Printf("Unable to start [odriveagent]: %s", err)
			}
		}
	}

	menuItems = make(map[string]*systray.MenuItem)

	systray.Run(generateMenu, onExit)
}

func generateMenu() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("Awesome App")
	systray.SetTooltip("Lantern")
	startOdrive := systray.AddMenuItem("Start odrive agent", "Start odrive agent")
	stopOdrive := systray.AddMenuItem("Stop odrive agent", "Stop odrive agent")
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
				startAgent(odriveHandler)
			case <-stopOdrive.ClickedCh:
				stopAgent(odriveHandler)
			}
		}
	}()
	onReady()
}

func startAgent(odriveHandler godrive.IOdriveAgentHandler) {
	if err := odriveHandler.Start(); err != nil {
		fmt.Printf("Unable to start [odriveagent]: %s", err)
	} else {
		menuItems["startOdrive"].Hide()
		menuItems["stopOdrive"].Show()
	}
}

func stopAgent(odriveHandler godrive.IOdriveAgentHandler) {
	if err := odriveHandler.Stop(); err != nil {
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
