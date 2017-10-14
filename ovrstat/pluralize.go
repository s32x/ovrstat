package ovrstat

import (
	"sort"
	"strings"
	"unicode"

	"github.com/jinzhu/inflection"
)

var keywords = []string{
	"kill",
	"multikill",
	"death",
	"generator",
	"shield",
	"enemy",
	"turret",
	"hit",
	"pad",
	"blow",
	"assist",
	"elimination",
	"card",
	"dragonblade",
	"player",
	"bomb",
}

func transformKey(str string) string {
	split := splitKeywords(str)

	sort.Sort(split)

	var lowerKeyPart string

	for _, keyPart := range split {
		lowerKeyPart = strings.ToLower(keyPart)

		for _, keyword := range keywords {
			if keyword == lowerKeyPart {
				replacement := inflection.Plural(keyword)

				if unicode.IsUpper(rune(keyPart[0])) {
					replacement = strings.ToUpper(replacement[0:1]) + replacement[1:]
				}

				str = strings.Replace(str, keyPart, replacement, 1)

				return str
			}
		}
	}

	if str == "allDamageDone" {
		return "damageDone"
	}

	return str
}

type Keywords []string

func (s Keywords) Len() int {
	return len(s)
}

func (s Keywords) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Keywords) Less(i, j int) bool {
	return sliceIndexOf(keywords, s[i]) < sliceIndexOf(keywords, s[j])
}

func splitKeywords(str string) Keywords {
	var words Keywords
	l := 0
	for s := str; s != ""; s = s[l:] {
		l = strings.IndexFunc(s[1:], unicode.IsUpper) + 1
		if l <= 0 {
			l = len(s)
		}
		words = append(words, s[:l])
	}
	return words
}

func sliceIndexOf(s []string, str string) int {
	for i, val := range s {
		if strings.ToLower(val) == strings.ToLower(str) {
			return i
		}
	}
	return len(s)
}
