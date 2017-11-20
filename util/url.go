/*
package util provides utilities
*/
package util

import (
	"net/url"
	"strings"
)

// StringToUrls parse comma saperated url string to []url.URL
func StringToUrls(raw string) ([]url.URL, error) {
	raws := strings.Split(raw, ",")
	res := make([]url.URL, 0)
	for _, s := range raws {
		if len(s) == 0 {
			continue
		}
		u, err := url.Parse(s)
		if err != nil {
			return nil, err
		}
		res = append(res, *u)
	}
	return res, nil
}

// UrlsToStrings parse []url.URL to []string, useful for get endpoints
func UrlsToStrings(urls []url.URL) []string {
	ret := make([]string, 0)
	for _, u := range urls {
		ret = append(ret, u.String())
	}
	return ret
}

// RemoveScheme removes http or https scheme from address, e.g. http://1.1.1.1 to 1.1.1.1
func RemoveScheme(addr string) string {
	ret := strings.TrimPrefix(addr, "http://")
	ret = strings.TrimPrefix(ret, "https://")
	return ret
}
