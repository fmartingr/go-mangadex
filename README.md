# go-mangadex

Mangadex API client in Golang.

## Chaching

...

## Usage

``` go
import (
    "log"

    "code.fmartingr.dev/fmartingr/go-mangadex"
)

func main() {
    // Retrieve manga information
    manga, err := mangadex.GetManga(123)
    if errManga != nil {
        log.Println("Error retrieving manga: %s", errManga)
    }

    // Retrieve a list of chapters 
    chaptersRequest := NewGetChaptersParams()
    chapters, errChapterList = manga.getChapters(chaptersRequest)
    if errChapterList != nil {
        log.Println("Error retrieving chapters page %d: %s", chaptersRequest.Page, errChapterList)
    }
        
    // Disables chache reads for requests beyond this point
    mangadex.DisableCache()

    // Retrieve a specific chapter detail
    // This will return more information than the list (the pages, server, etc)
    chapter, err := manga.GetChapter(1)
    if errChapter != nil {
        log.Println("Error retrieving chapter: %s", errChapter)
    }

    // Re-enables the cache
    mangadex.EnableCache()

    // Get all covers for this manga
    covers, errCovers := manga.GetCovers()
    if errCovers != nil {
        log.Println("Error retreiving covers: %s", errCovers)
    }
}
```
