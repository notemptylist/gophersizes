package linkparse

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

// findText recursively finds a text element.
func findText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}

	var text string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += findText(c)
	}
	return text
}

// inspectNode inspects a single node to determine if it is a link,
// it then recursively calls itself on the first child of the node and
// every sibling of that node.
func inspectNode(n *html.Node, padding string) []Link {
	var links []Link
	// if this node is a link, extract href and text
	// and create a Link value
	if n.Type == html.ElementNode && n.Data == "a" {
		// could be a valid link
		var link Link
		for _, a := range n.Attr {
			if a.Key == "href" {
				// found the address
				link.Href = a.Val
				break
			}
		}
		t := findText(n)
		link.Text = strings.Trim(t, "\t\n ")
		links = append(links, link)
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		// descend down the tree first
		links = append(links, inspectNode(child, padding+" ")...)
		// when this returns, we will follow siblings and descend their tree
	}
	return links
}

// ParseLinks returns all the links which were discovered by parsing the supplied page.
func ParseLinks(page io.Reader) []Link {
	var links []Link
	node, err := html.Parse(page)
	if err != nil {
		return nil
	}
	links = inspectNode(node, "")
	return links
}
