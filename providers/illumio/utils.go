package illumio

import "strings"

func stripIdFromHref(href string) string {
	if href == "" || !strings.Contains(href, "/") {
		return href
	}
	return strings.Trim(href[strings.LastIndex(href, "/"):], "/")
}

func stripServiceNameQualifiers(serviceName string) string {
	if !strings.Contains(serviceName, "[") {
		return serviceName
	}
	return serviceName[:strings.Index(serviceName, "[")]
}
