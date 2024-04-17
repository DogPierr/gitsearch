package dataparser

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/DogPierr/gitsearch/gitsearch/internal/config"
)

type lineCommitsFiles struct {
	lines     int
	commit    int
	fileCount int
}

type language struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Extensions []string `json:"extensions"`
}

var (
	info         parseInfo
	languagesMap map[string]*language
)

type parseInfo struct {
	repository   string
	revision     string
	orderby      string
	usecommitter bool
	format       string
	extensions   []string
	languages    []string
	exclude      []string
	restrictto   []string
}

func checkIfLanguageIsValid(lang string) bool {
	if _, ok := languagesMap[lang]; !ok || languagesMap[lang] == nil {
		var l = log.New(os.Stderr, "", 0)
		if lang != "" {
			l.Printf("Unknown language %s", lang)
		}
		return false
	}
	return true
}

func getParseInfo(
	repository string,
	revision string,
	orderby string,
	usecommitter bool,
	format string,
	extensions string,
	languages string,
	exclude string,
	restrictto string,
) parseInfo {
	extensionsParsed := strings.Split(extensions, ",")
	if len(extensionsParsed) == 1 && extensionsParsed[0] == "" {
		extensionsParsed = []string{}
	}
	languagesParsed := make([]string, 0)
	for _, lang := range strings.Split(languages, ",") {
		if checkIfLanguageIsValid(lang) {
			languagesParsed = append(languagesParsed, lang)
		}
	}
	if len(languagesParsed) == 1 && languagesParsed[0] == "" {
		languagesParsed = []string{}
	}
	restricttoParsed := strings.Split(restrictto, ",")
	if len(restricttoParsed) == 1 && restricttoParsed[0] == "" {
		restricttoParsed = []string{}
	}
	excludeParsed := strings.Split(exclude, ",")
	if len(excludeParsed) == 1 && excludeParsed[0] == "" {
		excludeParsed = []string{}
	}
	return parseInfo{
		repository:   repository,
		revision:     revision,
		orderby:      orderby,
		usecommitter: usecommitter,
		format:       format,
		extensions:   extensionsParsed,
		languages:    languagesParsed,
		exclude:      excludeParsed,
		restrictto:   restricttoParsed,
	}
}

func (p parseInfo) isExtensionFine(file string) bool {
	if len(p.extensions) == 0 {
		return true
	}
	for _, ext := range p.extensions {
		if filepath.Ext(file) == ext {
			return true
		}
	}
	return false
}

func (p parseInfo) isGlobFine(file string) bool {
	if len(p.exclude) == 0 {
		return true
	}
	for _, glob := range p.exclude {
		if matched, err := filepath.Match(glob, file); err == nil && matched {
			return false
		}
	}
	return true
}

func (p parseInfo) isLanguageFine(file string) bool {
	if len(p.languages) == 0 {
		return true
	}
	fileExt := filepath.Ext(file)
	for _, lang := range p.languages {
		for _, ext := range languagesMap[lang].Extensions {
			if ext == fileExt {
				return true
			}
		}
	}
	return false
}

func (p parseInfo) isRestrictFine(file string) bool {
	if len(p.restrictto) == 0 {
		return true
	}
	for _, rest := range p.restrictto {
		if matched, err := filepath.Match(rest, file); err == nil && matched {
			return true
		}
	}
	return false
}

func (p parseInfo) isFileFine(file string) bool {
	return p.isExtensionFine(file) &&
		p.isGlobFine(file) &&
		p.isLanguageFine(file) &&
		p.isRestrictFine(file)
}

func handleWrongArguments() error {
	if info.format != "tabular" && info.format != "csv" && info.format != "json" && info.format != "json-lines" {
		return fmt.Errorf("wrong format flag: %s", info.format)
	}
	if info.orderby != "lines" && info.orderby != "commits" && info.orderby != "files" {
		return fmt.Errorf("wrong orderby flag: %s", info.orderby)
	}
	return nil
}

func initParseInfo(
	repository string,
	revision string,
	orderby string,
	usecommitter bool,
	format string,
	extensions string,
	languages string,
	exclude string,
	restrictto string,
) error {
	languagesMap = make(map[string]*language)
	var rawLanguages []language
	err := json.Unmarshal([]byte(config.Config), &rawLanguages)
	if err != nil {
		return err
	}

	for _, lang := range rawLanguages {
		if _, ok := languagesMap[strings.ToLower(lang.Name)]; !ok {
			pLang := lang
			languagesMap[strings.ToLower(lang.Name)] = &pLang
		}
	}
	info = getParseInfo(
		repository,
		revision,
		orderby,
		usecommitter,
		format,
		extensions,
		languages,
		exclude,
		restrictto,
	)

	err = handleWrongArguments()
	if err != nil {
		return err
	}
	return nil
}
