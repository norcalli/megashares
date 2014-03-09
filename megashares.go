package megashares

import (
	"fmt"
	// "io"
	gq "github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

const (
	loginURL  = "http://d01.megashares.com/myms_login.php"
	searchURL = "http://www.megashares.com/search.php"
)

type Megashares struct {
	CookieJar *cookiejar.Jar
	Client    *http.Client
}

func New() *Megashares {
	cj, _ := cookiejar.New(nil)
	return &Megashares{cj, &http.Client{Jar: cj}}
	// m := new()
	// m.cookieJar, _ := cookiejar.New(nil)
	// m.Client := &http.Client{Jar: cookieJar}
}

type MegasharesEntry struct {
	Url      string
	Filename string
}

func (m *MegasharesEntry) String() string {
	return m.Filename
}

func EntryFromURL(url string) (*MegasharesEntry, error) {
	// TODO: Double check url format.
	i := strings.LastIndex(url, `fln=/`)
	if i < 0 {
		return nil, fmt.Errorf("Download url doesn't conform to 'fln=/'.")
	}
	return &MegasharesEntry{url, url[i+5:]}, nil
}

func (m *Megashares) Login(username, password string) error {
	values := url.Values{
		"mymslogin_name": {username},
		"mymspassword":   {password},
		"myms_login":     {"Login"},
	}
	r, err := m.Client.PostForm(loginURL, values)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer r.Body.Close()
	return nil
}

func (m *Megashares) Search(query string) ([]byte, error) {
	values := url.Values{
		"q":      {query},
		"simple": {"Submit"},
	}
	r, err := m.Client.PostForm(searchURL, values)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer r.Body.Close()
	return ioutil.ReadAll(r.Body)
}

func (m *Megashares) SearchResponse(query string) (*http.Response, error) {
	values := url.Values{
		"q":      {query},
		"simple": {"Submit"},
	}
	r, err := m.Client.PostForm(searchURL, values)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return r, nil
}

func (m *Megashares) SearchEntries(query string) ([]*MegasharesEntry, error) {
	values := url.Values{
		"q":      {query},
		"simple": {"Submit"},
	}
	r, err := m.Client.PostForm(searchURL, values)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	d, err := gq.NewDocumentFromResponse(r)
	if err != nil {
		log.Fatal(err)
	}
	urls := d.Find("div.float-r a img").Parent()
	entries := make([]*MegasharesEntry, urls.Length())
	urls.Each(func(i int, s *gq.Selection) {
		v, _ := s.Attr(`href`)
		entries[i] = EntryFromURL(v)
	})
	return entries, nil
}
