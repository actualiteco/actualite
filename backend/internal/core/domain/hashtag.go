package domain

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

const (
	MinTagLength = 2
	MaxTagLength = 50
)

var (
	hashtagRe = regexp.MustCompile(`#([\p{L}\p{N}_]+)`)
	urlRe     = regexp.MustCompile(`(?i)https?://\S+`)
)

type Tag string

func (t Tag) String() string { return string(t) }

func NewTag(raw string) (Tag, error) {
	s := normalize(raw)
	if s == "" {
		return "", fmt.Errorf("hashtag is empty: %w", ErrInvalidRequest)
	}

	rl := utf8.RuneCountInString(s)
	if rl < MinTagLength {
		return "", fmt.Errorf("hashtag=%s is too short: %w", s, ErrInvalidRequest)
	}

	if rl > MaxTagLength {
		return "", fmt.Errorf("hashtag=%s is too long: %w", s, ErrInvalidRequest)
	}

	if isNumericOnly(s) {
		return "", fmt.Errorf("hashtag=%s cannot consist only of numbers: %w", s, ErrInvalidRequest)
	}

	return Tag(s), nil
}

func normalize(s string) string {
	return strings.ToLower(norm.NFKC.String(strings.TrimPrefix(s, "#")))
}

func isNumericOnly(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func ExtractTags(text string) []Tag {
	if text == "" {
		return nil
	}

	cleaned := urlRe.ReplaceAllString(text, " ")
	matches := hashtagRe.FindAllStringSubmatchIndex(cleaned, -1)
	if len(matches) == 0 {
		return nil
	}

	seen := make(map[Tag]struct{}, len(matches))
	out := make([]Tag, 0, len(matches))

	for _, m := range matches {
		hashPos, start, end := m[0], m[2], m[3]

		if hashPos > 0 {
			prev, _ := utf8.DecodeLastRuneInString(cleaned[:hashPos])
			if prev == '#' || unicode.IsLetter(prev) || unicode.IsNumber(prev) || prev == '_' {
				continue
			}
		}

		tag, err := NewTag(cleaned[start:end])
		if err != nil {
			continue
		}

		if _, ok := seen[tag]; ok {
			continue
		}

		seen[tag] = struct{}{}
		out = append(out, tag)
	}

	return out
}

func StringsToTags(strs []string) ([]Tag, error) {
	tags := make([]Tag, 0, len(strs))

	for _, str := range strs {
		tag, err := NewTag(str)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func TagsToStrings(tags []Tag) []string {
	strs := make([]string, 0, len(tags))

	for _, tag := range tags {
		strs = append(strs, tag.String())
	}

	return strs
}
