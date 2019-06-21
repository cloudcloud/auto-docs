package server

import (
	"bytes"
	"net/http"

	"github.com/cloudcloud/auto-docs/auto-docs/docs"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/gin"
)

// handleFiles will add the handling methods for each of
// the available assets.
func handleFiles(e *gin.Engine) {
	e.StaticFS("/_assets",
		&assetfs.AssetFS{
			Asset:     Asset,
			AssetDir:  AssetDir,
			AssetInfo: AssetInfo,
			Prefix:    "/_assets",
		},
	)
}

// health checks the state of auto-docs and provides a
// response based on this.
func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// page will retrieve the data for a single page.
func page(c *gin.Context) {
	if p, ok := docs.S.Pages[c.Param("path")]; !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Path not found"})
	} else {
		c.JSON(http.StatusOK, p)
	}
}

// pages will provide a full list of the tree structure
// of pages currently available.
func pages(c *gin.Context) {
	c.JSON(http.StatusOK, docs.S)
}

// root will serve the base shell, and then filter out
// generated paths after.
func root(c *gin.Context) {
	r := bytes.NewReader(MustAsset("index.html"))
	c.DataFromReader(
		http.StatusOK,
		int64(len(MustAssetString("index.html"))),
		"text/html",
		r,
		map[string]string{},
	)
}
