package illumio

import (
	"strings"
)

const (
	DRAFT string = "draft"
)

type Reference struct {
	Href string `json:"href,omitempty"`
}

func stripIdFromHref(href string) string {
	if href == "" || !strings.Contains(href, "/") {
		return href
	}
	return strings.Trim(href[strings.LastIndex(href, "/"):], "/")
}
