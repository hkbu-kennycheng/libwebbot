package libwebbot

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tebeka/selenium"
)

var (
	SeleniumPath     = "/tmp/selenium-server-standalone-3.4.jar"
	SeleniumPort     = 8080
	ChromeDriverPath = "/usr/local/bin/chromedriver"
	WindowWidth      = 1600
	WindowHeight     = 2700
	IsDebug          = false
	ActionInterval   = time.Second
)

type KeyAction uint8

const (
	Click KeyAction = iota
	SendKeys
	Submit
	Clear
	ExecuteScript
	ElementScreenshot

	Go
	GoBack
	GoForward
	Refresh
	AddjQuery
	LogCurrentURL
	WindowScreenshot
)

type BotAction struct {
	XPath     string
	Action    KeyAction
	ActionArg string
}

func findVisibleElementByXPATH(wd selenium.WebDriver, xpath string) selenium.WebElement {
	if elements, err := wd.FindElements(selenium.ByXPATH, xpath); err == nil && len(elements) > 0 {
		for _, element := range elements {
			if isDisplayed, err := element.IsDisplayed(); err == nil && isDisplayed {
				return element
			}
		}
	}
	return nil
}

func saveFile(filename string, b []byte) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	if _, err := w.Write(b); err != nil {
		return err
	}
	return w.Flush()
}

func SetDebug(value bool) {
	IsDebug = value
}

func SetWindowSize(width, height int) {
	WindowWidth = width
	WindowHeight = height
}

func ChromeBot(url string, actions ...BotAction) error {
	opts := []selenium.ServiceOption{
		selenium.StartFrameBuffer(),             // Start an X frame buffer for the browser to run in.
		selenium.ChromeDriver(ChromeDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
		selenium.Output(os.Stderr),              // Output debug information to STDERR.
	}
	selenium.SetDebug(IsDebug)
	service, err := selenium.NewSeleniumService(SeleniumPath, SeleniumPort, opts...)
	if err != nil {
		return err
	}
	defer service.Stop()

	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", SeleniumPort))
	if err != nil {
		return err
	}
	defer wd.Quit()

	if err := wd.ResizeWindow("", WindowWidth, WindowHeight); err != nil {
		return err
	}

	if err := wd.Get(url); err != nil {
		return err
	}

	time.Sleep(ActionInterval)

	for _, action := range actions {
		if action.Action >= Go {
			switch action.Action {
			case Go:
				if err := wd.Get(action.ActionArg); err != nil {
					return err
				}
			case GoBack:
				if err := wd.Back(); err != nil {
					return err
				}
			case GoForward:
				if err := wd.Forward(); err != nil {
					return err
				}
			case Refresh:
				if err := wd.Refresh(); err != nil {
					return err
				}
			case AddjQuery:
				if _, err := wd.ExecuteScript(`if (!window.jQuery)  document.body.innerHTML += '<script src="https://code.jquery.com/jquery-1.12.4.min.js" integrity="sha256-ZosEbRLbNQzLpnKIkEdrPv7lOy9C27hHQ+Xp8a4MxAQ=" crossorigin="anonymous"></script>';`, nil); err != nil {
					return err
				}
			case LogCurrentURL:
				if currentURL, err := wd.CurrentURL(); err != nil {
					return err
				} else {
					log.Printf(action.ActionArg, currentURL)
				}
			case WindowScreenshot:
				if b, err := wd.Screenshot(); err != nil {
					return err
				} else {
					if err := saveFile(action.ActionArg, b); err != nil {
						return err
					}
				}
			}
			wd.AcceptAlert()
			time.Sleep(ActionInterval)
			continue
		}

		if element := findVisibleElementByXPATH(wd, action.XPath); element != nil {
			switch action.Action {
			case Click:
				if err := element.Click(); err != nil {
					return err
				}
			case SendKeys:
				if err := element.SendKeys(action.ActionArg); err != nil {
					return err
				}
			case Submit:
				if err := element.Submit(); err != nil {
					return err
				}
			case Clear:
				if err := element.Clear(); err != nil {
					return err
				}
			case ExecuteScript:
				if _, err := wd.ExecuteScript(action.ActionArg, nil); err != nil {
					return err
				}
			case ElementScreenshot:
				if b, err := element.Screenshot(true); err != nil {
					return err
				} else {
					if err := saveFile(action.ActionArg, b); err != nil {
						return err
					}
				}

			}
			wd.AcceptAlert()
			time.Sleep(ActionInterval)
		}
	}

	return nil
}

func ChromeGoAndFill() {

}
