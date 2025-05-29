package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/icons"
)

var (
	windowShowing bool
	currentWindow *application.WebviewWindow
	app           *application.App
)

func createWindow(app *application.App) *application.WebviewWindow {
	// Log the time taken to create the window
	window := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Name:            "BeoEcho",
		AlwaysOnTop:     true,
		Hidden:          false,
		Frameless:       false,
		DevToolsEnabled: false,
		Windows: application.WindowsWindow{
			HiddenOnTaskbar: true,
		},
	})

	window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
		fmt.Println("Window closing event triggered")
		windowShowing = false
		currentWindow = nil
	})

	window.OnWindowEvent(events.Common.WindowMinimise, func(e *application.WindowEvent) {
		fmt.Println("Window minimise event triggered")
		windowShowing = true
	})
	window.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
		fmt.Println("Window lost focus event triggered")
		windowShowing = false
		window.Minimise()
	})

	window.Show()
	return window
}

// checkAndShowWindow checks if a window exists and is showing, creates one if needed
func checkAndShowWindow() {
	if currentWindow == nil || !windowShowing {
		fmt.Println("Creating new window")
		currentWindow = createWindow(app)
		windowShowing = true
		app.Show()
	} else {
		fmt.Println("Showing existing window")
		currentWindow.Show()
		app.Show()
	}
}

// hideWindow hides the current window if it exists
func hideWindow() {
	if currentWindow != nil && windowShowing {
		fmt.Println("Hiding window")
		currentWindow.Hide()
		app.Hide()
		windowShowing = false
	}
}

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed all:frontend
var assets embed.FS

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {
	// Setup logging to file for desktop app debugging
	setupLogging()

	log.Println("🚀 Starting BeoEcho Desktop Application...")
	log.Printf("Current working directory: %s", getCurrentWorkingDir())
	log.Printf("Executable path: %s", getExecutablePath())
	log.Printf("Environment PATH: %s", os.Getenv("PATH"))

	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	backendService := NewBackendService()
	app = application.New(application.Options{
		Name:        "BeoEcho",
		Description: "Desktop API Mocking Service",
		Services: []application.Service{
			application.NewService(backendService, application.ServiceOptions{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Windows: application.WindowsOptions{
			DisableQuitOnLastWindowClosed: true,
		},
		Mac: application.MacOptions{
			ActivationPolicy: application.ActivationPolicyAccessory,
		},
	})

	// Create a goroutine that emits an event containing the current time every second.
	// The frontend can listen to this event and update the UI accordingly.
	go func() {
		for {
			now := time.Now().Format(time.RFC1123)
			app.EmitEvent("time", now)
			time.Sleep(time.Second)
		}
	}()

	systemTray := app.NewSystemTray()
	menu := app.NewMenu()
	menu.Add("Quit").OnClick(func(data *application.Context) {
		app.Quit()
	})
	systemTray.SetMenu(menu)

	if runtime.GOOS == "darwin" {
		systemTray.SetTemplateIcon(icons.SystrayMacTemplate)
	}

	systemTray.OnClick(func() {
		if windowShowing {
			hideWindow()
		} else {
			checkAndShowWindow()
		}
	})

	// Create initial window
	currentWindow = createWindow(app)
	windowShowing = true

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
