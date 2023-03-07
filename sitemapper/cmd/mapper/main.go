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
	Xmlns   string      `xml:"xmlns,attr"`
	XMLName xml.Name    `xml:"urlset"`
	Urls    []UrlString `xml:"url"`
}

// ExternalUrlError is an error indicating that the URL is pointing to an
// external domain.
type ExternalUrlError struct{}

func (e ExternalUrlError) Error() string {
	return "external URL"
}

// NotParseableError is an error indicating that the URL is not pointing to an HTML
// document which can be parsed for links.
type NotParseableError struct{}

func (e NotParseableError) Error() string {
	return "not parseable"
}

var bins = []string{".pdf", ".jpg", ".png", ".exe", ".bin", ".jpeg", ".so", ".js", ".gif", ".bmp"}

// parseableLink returns true if the link is not to a binary file.
func parseableLink(url string) bool {
	for _, ext := range bins {
		if strings.HasSuffix(url, ext) {
			return false
		}
	}
	return true
}

// externalUrl returns true of the URL points to a different domain than the given website.
func externalUrl(url, website string) bool {
	// if either url ends in "/", lets trim it to normalize them.
	if strings.HasPrefix(url, "http") {
		url = strings.TrimRight(url, "/")
		website = strings.TrimRight(website, "/")
		urlparts := strings.Split(url, "://")
		websiteparts := strings.Split(website, "://")
		// the website URL should contain less paths,
		// otherwise the url is not within its map.
		if strings.HasPrefix(urlparts[1], websiteparts[1]) {
			return false
		}
	}
	return true
}

// hrefs Returns links parsed from the passed in io.Reader.
func hrefs(r io.Reader, base string) []string {
	var ret []string
	for _, link := range linkparse.ParseLinks(r) {
		log.Printf("Found link: %s", link.Href)
		if strings.HasPrefix(link.Href, "/") || strings.HasPrefix(link.Href, "#") || !strings.HasPrefix(link.Href, "http") {
			b, err := url.Parse(base)
			if err != nil {
				log.Printf("Unparseable base path: %s", err)
				continue
			}
			p, err := url.Parse(link.Href)
			if err != nil {
				log.Printf("Corrupted url? %s", err)
				continue
			}
			normalized := b.ResolveReference(p).String()
			log.Printf("Normalizing to: %s", normalized)
			ret = append(ret, normalized)
		} else {
			ret = append(ret, link.Href)
		}
	}
	return ret
}

// generateXML returns an XML document describing the links.
func generateXML(sm sitemap) string {

	xmlDoc := &XMLDoc{Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9"}
	for k, v := range sm {
		if !errors.Is(v.err, ExternalUrlError{}) {
			xmlDoc.Urls = append(xmlDoc.Urls, UrlString{k})
		}
	}
	out, _ := xml.MarshalIndent(xmlDoc, " ", " ")
	return xml.Header + string(out)
}

// extractDomain returns the URL with just the domain part.
func extractDomain(website url.URL) string {
	website.Path = ""
	website.RawQuery = ""
	return website.String()
}

func main() {

	website := flag.String("url", "", "The URL of the website to map out.")
	flag.Parse()
	if *website == "" {
		flag.Usage()
		os.Exit(1)
	}
	log.Printf("Verifying URL: %s", *website)
	resp, err := http.Get(*website)
	if err != nil {
		panic(err)
	}
	// let's use the URL where the request ultimately ended up.
	reqUrl := resp.Request.URL
	fmt.Printf("Mapping site: %s\n", reqUrl.String())

	sm := sitemap{}

	found := []string{reqUrl.String()}
	domain := extractDomain(*reqUrl)
	for {
		if len(found) == 0 {
			// If no new links were found we quit
			break
		} else {
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
				// Any URL which has a status.parsed == true
				// or status.err != nil should be skipped early
				if sm[currentUrl].err != nil || sm[currentUrl].parsed {
					continue
				}
				// Now lets tag the URL with an error if we know ahead of time
				// it shouldn't be parsed.
				if !parseableLink(currentUrl) {
					sm[currentUrl].err = NotParseableError{}
					continue
				}
				// URLs which link to external domains should not be parsed.
				if externalUrl(currentUrl, domain) {
					sm[currentUrl].err = ExternalUrlError{}
					continue
				}

				// if checks pass, lets actually download the body of the URL
				log.Printf("Parsing URL: %s\n", currentUrl)
				body, err := getUrl(currentUrl)
				if err != nil {
					sm[currentUrl].err = err
					continue
				}
				found = append(found, hrefs(bytes.NewReader(body), currentUrl)...)
				// Mark it as parsed so on the next iteration we will skip it.
				sm[currentUrl].parsed = true
			}
		}
	}
	fmt.Println(generateXML(sm))
}
