package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/gocolly/colly"
)

type Album struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
}

type Song struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
}

// Able to retrieve new-release album from "melon.com" website
func main() {
	c := colly.NewCollector()

	albums := scrapeNewestAlubum(c)
	writeDataToJSON("newAlbums.json", albums)

	hipHopSongs := scrapeNewestHipHopSongs(c)
	writeDataToJSON("newHiphopSongs.json", hipHopSongs)
}

// Scrape Methods -------------------------------------------------

func scrapeNewestAlubum(c *colly.Collector) []Album {
	var albums []Album

	c.OnHTML("div.info", func(h *colly.HTMLElement) {
		albumName := h.ChildText("a.album_name")
		artistName := h.ChildText("span.checkEllipsis a.artist_name")

		if albumName != "" || artistName != "" {

			albumInstance := Album{
				Name:   removeBetweenBrackets(albumName),
				Artist: removeBetweenBrackets(artistName),
			}
			albums = append(albums, albumInstance)
		}

	})
	c.Visit("https://www.melon.com/new/album/index.htm")
	return albums
}

func scrapeNewestHipHopSongs(c *colly.Collector) []Song {

	var songs []Song
	var currentSong Song

	c.OnHTML("div.wrap_song_info", func(h *colly.HTMLElement) {
		h.ForEach("div", func(_ int, div *colly.HTMLElement) {
			// Find the nested <a> tag within the <span> tag
			a := div.ChildText("span a")

			// Check if the <a> tag is the first or second one based on your HTML structure
			if div.Index == 0 {
				currentSong.Title = a
			} else if div.Index == 1 {
				currentSong.Artist = a

				// Add to list of songs
				songs = append(songs, currentSong)

				// Reset the current song for the next iteration
				currentSong = Song{}
			}
		})
	})

	c.Visit("https://www.melon.com/genre/song_list.htm?gnrCode=GN0300&dtlGnrCode=")
	return songs
}

// Util Methods -----------------------------------------

func removeBetweenBrackets(input string) string {
	re := regexp.MustCompile(`\(.*?\)`)
	return re.ReplaceAllString(input, "")
}

func writeDataToJSON(fileName string, data any) {
	content, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error occured during json marshal process")
		return
	}

	os.WriteFile(fileName, content, 0644)
}
