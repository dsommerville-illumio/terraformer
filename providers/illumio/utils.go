package illumio

import (
	"strings"

	"github.com/brian1917/illumioapi"
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

func convertLabelsToReferenceSlice(labels []*illumioapi.Label) []Reference {
	refs := []Reference{}
	for _, label := range labels {
		refs = append(refs, Reference{Href: label.Href})
	}
	return refs
}
