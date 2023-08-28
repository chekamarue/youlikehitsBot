package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/playwright-community/playwright-go"
)

func main() {
	err := playwright.Install()
	if err != nil {
		log.Fatalf("could not install playwright: %v", err)
	}
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	launchOptions := playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	}
	browser, err := pw.Chromium.Launch(launchOptions)
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	defer func() {
		if err = browser.Close(); err != nil {
			log.Fatalf("could not close browser: %v", err)
		}
		if err = pw.Stop(); err != nil {
			log.Fatalf("could not stop Playwright: %v", err)
		}
	}()
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	var cookieData playwright.OptionalCookie
	cookieFile, err := os.ReadFile("cookies.json")
	if err != nil {
		fmt.Println(err)
	} else {
		json.Unmarshal(cookieFile, &cookieData)
		browser.Contexts()[0].AddCookies(cookieData)
	}
	if _, err = page.Goto("https://www.youlikehits.com/youtubenew2.php"); err != nil {
		log.Fatalf("could not goto: %v", err)
	}
	isNotLoggedIn, _ := page.GetByText("You're not logged in!").IsVisible()
	if isNotLoggedIn {
		if _, err = page.Goto("https://www.youlikehits.com/login.php"); err != nil {
			log.Fatalf("could not goto: %v", err)
		}
		fmt.Println("You're not logged in! Please login and press enter once done")
		fmt.Scanln()
		if _, err = page.Goto("https://www.youlikehits.com/youtubenew2.php"); err != nil {
			log.Fatalf("could not goto: %v", err)
		}
	}
	cookies, err := browser.Contexts()[0].Cookies()
	if err != nil {
		log.Fatalf("could not get cookies: %v", err)
	}
	cookiesJSON, err := json.Marshal(cookies)
	if err != nil {
		log.Fatalf("could not marshal cookies: %v", err)
	}
	os.WriteFile("cookies.json", cookiesJSON, 0644)
	isCaptchaShown, _ := page.GetByText("Are You Human?").IsVisible()
	if isCaptchaShown {
		fmt.Println("captcha is shown, please solve and press enter once done")
		fmt.Scanln()
	}
	for {
		viewButton := page.Locator(".followbutton")
		if err = viewButton.Click(); err != nil {
			log.Fatalf("could not click view button: %v", err)
		}
		for {
			timerTextIsVisible, _ := page.GetByText("Timer").IsVisible()
			if timerTextIsVisible {
				time.Sleep(1 * time.Second)
			} else {
				break
			}
		}
		time.Sleep(1 * time.Second)
	}
}
