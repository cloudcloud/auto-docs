package docs

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type pathStruct struct {
	CountPages int
	CountDirs  int
	InpDir     string
	Buffer     *bytes.Buffer
	ExpLogs    string
	M          string
	N          string
}

func TestUpdateFromPath(t *testing.T) {
	assert := assert.New(t)
	x := []pathStruct{
		{
			CountPages: 3,
			CountDirs:  2,
			InpDir:     getTestMarkdownDir(),
			Buffer:     bytes.NewBufferString(""),
			ExpLogs:    "",
			M:          "Standard processing should give 3 pages",
			N:          "Standard processing should give 2 dirs",
		},
		{
			CountPages: 2,
			CountDirs:  2,
			InpDir:     getTestMarkdownDir() + "first/",
			Buffer:     bytes.NewBufferString(""),
			ExpLogs:    "",
			M:          "Sub processing should give 2 pages",
			N:          "Sub processing should give 2 dirs",
		},
	}

	for _, a := range x {
		s := &Store{
			Dirs:  []*Dir{},
			Pages: map[string]*Page{},
		}
		log.SetOutput(a.Buffer)

		s.UpdateFromPath(a.InpDir)
		assert.Equal(a.CountPages, len(s.Pages), a.M)
		assert.Equal(a.CountDirs, len(s.Dirs), a.N)
		//assert.Equal(a.ExpLogs, a.Buffer.String())
	}
}

type walkerStruct struct {
	Path    string
	I       os.FileInfo
	Err     error
	ExpDirs []*Dir
	ExpErr  error
	Store   *Store
	M       string
}

type finfo struct {
	Dir bool
}

func (f *finfo) Name() string {
	return ""
}
func (f *finfo) Size() int64 {
	return int64(5)
}
func (f *finfo) Mode() os.FileMode {
	return os.FileMode(0755)
}
func (f *finfo) ModTime() time.Time {
	return time.Now()
}
func (f *finfo) IsDir() bool {
	return f.Dir
}
func (f *finfo) Sys() interface{} {
	return nil
}

func TestWalker(t *testing.T) {
	assert := assert.New(t)
	x := []walkerStruct{
		{
			Path:    "/",
			I:       &finfo{Dir: true},
			Err:     nil,
			ExpDirs: []*Dir{},
			ExpErr:  nil,
			Store:   &Store{Dirs: []*Dir{}},
			M:       "Base dir adding succeeds (by skipping).",
		},
		{
			Path:    "/",
			I:       &finfo{Dir: false},
			Err:     fmt.Errorf("..."),
			ExpDirs: []*Dir{},
			ExpErr:  fmt.Errorf("..."),
			Store:   &Store{Dirs: []*Dir{}},
			M:       "An error when processing won't add a dir.",
		},
		{
			Path:    getTestMarkdownDir() + "not-found-file.md",
			I:       &finfo{Dir: false},
			Err:     nil,
			ExpDirs: []*Dir{},
			ExpErr:  nil,
			Store:   &Store{Dirs: []*Dir{}, path: getTestMarkdownDir()},
			M:       "Erroneous file not added.",
		},
	}

	for _, a := range x {
		e := a.Store.walker(a.Path, a.I, a.Err)

		assert.Equal(a.ExpErr, e, a.M)
		assert.Equal(a.ExpDirs, a.Store.Dirs, a.M)
	}
}

type addStruct struct {
	CurDirs  []*Dir
	AddPath  string
	OrigPath string
	ExpDirs  []*Dir
	M        string
}

