package main

import (
	"strings"

	pluralize "github.com/gertd/go-pluralize"
)

func getWordForms(guess string) []string {
	var allForms []string
	allForms = append(allForms, guess)
	pluralize := pluralize.NewClient()
	if pluralize.IsPlural(guess) {
		allForms = append(allForms, pluralize.Singular(guess))
	} else if pluralize.IsSingular(guess) {
		allForms = append(allForms, pluralize.Plural(guess))
	}

	return allForms
}
func checkEitherIsSubString(currentTag string, acceptableAnswer string) bool {
	if strings.Contains(acceptableAnswer, currentTag) || (strings.Contains(currentTag, acceptableAnswer) && len(acceptableAnswer) >= 4) {
		return true
	}
	return false

}
