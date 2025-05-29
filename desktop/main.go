package main

import (
	"embed"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/icons"
)

var windowShowing bool

func createWindow(app *application.App) {
	if windowShowing {
		return
	}
	// Log the time taken to create the window
	window := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Name:            "BeoEcho",
		AlwaysOnTop:     false,
		Hidden:          true,
		Frameless:       false,
		DevToolsEnabled: true,
		Windows: application.WindowsWindow{
			HiddenOnTaskbar: true,
		},
	})

	windowShowing = true

	window.OnWindowEvent(events.Common.WindowClosing, func(e *application.WindowEvent) {
		windowShowing = false
	})

	window.Show()
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
	app := application.New(application.Options{
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

	// Create a new window with the necessary options.
	// 'Title' is the title of the window.
	// 'Mac' options tailor the window when running on macOS.
	// 'BackgroundColour' is the background colour of the window.
	// 'URL' is the URL that will be loaded into the webview.
	// app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
	// 	Title: "Window 1",
	// 	Mac: application.MacWindow{
	// 		InvisibleTitleBarHeight: 50,
	// 		Backdrop:                application.MacBackdropTranslucent,
	// 		TitleBar:                application.MacTitleBarHiddenInset,
	// 	},
	// 	BackgroundColour: application.NewRGB(27, 38, 54),
	// 	URL:              "/",
	// })

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
		createWindow(app)
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
