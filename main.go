package main

import (
	"encoding/xml"
	"errors"
	"flag"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ethanefung/linkparser"
)

var baseUrl string

var noneError error = errors.New("none left")

type URL struct {
	Loc string `xml:"loc"`
}

type UrlSet struct {
	XMLName xml.Name `xml:"urlset"`
	XMLNS   string   `xml:"xmlns,attr"`
	URLs    []URL    `xml:"url"`
}

type SiteMap struct {
	URLSet UrlSet `xml:"urlset"`
}

func main() {
	flag.StringVar(&baseUrl, "url", "", "the domain to generate a sitemap for")

	flag.Parse() // don't forget to parse!

	hostURL, err := url.Parse(baseUrl)
	if err != nil {
		panic(errors.New("Invalid url"))
	}
	if hostURL.String() == "" {
		panic(errors.New("please provide a `-url` with a valid url path"))
	}
	if hostURL.Scheme == "" {
		panic(errors.New("please use the http scheme in the provided url"))
	}

	lp := linkparser.New()

	urlset := UrlSet{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  []URL{{hostURL.String()}},
	}
	q := newQueue()
	q.enqueue(hostURL.String())

	for !q.empty() {
		currUrl, err := q.dequeue()
		if err != nil {
			panic(err)
		}
		links, err := getUrls(lp, currUrl)
		if err != nil {
			continue
		}
		for _, link := range links {
			u, err := url.Parse(link.Href)
			if err != nil {
				continue
			}
			if u.Scheme == "" || u.Scheme == "https" {
				u.Scheme = "http"
			}
			if u.Hostname() == "" {
				u.Host = hostURL.Hostname()
			}
			u.Path = strings.TrimSuffix(u.Path, "/")
			if u.RawQuery != "" {
				u.RawQuery = ""
			}
			if u.Fragment != "" {
				u.Fragment = ""
			}

			curr := u.String()
			if u.Scheme == "http" && u.Hostname() == hostURL.Hostname() && !q.queued(curr) {
				urlset.URLs = append(urlset.URLs, URL{curr})
				q.enqueue(curr)
			}
		}
	}

	b, err := xml.MarshalIndent(urlset, "", "    ")
	if err != nil {
		panic(err)
	}
	os.Stdout.WriteString(xml.Header)
	os.Stdout.Write(b)
}

type queue interface {
	empty() bool
	queued(path string) bool
	enqueue(path string) error
	dequeue() (string, error)
}

func newQueue() queue {
	return &pathQueue{
		seen: NewTrie(),
	}
}

type pathQueue struct {
	seen Trie
	head *node
	tail *node
}

func (q *pathQueue) empty() bool {
	return q.head == nil
}

func (q *pathQueue) queued(path string) bool {
	return q.seen.Search(path)
}

func (q *pathQueue) enqueue(path string) error {
	q.seen.Insert(path)
	n := &node{
		val: path,
	}
	if q.head == nil {
		q.head = n
		q.tail = n
	} else if q.tail.next != nil {
		return errors.New("tried to assign a tail to a preexisting node")
	} else {
		q.tail.next = n
		q.tail = q.tail.next
	}
	return nil
}

func (q *pathQueue) dequeue() (string, error) {
	if q.head == nil {
		return "", noneError
	}
	val := q.head.val
	q.head = q.head.next
	if q.head == nil {
		q.tail = nil
	}
	return val, nil
}

type node struct {
	val  string
	next *node
}

func getUrls(lp linkparser.LinkParser, url string) ([]linkparser.Link, error) {
	none := []linkparser.Link{}

	res, err := http.Get(url)
	ctype := res.Header.Get("Content-Type")

	if !strings.Contains(ctype, "text/html") {
		return none, nil
	}

	if err != nil {
		return none, err
	}
	data, err := ioutil.ReadAll(res.Body)
	sr := strings.NewReader(string(data))

	err = lp.UseReader(sr)
	if err != nil {
		return none, err
	}

	links, err := lp.Parse()
	if err != nil {
		return none, err
	}
	return links, err
}
