package api

import (
	"fmt"
	"net/http"
	"strings"

	cxcore "github.com/skycoin/cx/cx"

	"net/http/httputil"
)

// API represents an HTTP API.
type API struct {
	root string
	pg   *cxcore.CXProgram
}

// NewAPI returns a new API given a CX Program.
func NewAPI(root string, pg *cxcore.CXProgram) *API {
	if root == "" {
		root = "/"
	}
	return &API{root: root, pg: pg}
}

// ServeHTTP implements http.Handler
func (a *API) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	mux := http.NewServeMux()
	mux.HandleFunc(a.root+"/meta", ProgramMeta(a.pg))
	mux.HandleFunc(a.root+"/packages", PackagesOfProgram(a.pg))
	mux.HandleFunc(a.root+"/packages/", func(w http.ResponseWriter, req *http.Request) {
		split := strings.Split(req.URL.EscapedPath(), "/")
		base := split[len(split)-1]
		ExportedSymbolsOfPackage(a.pg, base)(w, req)
	})
	mux.ServeHTTP(w, req)
}

// ProgramMeta returns the program meta data.
func ProgramMeta(pg *cxcore.CXProgram) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		resp := extractProgramMeta(pg)
		httputil.WriteJSON(w, req, http.StatusOK, resp)
	}
}

// PackagesOfProgram returns an array of package names of a given program.
func PackagesOfProgram(pg *cxcore.CXProgram) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		pkgNames := make([]string, 0, len(pg.Packages))
		for _, pkg := range pg.Packages {
			pkgNames = append(pkgNames, pkg.Name)
		}

		httputil.WriteJSON(w, req, http.StatusOK, pkgNames)
	}
}

// ExportedSymbolsOfPackage returns exported symbols of a given package.
func ExportedSymbolsOfPackage(pg *cxcore.CXProgram, pkgName string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		for _, pkg := range pg.Packages {
			if pkg.Name != pkgName {
				continue
			}

			resp := extractExportedSymbols(pkg)
			httputil.WriteJSON(w, req, http.StatusOK, resp)
			return
		}

		err := fmt.Errorf("package '%s' is not found in program '%s'", pkgName, pg.Path)
		httputil.WriteJSON(w, req, http.StatusNotFound, err)
	}
}
