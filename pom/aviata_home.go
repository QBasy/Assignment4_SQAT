package pom

import (
	"github.com/tebeka/selenium"
	"time"
)

// AviataHomePage defines the structure for Aviata's home page.
type AviataHomePage struct {
	driver selenium.WebDriver
}

// NewAviataHomePage initializes a new instance of AviataHomePage.
func NewAviataHomePage(driver selenium.WebDriver) *AviataHomePage {
	return &AviataHomePage{driver: driver}
}

// NavigateToHome navigates to the Aviata home page.
func (p *AviataHomePage) NavigateToHome() error {
	return p.driver.Get("https://aviata.kz")
}

// GetFromCityInput returns the 'From City' input element.
func (p *AviataHomePage) GetFromCityInput() (selenium.WebElement, error) {
	return p.driver.FindElement(selenium.ByXPATH, "//div[contains(@class, 'search')]//input[1]")
}

// GetToCityInput returns the 'To City' input element.
func (p *AviataHomePage) GetToCityInput() (selenium.WebElement, error) {
	return p.driver.FindElement(selenium.ByXPATH, `//*[@id="search-route-0"]/div[3]/label/input`)
}

// GetDatePicker returns the date picker element.
func (p *AviataHomePage) GetDatePicker() (selenium.WebElement, error) {
	return p.driver.FindElement(selenium.ByCSSSelector, ".date-picker-class")
}

// GetSearchButton returns the Search Button element.
func (p *AviataHomePage) GetSearchButton() (selenium.WebElement, error) {
	return p.driver.FindElement(selenium.ByCSSSelector, "button[type='submit']")
}

func (p *AviataHomePage) SetFromCity(city string) error {
	fromCityInput, err := p.GetFromCityInput()
	if err != nil {
		return err
	}
	if err = fromCityInput.Click(); err != nil {
		return err
	}
	time.Sleep(1 * time.Second)

	cityOption, err := p.driver.FindElement(selenium.ByXPATH, `//*[@id="search-route-0"]/div[1]/div[2]/div/ul/li[2]/div/div`)
	if err != nil {
		return err
	}
	return cityOption.Click()
}

func (p *AviataHomePage) SetToCity(city string) error {
	toCityInput, err := p.GetToCityInput()
	if err != nil {
		return err
	}
	if err = toCityInput.Click(); err != nil {
		return err
	}
	time.Sleep(2 * time.Second)

	// Adjust city option selector as necessary
	cityOption, err := p.driver.FindElement(selenium.ByXPATH, `//*[@id="search-route-0"]/div[3]/div[2]/div/ul/li[1]/div/div`)
	if err != nil {
		return err
	}
	return cityOption.Click()
}

// SelectDate sets the travel date.
func (p *AviataHomePage) SelectDate(date string) error {
	datePicker, err := p.GetDatePicker()
	if err != nil {
		return err
	}
	if err = datePicker.Click(); err != nil {
		return err
	}
	time.Sleep(1 * time.Second)

	dateOption, err := p.driver.FindElement(selenium.ByXPATH, `//div[@data-date='`+date+`']`)
	if err != nil {
		return err
	}
	return dateOption.Click()
}

func (p *AviataHomePage) SearchFlights() error {
	searchButton, err := p.GetSearchButton()
	if err != nil {
		return err
	}
	return searchButton.Click()
}
