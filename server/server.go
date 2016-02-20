package server

import (
	"encoding/json"
	"html/template"
	"net/http"
	"appengine"
	"appengine/memcache"
	"fmt"
)

const mouseKey string = "MOUSEKEY"

type xy struct { X, Y string }

func init() {
	http.HandleFunc("/", mainPageHandler)
	http.HandleFunc("/mouseposnupdate", receiveMousePositionUpdate)
	http.HandleFunc("/mouseposnquery", receiveMousePositionQuery)
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	parsed_template, err := template.ParseFiles("static/template/index.html")
	parsed_template.Execute(w, gui_data())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func receiveMousePositionUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	// Extract the mouse coordinates from the payload.
	mousePosn := xy{r.FormValue("X"), r.FormValue("Y")}
	ctx.Errorf("struct literal: %v %v", mousePosn.X, mousePosn.Y)
	// Save them in memcache
	storeCoordInMemcache(ctx, &mousePosn);
	fmt.Fprintf(w, "stored mouse to memcache ok");
}

func receiveMousePositionQuery(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	var posnFromMemcache xy;
	retrieveCoordFromMemcache(ctx, &posnFromMemcache)
	js, err := json.Marshal(posnFromMemcache)
	if err != nil {
		ctx.Errorf("error with json marshal: %v", err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func storeCoordInMemcache(ctx appengine.Context, mousePosn *xy) {
	item := &memcache.Item{Key: mouseKey, Object: *mousePosn}
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

func retrieveCoordFromMemcache(ctx appengine.Context, mousePosn *xy) {
	_, err := memcache.Gob.Get(ctx, mouseKey, mousePosn);
	if (err != nil) {
		ctx.Infof("error from gob set: %s", err)
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