package web

import (
	"fmt"
	"log"
	models "logGen/pkg/dbmodels"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// showAllLogs - renders all logs (no filters)
func showAllLogs(c *gin.Context) {
	entries, err := models.Query(DB, []string{}) // empty filter to get all logs
	if err != nil {
		c.HTML(http.StatusInternalServerError, "result.html", gin.H{"error": err.Error()})
		return
	}

	// Render with empty Form so template doesn't try to checkboxes
	c.HTML(http.StatusOK, "result.html", gin.H{
		"entries": entries,
		"count":   len(entries),
		"Form":    map[string][]string{},
	})
}

// filterLogs - reads form (including multiple checkbox values) and queries DB
func filterLogs(c *gin.Context) {
	// Ensure form is parsed so c.Request.Form is populated.
	// Gin usually parses the form for you when using c.PostForm, but ParseForm ensures the underlying
	// net/http request has Form and PostForm filled when we want to inspect raw Form slices.
	if err := c.Request.ParseForm(); err != nil {
		c.HTML(http.StatusInternalServerError, "result.html", gin.H{"error": "failed to parse form: " + err.Error()})
		return
	}

	// Read all values. We check both "level" and "level[]" variants so it works with either HTML naming style.
	levels := c.Request.Form["level"]
	if len(levels) == 0 {
		levels = c.Request.Form["level[]"]
	}
	components := c.Request.Form["component"]
	if len(components) == 0 {
		components = c.Request.Form["component[]"]
	}
	hosts := c.Request.Form["host"]
	if len(hosts) == 0 {
		hosts = c.Request.Form["host[]"]
	}

	// Single value inputs
	timestamp := c.PostForm("timestamp")
	requestID := c.PostForm("request_id")

	// DEBUG: dump raw form to server logs to verify what client sent
	for k, v := range c.Request.Form {
		log.Printf("FORM KEY=%s VALUES=%v\n", k, v)
	}
	log.Printf("Parsed filters - levels=%v components=%v hosts=%v timestamp=%q requestID=%q\n",
		levels, components, hosts, timestamp, requestID)

	// Build queryParts.
	// OPTION A (used here): create one "key=value" token per selected value
	// Example tokens: ["level=INFO", "level=WARN", "component=worker", "host=web01"]
	var queryParts []string
	for _, lv := range levels {
		lv = strings.TrimSpace(lv)
		if lv != "" {
			queryParts = append(queryParts, fmt.Sprintf("level=%s", lv))
		}
	}
	for _, comp := range components {
		comp = strings.TrimSpace(comp)
		if comp != "" {
			queryParts = append(queryParts, fmt.Sprintf("component=%s", comp))
		}
	}
	for _, h := range hosts {
		h = strings.TrimSpace(h)
		if h != "" {
			queryParts = append(queryParts, fmt.Sprintf("host=%s", h))
		}
	}
	if timestamp != "" {
		queryParts = append(queryParts, fmt.Sprintf("timestamp=%s", timestamp))
	}
	if requestID != "" {
		queryParts = append(queryParts, fmt.Sprintf("request_id=%s", requestID))
	}

	// Execute query using your models.Query which accepts []string tokens
	entries, err := models.Query(DB, queryParts)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "result.html", gin.H{"error": err.Error()})
		return
	}

	// Pass the raw Form map so the template can re-check previously selected checkboxes
	c.HTML(http.StatusOK, "result.html", gin.H{
		"entries": entries,
		"count":   len(entries),
		"Form":    c.Request.Form,
	})
}
