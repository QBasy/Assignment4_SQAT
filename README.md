# Assignment 4 (Practice with Advanced Selenium WebDriver)

## Author
- **Name:** Sayat Adilkhanov
- **Group:** SE-2215

## Technology Stack
- **Programming Language:** Go
- **Libraries & Frameworks:**
  - "github.com/tebeka/selenium"
  - "github.com/tebeka/selenium/chrome"

## Project Overview
This project implements a Page Object Model (POM) for automated testing of the Aviata.kz booking system using Selenium WebDriver with Go. The implementation focuses on creating a structured and maintainable test framework for the booking functionality.

## Project Structure

### Main Structures
The project is built around two primary structures in `aviata.go`:

```go
type AviataSearchForm struct {
    driver selenium.WebDriver
    logger *log.Logger
}

type SearchFormElements struct {
    FromInput         selenium.WebElement
    ToInput          selenium.WebElement
    FromCityCode     selenium.WebElement
    ToCityCode       selenium.WebElement
    DepartureDate    selenium.WebElement
    ReturnDate       selenium.WebElement
    PassengerClass   selenium.WebElement
    SearchButton     selenium.WebElement
    ComplexRouteBtn  selenium.WebElement
    AnywhereBtn      selenium.WebElement
    BookingCheckbox  selenium.WebElement
}
```

### Core Functions

#### Constructor
```go
func NewAviataSearchForm(driver selenium.WebDriver, logger *log.Logger) *AviataSearchForm {
    return &AviataSearchForm{
        driver: driver,
        logger: logger,
    }
}
```

#### Element Retrieval
```go
func (a *AviataSearchForm) GetElements() (*SearchFormElements, error) {
    elements := &SearchFormElements{}
    var err error

    elements.FromInput, err = a.driver.FindElement(selenium.ByCSSSelector, ".from-input")
    if err != nil {
        return nil, err
    }

    elements.ToInput, err = a.driver.FindElement(selenium.ByCSSSelector, ".to-input")
    if err != nil {
        return nil, err
    }

    elements.FromCityCode, err = a.driver.FindElement(selenium.ByCSSSelector, ".from-city-code")
    if err != nil {
        return nil, err
    }

    elements.ToCityCode, err = a.driver.FindElement(selenium.ByCSSSelector, ".to-city-code")
    if err != nil {
        return nil, err
    }

    elements.DepartureDate, err = a.driver.FindElement(selenium.ByCSSSelector, ".departure-date")
    if err != nil {
        return nil, err
    }

    elements.ReturnDate, err = a.driver.FindElement(selenium.ByCSSSelector, ".return-date")
    if err != nil {
        return nil, err
    }

    elements.PassengerClass, err = a.driver.FindElement(selenium.ByCSSSelector, ".passenger-class")
    if err != nil {
        return nil, err
    }

    elements.SearchButton, err = a.driver.FindElement(selenium.ByCSSSelector, ".search-button")
    if err != nil {
        return nil, err
    }

    elements.ComplexRouteBtn, err = a.driver.FindElement(selenium.ByCSSSelector, ".complex-route-btn")
    if err != nil {
        return nil, err
    }

    elements.AnywhereBtn, err = a.driver.FindElement(selenium.ByCSSSelector, ".anywhere-btn")
    if err != nil {
        return nil, err
    }

    elements.BookingCheckbox, err = a.driver.FindElement(selenium.ByCSSSelector, ".booking-checkbox")
    if err != nil {
        return nil, err
    }

    return elements, nil
}
```

### Search Flight Implementation
#### starting search testing, gets dates as an input data, if test succeed returns nil (null), if not returns error
```go
func (a *AviataSearchForm) SearchFlight(from, to, departDate string) error {
    elements, err := a.GetElements()
    if err != nil {
        return err
    }

    err = elements.FromInput.SendKeys(from)
    if err != nil {
        return err
    }

    err = a.waitForCitySuggestions()
    if err != nil {
        return err
    }

    err = a.selectCityFromDropdown()
    if err != nil {
        return err
    }

    err = elements.ToInput.SendKeys(to)
    if err != nil {
        return err
    }

    a.logger.Println("Handle destination city")
    err = a.waitForCitySuggestions()
    if err != nil {
        return err
    }

    err = a.selectCityFromDropdown()
    if err != nil {
        return err
    }

    err = elements.DepartureDate.SendKeys(departDate)
    if err != nil {
        return err
    }

    err = elements.SearchButton.Click()
    if err != nil {
        return err
    }

    return nil
}
```

### Helper Functions
#### works instead on WebWaitDriver, because there is no same in go
```go
func (a *AviataSearchForm) waitForCitySuggestions() error {
    timeout := 10 * time.Second
    start := time.Now()
    for {
        if time.Since(start) > timeout {
            return fmt.Errorf("timeout waiting for city suggestions")
        }
        _, err := a.driver.FindElement(selenium.ByCSSSelector, ".city-suggestion")
        if err == nil {
            return nil
        }
        time.Sleep(100 * time.Millisecond)
    }
}

func (a *AviataSearchForm) selectCityFromDropdown() error {
    suggestions, err := a.driver.FindElements(selenium.ByCSSSelector, ".city-suggestion")
    if err != nil {
        return err
    }
    if len(suggestions) > 0 {
        return suggestions[0].Click()
    }
    return fmt.Errorf("no city suggestions found")
}
```

### Main Test Implementation in main.go
```go
func bookFlightUsingPOM(driver selenium.WebDriver) {
    defer wg.Done()
    
    logTest3 := log.New(os.Stdout, "AviataBooking Test: ", log.LstdFlags|log.Lshortfile)
    
    err := driver.Get("https://aviata.kz")
    if err != nil {
        logTest3.Printf("Failed to load Aviata: %v", err)
        return
    }
    
    time.Sleep(3 * time.Second)
    
    searchForm := pom.NewAviataSearchForm(driver, logTest3)
    
    _, err = searchForm.GetElements()
    if err != nil {
        logTest3.Printf("Failed to get form elements: %v", err)
        return
    }
    logTest3.Println("Successfully initialized all form elements")
    
    logTest3.Println("Starting flight search...")
    err = searchForm.SearchFlight("Астана", "Алматы", "2025-01-17")
    if err != nil {
        logTest3.Printf("Failed to perform search: %v", err)
        return
    }
    logTest3.Println("Successfully submitted flight search")
    
    time.Sleep(5 * time.Second)
    
    _, err = driver.FindElement(selenium.ByCSSSelector, ".search-results")
    if err != nil {
        logTest3.Printf("Failed to find search results: %v", err)
        return
    }
    logTest3.Println("Search results loaded successfully")
    
    logTest3.Println("Task 3: Aviata search test completed successfully")
}
```

### Main Function
```go
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
```

## Test Results
The implementation successfully demonstrates:
1. Proper initialization of page elements
2. Execution of flight search functionality
3. Verification of search results

My output from terminal after test execution:
```
AviataBooking Test: 2025/01/16 22:01:07 main.go:243: Successfully initialized all form elements
AviataBooking Test: 2025/01/16 22:01:07 main.go:245: Starting flight search...
AviataBooking Test: 2025/01/16 22:01:27 aviata_home.go:135: Handle destination city
AviataBooking Test: 2025/01/16 22:01:50 main.go:251: Successfully submitted flight search
AviataBooking Test: 2025/01/16 22:01:55 main.go:260: Search results loaded successfully
AviataBooking Test: 2025/01/16 22:01:55 main.go:262: Task 3: Aviata search test completed successfully
```

### You can get better report in SQAT4.docx
