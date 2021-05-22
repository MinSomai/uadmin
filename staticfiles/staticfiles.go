package staticfiles

import (
	"fmt"
	abtestmodel "github.com/uadmin/uadmin/blueprint/abtest/models"
	"net/http"
	"os"
	"strings"
	"time"
)

func containsDotDot(v string) bool {
	if !strings.Contains(v, "..") {
		return false
	}
	for _, ent := range strings.FieldsFunc(v, func(r rune) bool { return r == '/' || r == '\\' }) {
		if ent == ".." {
			return true
		}
	}
	return false
}

// StaticHandler is a function that serves static files
func StaticHandler(w http.ResponseWriter, r *http.Request) {
	if containsDotDot(r.URL.Path) {
		w.WriteHeader(404)
		return
	}
	var modTime time.Time
	ab := false
	var midnightDelta int
	for k := range abtestmodel.StaticABTests {
		if k == r.URL.Path && len(abtestmodel.StaticABTests[k]) != 0 {
			index := abtestmodel.GetABT(r) % len(abtestmodel.StaticABTests[k])
			r.URL.Path = abtestmodel.StaticABTests[k][index].v

			// Calculate number of seconds until midnight if no calculated yet
			if midnightDelta == 0 {
				midnight := time.Now()
				midnight = time.Date(midnight.Year(), midnight.Month(), midnight.Day(), 0, 0, 0, 0, midnight.Location())
				midnightDelta = int(time.Until(midnight).Seconds())
			}
			// Add a header to expire the satic content at midnigh
			w.Header().Add("Cache-Control", "private, max-age="+fmt.Sprint(midnightDelta))

			ab = true

			go func(index int) {
				abtestmodel.AbTestsMutex.Lock()
				t := abtestmodel.StaticABTests[k]
				t[index].imp++
				abtestmodel.StaticABTests[k] = t
				abtestmodel.AbTestsMutex.Unlock()
			}(index)
			break
		}
	}

	f, err := os.Open("." + r.URL.Path)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	if !ab {
		stat, err := os.Stat("." + r.URL.Path)
		if err != nil || stat.IsDir() {
			w.WriteHeader(404)
			return
		}
		modTime = stat.ModTime()
		w.Header().Add("Cache-Control", "private, max-age=3600")
	} else {
		modTime = time.Now()
	}

	http.ServeContent(w, r, "."+r.URL.Path, modTime, f)
}