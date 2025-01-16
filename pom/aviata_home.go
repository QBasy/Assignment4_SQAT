package pom

import (
	"github.com/tebeka/selenium"
	"time"
)

type AviataSearchForm struct {
	driver selenium.WebDriver
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

func NewAviataSearchForm(driver selenium.WebDriver) *AviataSearchForm {
	return &AviataSearchForm{
		driver: driver,
	}
}

func (a *AviataSearchForm) GetElements() (*SearchFormElements, error) {
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
		"#search-route-0 div:first-child span.text-gray-500")
	if err != nil {
		return nil, err
	}

	elements.ToCityCode, err = a.driver.FindElement(selenium.ByCSSSelector,
		"#search-route-0 div:nth-child(3) span.text-gray-500")
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
		".extra-params .relative")
	if err != nil {
		return nil, err
	}

	elements.SearchButton, err = a.driver.FindElement(selenium.ByCSSSelector,
		"button[type='submit']")
	if err != nil {
		return nil, err
	}

	elements.ComplexRouteBtn, err = a.driver.FindElement(selenium.ByCSSSelector,
		"button:contains('Составить сложный маршрут')")
	if err != nil {
		return nil, err
	}

	elements.AnywhereBtn, err = a.driver.FindElement(selenium.ByCSSSelector,
		"button:contains('Куда угодно')")
	if err != nil {
		return nil, err
	}

	elements.BookingCheckbox, err = a.driver.FindElement(selenium.ByCSSSelector,
		"input[type='checkbox']")
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
	if err := elements.FromInput.SendKeys(from); err != nil {
		return err
	}
	time.Sleep(500 * time.Millisecond)

	if err := elements.ToInput.Clear(); err != nil {
		return err
	}
	if err := elements.ToInput.SendKeys(to); err != nil {
		return err
	}
	time.Sleep(500 * time.Millisecond)

	if err := elements.DepartureDate.Click(); err != nil {
		return err
	}
	time.Sleep(500 * time.Millisecond)

	dateElement, err := a.driver.FindElement(selenium.ByCSSSelector,
		"[data-date='"+departDate+"']")
	if err != nil {
		return err
	}
	if err := dateElement.Click(); err != nil {
		return err
	}

	return elements.SearchButton.Click()
}

func (a *AviataSearchForm) ToggleBookingSearch() error {
	elements, err := a.GetElements()
	if err != nil {
		return err
	}
	return elements.BookingCheckbox.Click()
}

func (a *AviataSearchForm) OpenComplexRoute() error {
	elements, err := a.GetElements()
	if err != nil {
		return err
	}
	return elements.ComplexRouteBtn.Click()
}

func (a *AviataSearchForm) OpenAnywhereSearch() error {
	elements, err := a.GetElements()
	if err != nil {
		return err
	}
	return elements.AnywhereBtn.Click()
}
