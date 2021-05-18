package tests

import (
	"fmt"
	"github.com/uadmin/uadmin/utils"
	"io"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/html"

	"github.com/uadmin/uadmin/model"
)

type TestStruct struct {
	model.Model
	Name         string
	Children     []TestStruct `gorm:"foreignKey:ID"`
	Parent       *TestStruct
	ParentID     uint
	OtherModel   TestStruct1
	OtherModelID uint
}

type TestStruct1 struct {
	model.Model
	Name  string `uadmin:"search"`
	Value int
}

type TestType int

func (TestType) Active() TestType {
	return 1
}

func (TestType) Inactive() TestType {
	return 2
}

type TestStruct2 struct {
	model.Model
	Name           string
	Count          int
	Value          float64
	Start          time.Time
	End            *time.Time
	Type           TestType
	OtherModel     TestStruct1
	OtherModelID   uint
	AnotherModel   *TestStruct1
	AnotherModelID uint
	Active         bool
	Hidden         string `uadmin:"list_exclude"`
}

func traverse(n *html.Node, tag string) (string, map[string]string, bool) {
	if isTagElement(n, tag) {
		tempMap := map[string]string{}
		for i := range n.Attr {
			tempMap[n.Attr[i].Key] = n.Attr[i].Val
		}
		if n.FirstChild == nil {
			return "", tempMap, true
		}
		return strings.TrimSpace(n.FirstChild.Data), tempMap, true
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, attr, ok := traverse(c, tag)
		if ok {
			return result, attr, ok
		}
	}

	return "", map[string]string{}, false
}

func getHTMLTag(r io.Reader, tag string) (string, map[string]string, bool) {
	doc, err := html.Parse(r)
	if err != nil {
		utils.Trail(utils.ERROR, "Fail to parse html")
		return "", map[string]string{}, false
	}

	return traverse(doc, tag)
}

func isTagElement(n *html.Node, tag string) bool {
	return n.Type == html.ElementNode && n.Data == tag
}

func tagSearch(n *html.Node, tag string, path string, index int) ([]string, []string, []map[string]string) {
	paths := []string{}
	content := []string{}
	attr := []map[string]string{}

	if path == "" {
		if n.Data != "" {
			path = fmt.Sprintf("%s[%d]", n.Data, index)
		}
	} else {
		path = path + "/" + fmt.Sprintf("%s[%d]", n.Data, index)
	}

	if isTagElement(n, tag) {
		if n.FirstChild == nil {
			content = append(content, "")
		} else {
			content = append(content, strings.TrimSpace(n.FirstChild.Data))
		}
		paths = append(paths, path)
		tempMap := map[string]string{}
		for i := range n.Attr {
			tempMap[n.Attr[i].Key] = n.Attr[i].Val
		}
		attr = append(attr, tempMap)
	}

	index = 0
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		childPaths, childContent, childAttr := tagSearch(c, tag, path, index)
		paths = append(paths, childPaths...)
		content = append(content, childContent...)
		attr = append(attr, childAttr...)
		if c.Type == html.ElementNode {
			index++
		}
	}
	return paths, content, attr
}

// func getHTMLTagList(r io.Reader, tag string) (paths []string, content []string, attr []map[string]string) {
// 	doc, err := html.Parse(r)
// 	if err != nil {
// 		Trail(ERROR, "Failed to parse html")
// 		return
// 	}
// 	return tagSearch(doc, tag, "", 0)
// }

func parseHTML(r io.Reader, t *testing.T) (*html.Node, error) {
	doc, err := html.Parse(r)
	if err != nil {
		t.Errorf("Unable to parse html stream")
	}
	return doc, err
}

type TestModelA struct {
	model.Model
	Name string
}

