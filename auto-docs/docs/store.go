// Package docs provides storage and retrieval for a series of data
// points based around a tree structure of page content.
//
// The content loaded within this tree has no bearing on the source
// location, these 2 concepts should be kept separate (there may be
// loading of data from separate sources and combined into a single
// cohesive interface).
package docs

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	autodocs "github.com/cloudcloud/auto-docs"
	"gitlab.com/golang-commonmark/markdown"
)

var (
	// S is a singleton instance of the current docs storage.
	S = &Store{}
)

func init() {
	S.Dirs = []*Dir{}
	S.Pages = make(map[string]*autodocs.Page, 0)
}

// Dir gives a holder of further nodes.
type Dir struct {
	Children []*Dir `json:"children"`
	Icon     string `json:"icon"`
	IconAlt  string `json:"icon-alt"`
	Model    bool   `json:"model"`
	Path     string `json:"path"`
	Text     string `json:"text"`
}

// Store captures the doc file references and content.
type Store struct {
	// Dirs tracks the structure of pages under their paths.
	Dirs []*Dir `json:"pages"`

	// Pages captures the content for a full path page.
	Pages map[string]*autodocs.Page `json:"-"`

	// path is the base that this store is defined for.
	path string
}

// UpdateFromPath will accept a base path location and walk the
// directory structure to find appropriate files to be pulled in
// to memory for serving.
func (s *Store) UpdateFromPath(p string) {
	s.path = p
	filepath.Walk(p, s.walker)
}

// walker is the handler method for directory traversal.
func (s *Store) walker(path string, i os.FileInfo, err error) error {
	if err != nil {
		log.Println("unable to read path:", path)
		return err
	}

	// skip our own processing of folders and irrelevant files
	if i.IsDir() || !strings.HasSuffix(i.Name(), ".md") {
		return nil
	}

	x := strings.TrimPrefix(strings.TrimSuffix(strings.ToLower(path), ".md"), s.path)
	p, err := buildPage(x, path)
	if err == nil {
		s.Pages[x] = p
		s.Dirs = addToDir(s.Dirs, x, x)
	}

	return nil
}

// addToDir will push a path into the slice of Dir
// entries, working through entries in the slice to
// prevent duplication of entries.
func addToDir(d []*Dir, p, o string) []*Dir {
	b := tokenise(p)
	dir, yes := dirHasText(d, strings.Title(b[0]))

	if !yes {
		// create a dir and then descend
		dir = &Dir{
			Text: strings.Title(b[0]),
		}

		if len(b) > 1 {
			// have children
			dir.Icon = "keyboard_arrow_up"
			dir.IconAlt = "keyboard_arrow_down"
			dir.Children = []*Dir{}
		} else {
			// no children
			dir.Icon = "note"
			dir.IconAlt = "note"
			dir.Path = o
		}

		d = append(d, dir)
	}

	// dir is the node to descend into now
	if len(b) > 1 {
		// recurse using b[1:]
		dir.Children = addToDir(
			dir.Children,
			strings.Join(b[1:], string(os.PathSeparator)),
			o,
		)
	}

	return d
}

// buildPage
func buildPage(p, d string) (*autodocs.Page, error) {
	b := tokenise(p)
	m := markdown.New(markdown.XHTMLOutput(true))
	f, err := ioutil.ReadFile(d)
	if err != nil {
		return nil, fmt.Errorf("unable to load markdown file: %s", err)
	}

	return &autodocs.Page{
		Name:    b[len(b)-1],
		Content: m.RenderToString(f),
	}, nil
}

// dirHasText will look for an existing Dir in the slice
// that has the requested Text value.
func dirHasText(d []*Dir, t string) (*Dir, bool) {
	for _, x := range d {
		if t == x.Text {
			return x, true
		}
	}

	return nil, false
}

// tokenise will prepare the path for node creation
func tokenise(p string) []string {
	return strings.Split(
		strings.Trim(
			p,
			string(os.PathSeparator),
		),
		string(os.PathSeparator),
	)
}
