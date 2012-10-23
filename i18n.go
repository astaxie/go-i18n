package i18n

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrInvalidLocale is returned if a specified locale is invalid.
	ErrInvalidLocale = errors.New("invalid locale")
)

type IL struct {
	defaultLocale          string
	dirPath                string
	locales                []string
	locale2SourceKey2Trans map[string]map[string]string
}

func NewIL(translationsDirPath, locale string) (*IL, error) {
	locales := localesChainForLocale(locale)
	if len(locales) == 0 {
		return nil, ErrInvalidLocale
	}

	il := &Dict{
		defaultLocale:          locales[0],
		dirPath:                translationsDirPath,
		locales:                locales,
		locale2SourceKey2Trans: make(map[string]map[string]string),
	}

	if err := il.loadDefaultTranslations(locales); err != nil {
		return nil, err
	}

	if err := il.loadTranslations(translationsDirPath); err != nil {
		return nil, err
	}

	return il, nil
}

func (il *IL) loadDefaultTranslations(locales []string) {
	
}

// DirPath returns the translations directory path of the international.
func (il *IL) DirPath() string {
	return il.dirPath
}

// Locale returns the locale of the international.
func (il *IL) Locale() string {
	return il.locales[0]
}

// Translation returns a translation of the source string to the locale of the 
// Dict or the source string, if no matching translation was found.
//
// Provided context arguments are used for disambiguation.
func (il *IL) Translation(source string, context ...string) string {
	for _, locale := range il.locales {
		if sourceKey2Trans, ok := il.locale2SourceKey2Trans[locale]; ok {
			if trans, ok := sourceKey2Trans[sourceKey(source, context)]; ok {
				return trans
			}
		}
	}

	return source
}

func (il *IL) loadTranslation(reader io.Reader, locale string) error {
	var trf trfile

	if err := json.NewDecoder(reader).Decode(&trf); err != nil {
		return err
	}

	sourceKey2Trans, ok := il.locale2SourceKey2Trans[locale]
	if !ok {
		sourceKey2Trans = make(map[string]string)

		il.locale2SourceKey2Trans[locale] = sourceKey2Trans
	}

	for _, m := range trf.Messages {
		if m.Translation != "" {
			sourceKey2Trans[sourceKey(m.Source, m.Context)] = m.Translation
		}
	}

	return nil
}

func (il *IL) loadTranslations(dirPath string) error {
	dir, err := os.Open(dirPath)
	if err != nil {
		return err
	}
	defer dir.Close()

	names, err := dir.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		fullPath := path.Join(dirPath, name)

		fi, err := os.Stat(fullPath)
		if err != nil {
			return err
		}

		if fi.IsDir() {
			if err := il.loadTranslations(fullPath); err != nil {
				return err
			}
		} else if locale := il.matchingLocaleFromFileName(name); locale != "" {
			file, err := os.Open(fullPath)
			if err != nil {
				return err
			}
			defer file.Close()

			if err := il.loadTranslation(file, locale); err != nil {
				return err
			}
		}
	}

	return nil
}

func (il *IL) matchingLocaleFromFileName(name string) string {
	for _, locale := range il.locales {
		if strings.HasSuffix(name, fmt.Sprintf("-%s.tr", locale)) {
			return locale
		}
	}

	return ""
}

func sourceKey(source string, context []string) string {
	if len(context) == 0 {
		return source
	}

	return fmt.Sprintf("__%s__%s__", source, strings.Join(context, "__"))
}

func localesChainForLocale(locale string) []string {
	parts := strings.Split(locale, "_")
	if len(parts) > 2 {
		return nil
	}

	if len(parts[0]) != 2 {
		return nil
	}

	for _, r := range parts[0] {
		if r < rune('a') || r > rune('z') {
			return nil
		}
	}

	if len(parts) == 1 {
		return []string{parts[0]}
	}

	if len(parts[1]) < 2 || len(parts[1]) > 3 {
		return nil
	}

	for _, r := range parts[1] {
		if r < rune('A') || r > rune('Z') {
			return nil
		}
	}

	return []string{locale, parts[0]}
}