func TestAddToDir(t *testing.T) {
	assert := assert.New(t)
	x := []addStruct{
		{
			CurDirs:  []*Dir{},
			AddPath:  "/readme",
			OrigPath: "/readme",
			ExpDirs:  []*Dir{&Dir{Icon: "note", IconAlt: "note", Path: "/readme", Text: "Readme"}},
			M:        "First dir should be added.",
		},
		{
			CurDirs: []*Dir{
				&Dir{Icon: "keyboard_arrow_up", IconAlt: "keyboard_arrow_down", Path: "/hello", Text: "Hello"},
			},
			AddPath:  "/hello/world",
			OrigPath: "/hello/world",
			ExpDirs: []*Dir{
				&Dir{Children: []*Dir{
					&Dir{Icon: "note", IconAlt: "note", Path: "/hello/world", Text: "World"},
				},
					Icon:    "keyboard_arrow_up",
					IconAlt: "keyboard_arrow_down",
					Path:    "/hello",
					Text:    "Hello"},
			},
			M: "Adding a subdir should push the node as a child.",
		},
	}

	for _, a := range x {
		actDirs := addToDir(a.CurDirs, a.AddPath, a.OrigPath)
		assert.Equal(a.ExpDirs, actDirs, a.M)
	}
}

type pageStruct struct {
	ExpPage *Page
	ExpErr  error
	InpPag  string
	InpDir  string
	M       string
}

func TestBuildPage(t *testing.T) {
	assert := assert.New(t)
	x := []pageStruct{
		{
			ExpPage: &Page{Name: "root", Content: "<h1>root</h1>\n"},
			ExpErr:  nil,
			InpPag:  "/root",
			InpDir:  getTestMarkdownDir() + "root.md",
			M:       "Valid file should be processed",
		},
		{
			ExpPage: nil,
			ExpErr: fmt.Errorf(
				"unable to load markdown file: open %s: no such file or directory",
				getTestMarkdownDir()+"hello.md",
			),
			InpPag: "/hello",
			InpDir: getTestMarkdownDir() + "hello.md",
			M:      "Invalid file should error",
		},
	}

	for _, a := range x {
		actPage, actErr := buildPage(a.InpPag, a.InpDir)
		assert.Equal(a.ExpPage, actPage, a.M)
		assert.Equal(a.ExpErr, actErr, a.M)
	}
}

type dirStruct struct {
	ExpDir *Dir
	ExpSuc bool
	InpDir []*Dir
	InpStr string
	M      string
}

func TestDirHasText(t *testing.T) {
	assert := assert.New(t)
	x := []dirStruct{
		{
			ExpDir: nil,
			ExpSuc: false,
			InpDir: []*Dir{},
			InpStr: "",
			M:      "Empty list should fail",
		},
		{
			ExpDir: &Dir{Text: "monkey"},
			ExpSuc: true,
			InpDir: []*Dir{&Dir{Text: "monkey"}, &Dir{Text: "dishwasher"}},
			InpStr: "monkey",
			M:      "Matching dir should be returned",
		},
		{
			ExpDir: nil,
			ExpSuc: false,
			InpDir: []*Dir{&Dir{Text: "monkey"}, &Dir{Text: "dishwasher"}},
			InpStr: "purple",
			M:      "Non-matching text should fail",
		},
	}

	for _, a := range x {
		actDir, actSuc := dirHasText(a.InpDir, a.InpStr)
		assert.Equal(a.ExpDir, actDir, a.M)
		assert.Equal(a.ExpSuc, actSuc, a.M)
	}
}

type tokeniseStruct struct {
	Exp []string
	Inp string
	M   string
}

func TestTokenise(t *testing.T) {
	assert := assert.New(t)
	x := []tokeniseStruct{
		{
			Exp: []string{"a", "b", "c"},
			Inp: "a/b/c",
			M:   "basic string should split",
		},
		{
			Exp: []string{"purple", "monkey", "dishwasher"},
			Inp: "purple/monkey/dishwasher/",
			M:   "trailing slash should be removed",
		},
		{
			Exp: []string{"readme"},
			Inp: "/readme",
			M:   "slash prefix should be removed",
		},
	}

	for _, a := range x {
		act := tokenise(a.Inp)
		assert.Equal(a.Exp, act, a.M)
	}
}

func getTestMarkdownDir() string {
	wd, _ := os.Getwd()
	return wd + "/testdata/"
}
