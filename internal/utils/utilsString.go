package utils

import (
	"regexp"
	"strings"
)

// NormalizeSongName нормализует название песни, убирая лишние пробелы перед и после тире и запятой.
func NormalizeSongName(songName string) string {
	// Заменяем неразрывные пробелы на обычные пробелы
	songName = strings.ReplaceAll(songName, "\u00A0", " ")

	// Убираем пробелы перед и после тире с помощью регулярных выражений
	reDash := regexp.MustCompile(`\s*-\s*`)
	songName = reDash.ReplaceAllString(songName, "-")

	// Убираем пробелы перед и после запятой
	reComma := regexp.MustCompile(`\s*,\s*`)
	songName = reComma.ReplaceAllString(songName, ",")

	// Удаляем лишние пробелы между словами (например, двойные пробелы)
	songName = strings.Join(strings.Fields(songName), " ")

	// Убираем пробелы в начале и конце строки
	songName = strings.TrimSpace(songName)

	return songName
}
