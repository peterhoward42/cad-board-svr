package server

import (
	"html/template"
	"net/http"
	"appengine"
	"appengine/memcache"
	"appengine/log"
)

func init() {
	http.HandleFunc("/", mainPageHandler)
	http.HandleFunc("/mouse", mouseHandler)
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	parsed_template, err := template.ParseFiles("static/template/index.html")
	// how to use gui name in template exec fn
	parsed_template.Execute(w, gui_data())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func mouseHandler(w http.ResponseWriter, r *http.Request) {
	// extract the mouse x coordinate sent
	x := r.FormValue("xcoord")
	ctx := appengine.NewContext(r)
	ctx.Infof("request contains x coord string: %v", x)
	storeXCoordInMemcache(ctx, x);
}

func storeXCoordInMemcache(ctx appengine.Context, x string) void {
	newMemcacheItem := &memcache.Item{Key: "mouse_x", Value: []byte(x)}
	// Add the item to the cache (if not already present)
	err := memcache.Add(ctx, newMemcacheItem)
	// Ignore error that is already present in cache, but react to other errors
	if err != memcache.ErrNotStored {
		ctx.Infof("error adding item: %v", err)
	}
	// Now overwrite the cached item regardless
	err = memcache.Set(ctx, newMemcacheItem)
	if err {
		ctx.Errorf(ctx, "error setting item: %v", err)
	}
}

// A data structure for the model part of the example GUI's model-view pattern.
type GuiDataModel struct {
	Title       string
	Unwatch     int
	Star        int
	Fork        int
	Commits     int
	Branch      int
	Release     int
	Contributor int
	RowsInTable []TableRow
}

// A sub-model to the GuiDataModel that models a single row in the table
// displayed in the GUI.
type TableRow struct {
	File    string
	Comment string
	Ago     string
	Icon    string
}

// Provides an illustrative, hard-coded instance of a GuiDataModel.
func gui_data() *GuiDataModel {
	gui_data := &GuiDataModel{
		Title:       "Golang Standalone GUI Example",
		Unwatch:     3,
		Star:        0,
		Fork:        2,
		Commits:     31,
		Release:     1,
		Contributor: 1,
		RowsInTable: []TableRow{},
	}
	gui_data.RowsInTable = append(gui_data.RowsInTable,
		TableRow{"do_this.go", "Initial commit", "1 month ago", "file"})
	gui_data.RowsInTable = append(gui_data.RowsInTable,
		TableRow{"do_that.go", "Initial commit", "1 month ago", "file"})
	gui_data.RowsInTable = append(gui_data.RowsInTable,
		TableRow{"index.go", "Initial commit", "1 month ago", "file"})
	gui_data.RowsInTable = append(gui_data.RowsInTable,
		TableRow{"resources", "Initial commit", "2 months ago", "folder-open"})
	gui_data.RowsInTable = append(gui_data.RowsInTable,
		TableRow{"docs", "Initial commit", "2 months ago", "folder-open"})
	return gui_data
}