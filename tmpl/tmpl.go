package tmpl

import (
	"embed"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"text/template"
	"time"
)

//go:embed tmpls/*
var fs embed.FS
var embedTemplates map[string]*template.Template
var customTemplates map[string]*template.Template
var funcMap template.FuncMap

func init() {
	// func
	funcMap = template.FuncMap{
		"date": func(dt time.Time, zone string) string {
			loc, err := time.LoadLocation(zone)
			if err != nil {
				return err.Error()
			}
			dt = dt.In(loc)
			return dt.Format("2006-01-02 15:04:05")
		},
		"isNonZeroDate": func(dt time.Time) bool {
			return !(dt == time.Time{})
		},
		"in": func(m map[string]string, key string) bool {
			_, ok := m[key]
			return ok
		},
		"toUpper": strings.ToUpper,
		"toLink": func(s string) string {
			return fmt.Sprintf("[%s](%s)", s, s)
		},
		"displayKV": func(k, v string) string {
			_, err := url.ParseRequestURI(v)
			if err != nil {
				return fmt.Sprintf("%s:%s", k, v)
			}
			return fmt.Sprintf("[%s](%s)", k, v)
		},
		"contains": strings.Contains,
	}

	// embed
	dir, err := fs.ReadDir("tmpls")
	if err != nil {
		panic(err)
	}

	embedTemplates = make(map[string]*template.Template)
	for _, entry := range dir {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if !strings.HasSuffix(filename, ".tmpl") {
			continue
		}

		t, err := template.New(filename).Funcs(funcMap).ParseFS(fs, "tmpls/"+filename)
		if err != nil {
			panic(err)
		}

		embedTemplates[t.Name()] = t
	}

	// custom
	customTemplates = make(map[string]*template.Template)
}

func GetEmbedTemplate(filename string) (*template.Template, error) {
	if t, ok := embedTemplates[filename]; ok {
		return t, nil
	}

	return nil, errors.New("template not found")
}

func GetCustomTemplate(filename string) (*template.Template, error) {
	if t, ok := customTemplates[filename]; ok {
		return t, nil
	}

	t, err := template.New(filename).Funcs(funcMap).ParseFiles(filename)
	if err != nil {
		return nil, err
	}
	customTemplates[filename] = t

	return t, nil
}
