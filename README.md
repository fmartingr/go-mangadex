# go-mangadex

Mangadex API client in Golang.

## Installation

...

## Usage

``` go
import (
    mangadex "code.fmartingr.dev/fmartingr/go-mangadex"
)

func main() {
    manga, err := mangadex.GetManga(123)
    if err != nil {
        log.Println("[error] Retrieving manga: %s", err)
    }
    // manga.GetChapter(1)
    // manga.GetCovers()
    // manga.GetVolume()
}

```
