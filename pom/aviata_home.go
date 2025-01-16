package pom

import (
	"fmt"
	"github.com/tebeka/selenium"
	"log"
	"time"
)

type AviataSearchForm struct {
	driver selenium.WebDriver
	logger *log.Logger
}

type SearchFormElements struct {
	FromInput       selenium.WebElement
	ToInput         selenium.WebElement
	FromCityCode    selenium.WebElement
	ToCityCode      selenium.WebElement
	DepartureDate   selenium.WebElement
	ReturnDate      selenium.WebElement
	PassengerClass  selenium.WebElement
	SearchButton    selenium.WebElement
	ComplexRouteBtn selenium.WebElement
	AnywhereBtn     selenium.WebElement
	BookingCheckbox selenium.WebElement
}

func NewAviataSearchForm(driver selenium.WebDriver, logger *log.Logger) *AviataSearchForm {
	return &AviataSearchForm{
		driver: driver,
		logger: logger,
	}
}

func (a *AviataSearchForm) GetElements() (*SearchFormElements, error) {
	time.Sleep(3 * time.Second)

	elements := &SearchFormElements{}
	var err error

	elements.FromInput, err = a.driver.FindElement(selenium.ByCSSSelector,
		"#search-route-0 input[placeholder='Откуда']")
	if err != nil {
		return nil, err
	}

	elements.ToInput, err = a.driver.FindElement(selenium.ByCSSSelector,
		"#search-route-0 input[placeholder='Куда']")
	if err != nil {
		return nil, err
	}

	elements.FromCityCode, err = a.driver.FindElement(selenium.ByCSSSelector,
		"#search-route-0 > div:first-child span.text-gray-500")
	if err != nil {
		return nil, err
	}

	elements.ToCityCode, err = a.driver.FindElement(selenium.ByCSSSelector,
		"#search-route-0 > div:nth-child(3) span.text-gray-500")
	if err != nil {
		return nil, err
	}

	elements.DepartureDate, err = a.driver.FindElement(selenium.ByID, "desktop-main-date-from")
	if err != nil {
		return nil, err
	}

	elements.ReturnDate, err = a.driver.FindElement(selenium.ByID, "desktop-main-date-to")
	if err != nil {
		return nil, err
	}

	elements.PassengerClass, err = a.driver.FindElement(selenium.ByCSSSelector,
		".extra-params div.relative")
	if err != nil {
		return nil, err
	}

	elements.SearchButton, err = a.driver.FindElement(selenium.ByCSSSelector,
		"button[type='submit']")
	if err != nil {
		return nil, err
	}

	elements.ComplexRouteBtn, err = a.driver.FindElement(selenium.ByCSSSelector,
		"button.text-purple-300")
	if err != nil {
		return nil, err
	}

	elements.AnywhereBtn, err = a.driver.FindElement(selenium.ByXPATH,
		"//button[contains(@class, 'text-purple-300')][.//span[contains(text(), 'Куда угодно')]]")
	if err != nil {
		return nil, err
	}

	elements.BookingCheckbox, err = a.driver.FindElement(selenium.ByCSSSelector,
		"label > input[type='checkbox']")
	if err != nil {
		return nil, err
	}

	return elements, nil
}

func (a *AviataSearchForm) SearchFlight(from, to, departDate string) error {
	elements, err := a.GetElements()
	if err != nil {
		return err
	}

	if err := elements.FromInput.Clear(); err != nil {
		return err
	}
	time.Sleep(1 * time.Second)

	if err := elements.FromInput.SendKeys(from); err != nil {
		return err
	}
	time.Sleep(3 * time.Second)

	if err := a.waitForCitySuggestions(); err != nil {
		return err
	}
	time.Sleep(1 * time.Second)

	if err := a.selectCityFromDropdown(); err != nil {
		return err
	}
	time.Sleep(2 * time.Second)

	a.logger.Printf("Handle destination city")
	if err := elements.ToInput.Clear(); err != nil {
		return err
	}
	time.Sleep(1 * time.Second)

	if err := elements.ToInput.SendKeys(to); err != nil {
		return err
	}
	time.Sleep(3 * time.Second)

	if err := a.waitForCitySuggestions(); err != nil {
		return err
	}
	time.Sleep(1 * time.Second)

	if err := a.selectCityFromDropdown(); err != nil {
		return err
	}
	time.Sleep(2 * time.Second)

	if err := elements.DepartureDate.Click(); err != nil {
		return err
	}
	time.Sleep(2 * time.Second)

	dateElement, err := a.driver.FindElement(selenium.ByXPATH,
		fmt.Sprintf("//div[@role='gridcell'][@data-date='%s']", departDate))
	if err != nil {
		return err
	}
	time.Sleep(1 * time.Second)

	if err := dateElement.Click(); err != nil {
		return err
	}
	time.Sleep(2 * time.Second)

	return elements.SearchButton.Click()
}

func (a *AviataSearchForm) waitForCitySuggestions() error {
	timeout := 10 * time.Second
	interval := 500 * time.Millisecond
	start := time.Now()

	for time.Since(start) < timeout {
		_, err := a.driver.FindElement(selenium.ByCSSSelector, "div[role='listbox']")
		if err == nil {
			return nil
		}
		time.Sleep(interval)
	}
	return fmt.Errorf("city suggestions not found within %v", timeout)
}

func (a *AviataSearchForm) selectCityFromDropdown() error {
	cityOption, err := a.driver.FindElement(selenium.ByCSSSelector, "div[role='option']")
	if err != nil {
		return err
	}
	return cityOption.Click()
}
