##Titlecase

This is a production-quality package made for cleaning and formatting book titles, but it can be used for titlecasing anything. It is far better than any of the other titlecasing packages available at the time of writings.

This package is made for book titles from the pre-Internet era and so it does not, in its current state, support domain names. For example: `look at this cool thing i found on google.com` would be changed to `Look at This Cool Thing I Found on Google. Com`. I can add support for domain names on request, but I have no use for it which is why I didn't add it already, and I have no idea if anyone ever uses any of these packages I write.

##Features

* Supports multiple languages: English, French, German, Italian, Spanish, Portuguese & Generic
* Supports contractions
* Supports initials
* Supports academic honors (M.D., Ph.D, etc.)
* Supports common abbreviations (USA, USSR, YMCA, etc.)
* Supports Roman numerals, without mistaking words for roman numerals
* Supports hyphenation and slashes
* Repairs grammatical errors in English
* Redetermines whitespace
* Converts or strips inappropriate punctuation
* Decodes all HTML entities
* Fully UTF8 compliant
* Written for speed and efficiency - no regular expressions, minimal looping

##Installation

    go get github.com/AlasdairF/Titlecase

##Example 1

this is a title by a. forsythe, made to demonstrate a example. WRITTEN ON 5TH DECEMBER, 2014.

    formatted := titlecase.English(` this is a title by a. forsythe, made to demonstrate a example. WRITTEN ON 5TH DECEMBER, 2014.`)
    
This Is a Title by A. Forsythe, Made to Demonstrate an Example; Written on 5th December, 2014

##Example 2

a very many exceptions exist. e.g. ussr IS FIXED. some guy with a ph.d knows what's going on. people forget rules like: on an in.

    formatted := titlecase.English(`a very many exceptions exist. e.g. ussr IS FIXED. some guy with a ph.d knows what's going on. people forget rules like: on on on.`)
    
A Very Many Exceptions Exist; E.g. USSR Is Fixed; Some Guy With a Ph.D Knows What's Going On; People Forget Rules Like: On an In

##Example 3

della corte d'appello di roma nell'anno mdccclxxxi

    formatted := titlecase.Italian(`   della corte d'appello di roma nell'anno mdccclxxxi `)

Della Corte d'Appello di Roma nell'Anno MDCCCLXXXI
