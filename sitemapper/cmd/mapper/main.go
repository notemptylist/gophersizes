package main

import (
	"bytes"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/notemptylist/gophersizes/linkparser/pkg/linkparse"
)

// getURL returns the body of a webpage at the specified URL.
func getUrl(url string) ([]byte, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

type status struct {
	parsed bool
	err    error
}

type sitemap map[string]*status

type UrlString struct {
	Url string `xml:"loc"`
}
type XMLDoc struct {
	XMLName xml.Name    `xml:"urlset"`
	Urls    []UrlString `xml:"url"`
}

type ExternalUrlError struct{}

func (m *ExternalUrlError) Error() string {
	return "external URL"
}

// parseableLink returns true if the link is to an HTML file or a path.
func parseableLink(url string) bool {
	if strings.HasSuffix(url, ".html") || strings.HasSuffix(url, "/") {
		return true
	}
	return false
}

// generateXML returns an XML document describing the links.
func generateXML(sm sitemap) string {

	xmlDoc := &XMLDoc{}
	for k, v := range sm {
		if !errors.Is(v.err, &ExternalUrlError{}) {
			xmlDoc.Urls = append(xmlDoc.Urls, UrlString{k})
		}
	}
	out, _ := xml.MarshalIndent(xmlDoc, " ", " ")
	return xml.Header + string(out)
}

func main() {

	website := flag.String("url", "", "The URL of the website to map out.")
	flag.Parse()
	if *website == "" {
		flag.Usage()
		os.Exit(1)
	}
	fmt.Printf("Mapping site: %s\n", *website)

	sm := sitemap{*website: &status{parsed: false, err: nil}}
	found := []string{*website}
	for {
		if len(found) > 0 {
			// Merge new urls into sitemap
			for _, v := range found {
				if _, ok := sm[v]; ok {
					log.Printf("Duplicate URL: %v\n", v)
				} else {
					log.Printf("Adding new URL: %v\n", v)
					sm[v] = &status{parsed: false, err: nil}
				}
			}
			// clear the slice keeping allocated memory
			found = found[:0]

			for currentUrl := range sm {
				// Only parse HTML files or directories
				if !parseableLink(currentUrl) {
					sm[currentUrl].err = errors.New("non parsable")
					continue
				}
				if !sm[currentUrl].parsed && sm[currentUrl].err == nil {
					if !strings.HasPrefix(currentUrl, *website) {
						sm[currentUrl].err = &ExternalUrlError{}
						continue
					}
					log.Printf("Parsing URL: %s\n", currentUrl)
					body, err := getUrl(currentUrl)
					if err != nil {
						sm[currentUrl].err = err
						continue
					}
					parser, err := linkparse.New(bytes.NewReader(body))
					if err != nil {
						sm[currentUrl].err = err
						continue
					}

					for _, link := range parser.EmitLinks() {
						log.Printf("Found link: %s", link.Href)
						if strings.HasPrefix(link.Href, *website) {
							found = append(found, link.Href)
						} else if strings.HasPrefix(link.Href, "/") {
							normalized, err := url.JoinPath(*website, link.Href)
							if err != nil {
								log.Printf("Corrupted url? %s", err)
								continue
							}
							log.Printf("Normalizing to: %s", normalized)
							found = append(found, normalized)
						} else {
							normalized, err := url.JoinPath(currentUrl, link.Href)
							if err != nil {
								log.Printf("Corrupted url? %s", err)
								continue
							}
							log.Printf("Normalizing to: %s", normalized)
							found = append(found, normalized)

						}
					}
					// Mark it as parsed so on the next iteration we will skip it.
					sm[currentUrl].parsed = true
				}
			}
		} else {
			// If no new links were found we quit
			break
		}
	}
	fmt.Println(generateXML(sm))
}
