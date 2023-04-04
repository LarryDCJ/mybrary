package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

const openLibrarySearchURL = "https://openlibrary.org/search.json"
const openLibraryAuthorSearchURL = "https://openlibrary.org/search.json"

type Book struct {
	CoverID          int      `json:"cover_i"`
	HasFullText      bool     `json:"has_fulltext"`
	EditionCount     int      `json:"edition_count"`
	Title            string   `json:"title"`
	AuthorName       []string `json:"author_name"`
	FirstPublishYear int      `json:"first_publish_year"`
	Key              string   `json:"key"`
	IA               []string `json:"ia"`
	AuthorKey        []string `json:"author_key"`
	PublicScanB      bool     `json:"public_scan_b"`
}

type SearchResult struct {
	NumFound int    `json:"num_found"`
	Docs     []Book `json:"docs"`
}

func getBookByTitle(c echo.Context) error {
	title := c.Param("title")
	if title == "" {
		return c.String(http.StatusBadRequest, "Title parameter is required")
	}

	query := fmt.Sprintf("title=%s", url.QueryEscape(title))
	resp, err := http.Get(fmt.Sprintf("%s?%s", openLibrarySearchURL, query))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var searchResult SearchResult
	err = json.Unmarshal(body, &searchResult)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if searchResult.NumFound == 0 {
		return c.String(http.StatusNotFound, "No books found with the given title")
	}

	return c.JSON(http.StatusOK, searchResult.Docs[0])
}

func getBooksByAuthor(c echo.Context) error {
	author := c.Param("author")
	if author == "" {
		return c.String(http.StatusBadRequest, "Author parameter is required")
	}

	query := fmt.Sprintf("author=%s", url.QueryEscape(author))
	resp, err := http.Get(fmt.Sprintf("%s?%s", openLibrarySearchURL, query))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Print(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var searchResult SearchResult
	err = json.Unmarshal(body, &searchResult)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if searchResult.NumFound == 0 {
		return c.String(http.StatusNotFound, "No books found with the given author")
	}

	return c.JSON(http.StatusOK, searchResult.Docs)
}

func main() {
	e := echo.New()

	e.GET("/books/:title", getBookByTitle)
	e.GET("/books/:author", getBooksByAuthor)

	err := e.Start(":8080")
	if err != nil {
		return
	}
}
