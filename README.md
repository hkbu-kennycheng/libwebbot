# libwebbot
libwebbot is a library for me to run UAT of Web projects automatically.

## Usage
```go
package main

import (
	bot "github.com/hkbu-kennycheng/libwebbot"
	"log"
)

func main() {
	bot.SeleniumPath = "/path/selenium-server-standalone.jar"
	bot.SetDebug(true)
	bot.SetWindowSize(1280, 720)

	actions := []bot.BotAction{
		bot.BotAction{"//input[@name='q']", bot.SendKeys, "ChromeDriver"},
		bot.BotAction{"//input[@name='q']", bot.Submit, ""},
		bot.BotAction{"", bot.WindowScreenshot, "screenshot.png"},
	}

	if err := bot.ChromeBot("https://www.google.com/xhtml", actions...); err != nil {
		log.Fatalln(err.Error())
	}
}


```
