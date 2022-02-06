package xmlprefixer

import (
	"bytes"
	"github.com/antchfx/xmlquery"
	"strings"
)

func Parse(data []byte) ([]byte, error) {
	doc, err := xmlquery.Parse(bytes.NewBuffer(data))
	if err != nil {
		return data, err
	}

	nodes := doc.SelectElements("//*[@xmlns]")
	for _, node := range nodes {
		injectXmlParentPrefix(node)
	}

	return []byte(doc.OutputXML(false)), nil
}

func findXmlPrefix(ns string) string {
	dots := strings.Split(ns, ":")
	slashes := strings.Split(ns, "/")

	var last string
	if len(dots) > len(slashes) {
		last = dots[len(dots)-1]
		slashes = strings.Split(last, "/")
		last = slashes[len(slashes)-1]

	} else {
		last = slashes[len(slashes)-1]
		dots = strings.Split(last, ":")
		last = dots[len(dots)-1]
	}

	var cut int
	for i := len(last) - 1; i > 0; i-- {
		if last[i] == '.' {
			cut++
			continue
		}

		if last[i] == '-' {
			cut++
			continue
		}

		if last[i] >= '0' && last[i] <= '9' {
			cut++
			continue
		}

		break
	}

	return last[:len(last)-cut]
}

func injectXmlParentPrefix(parent *xmlquery.Node) {
	if parent == nil || parent.Type != xmlquery.ElementNode {
		return
	}

	prefix := findXmlPrefix(parent.NamespaceURI)

	if prefix == "" || prefix == "epp" {
		return
	}

	parent.Prefix = prefix

	for k, attr := range parent.Attr {
		if attr.Name.Local == "xmlns" {
			parent.Attr[k].Name.Local = "xmlns:" + prefix
		}
	}

	injectXmlChildPrefix(parent.FirstChild.NextSibling, prefix)
}

func injectXmlChildPrefix(child *xmlquery.Node, prefix string) {
	if child == nil || child.Type != xmlquery.ElementNode {
		return
	}

	child.Prefix = prefix

	if child.FirstChild != nil && child.FirstChild.NextSibling != nil {
		injectXmlChildPrefix(child.FirstChild.NextSibling, prefix)
	}

	if child.NextSibling != nil && child.NextSibling.NextSibling != nil {
		injectXmlChildPrefix(child.NextSibling.NextSibling, prefix)
	}
}
