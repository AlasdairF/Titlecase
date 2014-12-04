##Titlecase

This is a production-quality package made for cleaning and formatting book titles, but it can be used for titlecasing anything.

##Features

* Supports multiple languages: English, French, German, Italian, Spanish, Portuguese & Generic
* Supports contractions
* Supports initials, titles and abbreviations
* Supports Roman numerals, with exceptions for real words that look like roman numerals
* Supports hypenatation and slashes
* Repairs grammatical errors in English
* Decodes all HTML entities
* Fully UTF8 compliant
* Written for speed and efficiency - no regular expressions, minimal looping

##Installation

    go get github.com/AlasdairF/Titlecase

##Usage

    unformatted := ` this is a title by a. forsythe, made to demonstrate a example. WRITTEN ON 5TH DECEMBER, 2014.`
    formatted := titlecase.Format(unformatted, titlecase.English)
    // This Is a Title by A. Forsythe, Made to Demonstrate an Example; Written on 5th December, 2014
    
    unformatted = `   della corte d'appello di roma nell'anno mdccclxxxi `
    formatted = titlecase.Format(unformatted, titlecase.Italian)
    // Della Corte d'Appello di Roma nell'Anno MDCCCLXXXI