type TestModelB struct {
	model.Model
	Name         string     `uadmin:"help:This is a test help message;search;list_exclude"`
	ItemCount    int        `uadmin:"max:5;min:1;format:%03d;required;read_only:true,edit"`
	Phone        string     `uadmin:"default_value:09;pattern:[0-9+]{7,15};pattern_msg:invalid phone number;encrypt"`
	Active       bool       `uadmin:"hidden;read_only"`
	OtherModel   TestModelA `uadmin:"categorical_filter;filter;read_only:new"`
	OtherModelID uint
	ModelAList   []TestModelA `gorm:"foreignKey:ID"`
	Parent       *TestModelB
	ParentID     uint
	Email        string  `uadmin:"email"`
	Greeting     string  `uadmin:"multilingual"`
	Image        string  `uadmin:"image;upload_to:/media/home/me/images/"`
	File         string  `uadmin:"file;upload_to:/media/home/me/files"`
	Secret       string  `uadmin:"password"`
	Description  string  `uadmin:"html"`
	URL          string  `uadmin:"link"`
	Code         string  `uadmin:"code"`
	P1           int     `uadmin:"progress_bar"`
	P2           float64 `uadmin:"progress_bar"`
	P3           float64 `uadmin:"progress_bar:1.0"`
	P4           float64 `uadmin:"progress_bar:1.0:red"`
	P5           float64 `uadmin:"progress_bar:1.0:#f00"`
	P6           float64 `uadmin:"progress_bar:0.3:red,0.7:yellow,1.0:lime"`
	Price        float64 `uadmin:"money"`
	List         testList
}

type TestApproval struct {
	model.Model
	Name        string     `uadmin:"approval"`
	Start       time.Time  `uadmin:"approval"`
	End         *time.Time `uadmin:"approval"`
	Count       int        `uadmin:"approval"`
	Price       float64    `uadmin:"approval"`
	List        testList   `uadmin:"approval"`
	TestModel   TestModelA `uadmin:"approval"`
	TestModelID uint
	Active      bool `uadmin:"approval"`
}

var receivedEmail string

// Method__List__Form is a method to test method based properties for models
func (TestModelB) Method__List__Form() string {
	return "Value"
}

type testList int

func (testList) A() testList {
	return 1
}

func startEmailServer() {
	// Listen for incoming connections.
	l, err := net.Listen("tcp", "localhost:2525")
	if err != nil {
		utils.Trail(utils.ERROR, "listening: %s", err)
		return
	}

	// Close the listener when the application closes.
	defer l.Close()

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			utils.Trail(utils.ERROR, "startEmailServer error accepting connection. %s", err)
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.

	conn.Write([]byte("220 smtp.example.com ESMTP Postfix (Ubuntu)\n"))

	_, err := conn.Read(buf)
	if err != nil {
		utils.Trail(utils.ERROR, "reading: %s", err)
	}

	buf = make([]byte, 1024)
	conn.Write([]byte(`250-smtp.example.com
250-AUTH LOGIN PLAIN
250-PIPELINING
250-SIZE 102400000
250-VRFY
250-ETRN
250-ENHANCEDSTATUSCODES
250-8BITMIME
250 DSN
`))
	_, err = conn.Read(buf)
	if err != nil {
		utils.Trail(utils.ERROR, "reading: %s", err)
	}
	conn.Write([]byte("235 Authentication succeeded\n"))

	buf = make([]byte, 1024)
	_, err = conn.Read(buf)
	if err != nil {
		utils.Trail(utils.ERROR, "reading: %s", err)
	}
	conn.Write([]byte("250 2.1.0 Ok\n"))

	buf = make([]byte, 1024)
	_, err = conn.Read(buf)
	if err != nil {
		utils.Trail(utils.ERROR, "reading: %s", err)
	}
	conn.Write([]byte("250 2.1.5 Ok\n"))

	buf = make([]byte, 1024)
	_, err = conn.Read(buf)
	if err != nil {
		utils.Trail(utils.ERROR, "reading: %s", err)
	}
	conn.Write([]byte("354 End data with <CR><LF>.<CR><LF>\n"))

	buf = make([]byte, 1024)
	_, err = conn.Read(buf)
	if err != nil {
		utils.Trail(utils.ERROR, "reading: %s", err)
	}
	conn.Write([]byte("250 2.0.0 Ok: queued as 16756A11026D\n"))
	receivedEmail = string(buf)

	buf = make([]byte, 1024)
	_, err = conn.Read(buf)
	if err != nil {
		utils.Trail(utils.ERROR, "reading: %s", err)
	}
	conn.Write([]byte("221 2.0.0 Bye\n"))

	// Close the connection when you're done with it.
	conn.Close()
}
