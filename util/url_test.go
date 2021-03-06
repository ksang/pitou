package util

import (
	"reflect"
	"strings"
	"testing"
)

func TestStringToUrls(t *testing.T) {
	var tests = []struct {
		s   string
		err bool
	}{
		{
			"http://127.0.0.1,",
			false,
		},
		{
			"http://127.0.0.1,http://0.0.0.0:80/",
			false,
		},
		{
			"",
			false,
		},
	}

	for caseid, c := range tests {
		res, err := StringToUrls(c.s)
		if err != nil {
			if c.err {
				t.Logf("case #%d, err: %v", caseid+1, err)
			} else {
				t.Errorf("case #%d, err: %v", caseid+1, err)
			}
		}
		if c.err {
			t.Errorf("case #%d, no error returned, expecting error", caseid+1)
		}
		t.Logf("Result: %v", res)
	}
}

func TestUrlsToStrings(t *testing.T) {
	var tests = []struct {
		s string
		e []string
	}{
		{
			"http://127.0.0.1,",
			[]string{"http://127.0.0.1"},
		},
		{
			"http://127.0.0.1,http://0.0.0.0:80",
			[]string{"http://127.0.0.1", "http://0.0.0.0:80"},
		},
		{
			"",
			[]string{},
		},
		{
			"http://localhost:2379",
			[]string{"http://localhost:2379"},
		},
	}

	for caseid, c := range tests {
		res, err := StringToUrls(c.s)
		if err != nil {
			t.Errorf("case #%d, failed to parse case string err: %v", caseid+1, err)
		}
		r := UrlsToStrings(res)
		if !reflect.DeepEqual(r, c.e) {
			t.Errorf("case #%d, result incorrect: %v, expected: %v", caseid+1, r, c.e)

		}
		t.Logf("Result: %v", r)
	}
}

func TestRemoveScheme(t *testing.T) {
	var tests = []struct {
		s string
		e string
	}{
		{
			"http://127.0.0.1:1234",
			"127.0.0.1:1234",
		},
		{
			"https://127.0.0.1:1234",
			"127.0.0.1:1234",
		},
		{
			"http://127.0.0.1",
			"127.0.0.1",
		},
	}

	for caseid, c := range tests {
		res := RemoveScheme(c.s)
		if strings.Compare(res, c.e) != 0 {
			t.Errorf("case #%d, expected: %s, actual: %s", caseid+1, c.e, res)
		}
	}
}
