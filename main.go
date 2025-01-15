package main

import (
	"assignment4_SQA/pom"
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"log"
	"sync"
	"time"
)

var wg sync.WaitGroup

const (
	seleniumPath = "C:/chromedriver-win64/chromedriver.exe"
	port         = 5001
)

func main() {
	var opts []selenium.ServiceOption
	service, err := selenium.NewChromeDriverService(seleniumPath, port, opts...)
	if err != nil {
		log.Fatalf("Error starting the ChromeDriver service: %v", err)
	}
	defer service.Stop()

	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{
		Args: []string{
			"--start-maximized",
			"--disable-notifications",
			"--no-sandbox",
			"--disable-dev-shm-usage",
		},
	})
	waitGroup := 0

	driver1, _ := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	waitGroup += 1
	driver2, _ := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	waitGroup += 1
	driver3, _ := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	waitGroup += 1

	defer driver1.Quit()
	defer driver2.Quit()
	defer driver3.Quit()

	wg.Add(waitGroup)

	go testSearch(driver1)
	go testLoginLogoutOpencart(driver2)
	go bookFlightUsingPOM(driver3)

	wg.Wait()
}

func waitForElement(driver selenium.WebDriver, by, value string, timeout time.Duration) (selenium.WebElement, error) {
	var element selenium.WebElement
	var err error

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		element, err = driver.FindElement(by, value)
		if err == nil {
			return element, nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return nil, fmt.Errorf("element not found after %v: %v", timeout, err)
}

func testSearch(driver selenium.WebDriver) {
	defer wg.Done()

	err := driver.Get("https://www.google.com")
	if err != nil {
		log.Printf("Failed to load Google: %v", err)
		return
	}

	searchBox, err := waitForElement(driver, selenium.ByCSSSelector, "textarea[name='q']", 10*time.Second)
	if err != nil {
		log.Printf("Failed to find search box: %v", err)
		return
	}

	err = searchBox.SendKeys("Selenium with Go")
	if err != nil {
		log.Printf("Failed to enter search text: %v", err)
		return
	}

	err = searchBox.SendKeys(selenium.EnterKey)
	if err != nil {
		log.Printf("Failed to press Enter: %v", err)
		return
	}

	results, err := waitForElement(driver, selenium.ByXPATH, "//h3[1]", 10*time.Second)
	if err != nil {
		log.Printf("No results found: %v", err)
		return
	}

	title, err := results.Text()
	if err != nil {
		log.Printf("Failed to get result text: %v", err)
		return
	}

	fmt.Println("Task 1: First result title:", title)
}

func testLoginLogoutOpencart(driver selenium.WebDriver) {
	defer wg.Done()

	err := driver.Get("https://demo.opencart.com/index.php?route=account/login")
	if err != nil {
		log.Printf("Failed to load OpenCart: %v", err)
		return
	}

	email, err := waitForElement(driver, selenium.ByCSSSelector, "input#input-email", 10*time.Second)
	if err != nil {
		log.Printf("Failed to find email input: %v", err)
		return
	}

	password, err := waitForElement(driver, selenium.ByCSSSelector, "input#input-password", 10*time.Second)
	if err != nil {
		log.Printf("Failed to find password input: %v", err)
		return
	}

	err = email.SendKeys("220859@astanait.edu.kz")
	if err != nil {
		log.Printf("Failed to enter email: %v", err)
		return
	}

	err = password.SendKeys("TestPassword")
	if err != nil {
		log.Printf("Failed to enter password: %v", err)
		return
	}

	loginButton, err := waitForElement(driver, selenium.ByXPATH, `//*[@id="form-login"]/div[3]/button`, 10*time.Second)
	if err != nil {
		log.Printf("Failed to find login button: %v", err)
		return
	}

	err = loginButton.Click()
	if err != nil {
		log.Printf("Failed to click login button: %v", err)
		return
	}
	time.Sleep(3 * time.Second)
	err = loginButton.Click()
	if err != nil {
		log.Printf("Failed to click login button: %v", err)
		return
	}
	time.Sleep(3 * time.Second)
	time.Sleep(10 * time.Second)

	_, err = waitForElement(driver, selenium.ByXPATH, "//span[text()='My Account']", 10*time.Second)
	if err != nil {
		log.Printf("Login verification failed: My Account element not found: %v", err)
		return
	}

	log.Printf("Login verified successfully!")

	err = driver.Get("https://demo.opencart.com/index.php?route=common/home")
	if err != nil {
		log.Printf("Failed to navigate to home page: %v", err)
		return
	}
	time.Sleep(2 * time.Second)

	myAccountDropdown, err := waitForElement(driver, selenium.ByXPATH, "//span[text()='My Account']", 10*time.Second)
	if err != nil {
		log.Printf("Failed to find My Account dropdown after navigation: %v", err)
		return
	}

	err = myAccountDropdown.Click()
	if err != nil {
		log.Printf("Failed to click My Account dropdown: %v", err)
		return
	}
	time.Sleep(1 * time.Second)

	logoutLink, err := waitForElement(driver, selenium.ByXPATH, "//a[contains(@href, 'logout')]", 10*time.Second)
	if err != nil {
		log.Printf("Failed to find logout link: %v", err)
		return
	}

	err = logoutLink.Click()
	if err != nil {
		log.Printf("Failed to click logout link: %v", err)
		return
	}
	time.Sleep(5 * time.Second)

	_, err = waitForElement(driver, selenium.ByXPATH, "//a[contains(@href, 'login')]", 10*time.Second)
	if err != nil {
		log.Printf("Logout verification failed: Login link not found: %v", err)
		return
	}

	fmt.Println("Task 2: Login and logout cycle completed successfully")
}

func bookFlightUsingPOM(driver selenium.WebDriver) {
	homePage := pom.NewAviataHomePage(driver)

	if err := homePage.NavigateToHome(); err != nil {
		log.Printf("Error navigating to Aviata: %v", err)
		return
	}

	if err := homePage.SetFromCity("Astana"); err != nil {
		log.Printf("Error setting From city: %v", err)
		return
	}

	if err := homePage.SetToCity("Almaty"); err != nil {
		log.Printf("Error setting To city: %v", err)
		return
	}

	if err := homePage.SearchFlights(); err != nil {
		log.Printf("Error searching flights: %v", err)
		return
	}

	time.Sleep(10 * time.Second)
	log.Println("Flight booking test completed successfully.")
}
