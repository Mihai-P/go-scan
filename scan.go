// scan_woocommerce.go
// https://www.devdungeon.com/content/web-scraping-go
// https://gist.github.com/salmoni/27aee5bb0d26536391aabe7f13a72494
// https://gobyexample.com/environment-variables
// https://blog.drewolson.org/dependency-injection-in-go
package main

import (
    "os"
    //"io"
    "fmt"
    "log"
    "strings"
    // "io/ioutil"
    // "database/sql"
    // _ "github.com/go-sql-driver/mysql"
    "net/http"
    "github.com/PuerkitoBio/goquery"
)

// var DB *sql.DB

// type SQLConfig struct {
//     Server       string
//     Username     string
//     Password     string
//     Db           string
// }

// func NewSQLConfig() *SQLConfig {
//   return &SQLConfig {
//     Server:      os.Getenv("MYSQL_SERVER"),
//     Username:    os.Getenv("MYSQL_USER"),
//     Password:    os.Getenv("MYSQL_PASSWORD"),
//     Db:          os.Getenv("MYSQL_DB"),
//   }
// }

// func OpenDBConnection(config *SQLConfig) (*sql.DB, error) {
//     connectionString := config.Username + ":" + config.Password + "@tcp(" + config.Server + ")/" + config.Db
//     // Open up our database connection.
//     return sql.Open("mysql", connectionString)
// }


// func insert(db *sql.DB, value string) {
//     // perform a db.Query insert
//     _, err := db.Exec("INSERT INTO test VALUES ( now(), '" + value + "' )")

//     // if there is an error inserting, handle it
//     if err != nil {
//         panic(err.Error())
//     }
// }

// Extract tag content from within a string
// good enough for simple tags, not to be used for complicated or nested ones
func processHtmlContent(content, search string) (result string) {
    // Find the start of the pot string
    elementStartIndex := strings.Index(content, search)
    if elementStartIndex == -1 {
        fmt.Println("No element found")
        os.Exit(0)
    }
    // The start index of the title is the index of the first
    // character, the < symbol. We don't want to include
    // <strong>Pot Size</strong> as part of the final value, so let's offset
    // the index by the number of characers in <strong>Pot Size</strong>
    elementStartIndex += len(search)
    
    strippedContent := content[elementStartIndex:]
    // Find the index of the closing tag
    elementEndIndex := strings.Index(strippedContent, "<")
    
    if elementEndIndex == -1 {
        result = strippedContent
    } else {
        result = strippedContent[:elementEndIndex]
    }

    return result
}

/*
// This will get called for each HTML element found
func processPlantLink(index int, element *goquery.Selection) {
    // See if the href attribute exists on the element
    href, exists := element.Attr("href")
    if exists {
        response, err := http.Get(href)
        if err != nil {
            log.Fatal(err)
        }
        defer response.Body.Close()
        // Create a goquery document from the HTTP response
        document, err := goquery.NewDocumentFromReader(response.Body)
        if err != nil {
            log.Fatal("Error loading HTTP response body. ", err)
        }
        name, exists := document.Find("h1.product_title.entry-title").First();
        fullDescription, exists := document.Find("div.entry-content").First();
        potSizes = processPlantDescription(fullDescription, "<strong>Pot Size</strong>")

        // Get the response body as a string
        dataInBytes, err := ioutil.ReadAll(fullDescription)
        descriptionContent := string(dataInBytes)

        document.Find("a.woocommerce-LoopProduct-link.woocommerce-loop-product__link").Each(processElement)
    }
}
*/

// // This will get called for each HTML element found
// func processListPage(index int, element *goquery.Selection) {
//     // See if the href attribute exists on the element
//     href, exists := element.Attr("href")
//     if exists {
//         fmt.Println(href)
//     }
//     for {
//         document.Find("a.woocommerce-LoopProduct-link.woocommerce-loop-product__link").Each(processElement)
        
//         href, exists := document.Find("a.next.page-numbers").First()
//         if !exists {
//             break
//         }
//         href
//     }
// }

// 
func processElement(index int, element *goquery.Selection) {
    // See if the href attribute exists on the element
    href, exists := element.Attr("href")
    if exists {
        fmt.Println(href)
    }
}

// Get  a goquery document from an url
func getDocument(url string) (document *goquery.Document) {
    fmt.Println("Visiting page:" + url)
    response, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    defer response.Body.Close()

    document, err = goquery.NewDocumentFromReader(response.Body)
    if err != nil {
        panic("Error loading HTTP response body.")
    }
    return
}

// Indentify all the product links on the page and pass them along to 
func processListingPage(document *goquery.Document) {
    // get all the product links from the page and process them.
    
}

func getNextListingPage(document *goquery.Document) (href string, exists bool) {
    href, exists = document.Find("a.next.page-numbers").First().Attr("href")
    return
}

// Go to the listing page, process it and move to the next link. 
// This is done by pressing the next page link.
// Sto processing when there is no next page link.
func visitListingPage(url string) {
    // See if the href attribute exists on the element
    for {
        document := getDocument(url)
        document.Find("a.woocommerce-LoopProduct-link.woocommerce-loop-product__link").Each(visitProductPage)
        nextPageURL, exists := getNextListingPage(document)
        
//        if !exists {
        if exists {
            break
        } else {
            url = nextPageURL
        }
    }
}

// This will get called for each HTML element found
func visitProductPage(index int, element *goquery.Selection) {
    // See if the href attribute exists on the element
    href, exists := element.Attr("href")
    if !exists {
        log.Fatal("url does not exist")
    }
    document := getDocument(href)
    title := document.Find("h1").First().Contents().Text()
    fmt.Println("Page Title: " + title)

    details := getTableCellValues(document.Find("table.shop_attributes").First())
    potSize, ok := details["Pot Size"]
    if ok {
        potSizes := strings.Split(potSize, ",")
        for i := range potSizes {
            potSizes[i] = strings.TrimSpace(potSizes[i])
        }
        fmt.Println(potSizes)
    }
}

func getTableCellValues(tablehtml *goquery.Selection) map[string]string {
    rows := make(map[string]string) 
    tablehtml.Find("tr").Each(func(indextr int, rowhtml *goquery.Selection) {
        heading := rowhtml.Find("th").First().Contents().Text()
        value := rowhtml.Find("td").First().Find("p").First().Contents().Text()
        rows[heading] = value
    })
    return rows
}

func main() {
    visitListingPage(os.Getenv("MYSQL_SERVER"))
}
