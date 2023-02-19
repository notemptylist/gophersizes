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

type LinkParser struct {
	r     io.Reader
	node  *html.Node
	links []Link
}

// New returns a LinkParser initialized with the first node parsed
// from the passed in io.Reader.
func New(r io.Reader) (*LinkParser, error) {
	node, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	return &LinkParser{r: r, node: node}, nil
}

// inspectNode inspects a single node to determine if it is a link,
// it then recursively calls itself on the first child of the node and
// every sibling of that node.
func (l *LinkParser) inspectNode(n *html.Node) {
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
		if n.FirstChild != nil {
			link.Text = strings.Trim(n.FirstChild.Data, "\t\n ")
			l.links = append(l.links, link)
		}
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		l.inspectNode(child)
	}
}

// EmitLinks returns all the links which were discovered during parsing.
func (l *LinkParser) EmitLinks() []Link {
	n := l.node
	l.inspectNode(n)
	return l.links
}
