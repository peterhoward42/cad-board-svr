package server

import (
	"html/template"
	"net/http"
	"appengine"
	"appengine/memcache"
	"fmt"
)

func init() {
	http.HandleFunc("/", mainPageHandler)
	http.HandleFunc("/mouse", mouseHandler)
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	parsed_template, err := template.ParseFiles("static/template/index.html")
	parsed_template.Execute(w, gui_data())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func mouseHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the mouse coordinates from the payload.
	x := r.FormValue("X")
	y := r.FormValue("Y")
	// Save them in memcache
	ctx := appengine.NewContext(r)
	storeCoordInMemcache(ctx, x, y)
	// Send reply that all is well.

	// Temporarily todo - retrieve it again too
	position := retrieveCoordFromMemcache(ctx);
	fmt.Fprintf(w, "ok, retrieved x: %s, y: %s", position.X, position.Y);
}

// Todo - should have proper type for the stored item with single definition.
// Todo - should have single definition of the key string.
// Todo - need error handling on memcache api calls.
func storeCoordInMemcache(ctx appengine.Context, x string, y string) {
	var positionToStore struct{X string; Y string}
	positionToStore.X = x
	positionToStore.Y = y
	item := &memcache.Item{Key: "mousePosition", Object: positionToStore}
	// We have to attempt to add it before setting it, to cope with first time,
	// and with it having been evicted. We tolerate not-stored error because this will
	// be the general case when it is present.
	if err := memcache.Gob.Add(ctx, item); err != memcache.ErrNotStored {
		ctx.Errorf("error adding item: %v", err)
	}
	if err := memcache.Gob.Set(ctx, item); err != nil {
		ctx.Errorf("error setting item: %v", err)
	}
}

func retrieveCoordFromMemcache(ctx appengine.Context) struct{X string; Y string} {
	var positionRetrieved struct{X string; Y string}
	_, err := memcache.Gob.Get(ctx, "mousePosition", &positionRetrieved);
	if (err != nil) {
		ctx.Infof("error from gob set: %s", err)
	}
	return positionRetrieved;
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