package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Открываем файл input.txt для чтения
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Printf("Ошибка при открытии файла: %v\n", err)
		return
	}
	defer file.Close() // Закрываем файл после завершения работы

	// Создаем новый сканер для чтения из файла
	scanner := bufio.NewScanner(file)

	var textBuilder strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		textBuilder.WriteString(line + "\n") // Добавляем строку к тексту с переводом строки
	}

	// Проверяем на ошибки при чтении файла
	if err := scanner.Err(); err != nil {
		fmt.Printf("Ошибка при чтении файла: %v\n", err)
		return
	}

	// Получаем текст из Builder
	text := textBuilder.String()

	// Разбиваем текст на куплеты, используя двойной перевод строки
	verses := strings.Split(text, "\n\n")

	// Формируем массив строк для вывода
	var formatted []string
	for _, verse := range verses {
		// Удаляем лишние пробелы и добавляем к массиву, если не пусто
		trimmed := strings.TrimSpace(verse)
		if trimmed != "" {
			// Удаляем все переводы строк внутри куплета
			noLineBreaks := strings.ReplaceAll(trimmed, "\n", " ")
			formatted = append(formatted, fmt.Sprintf("\"%s\"", noLineBreaks))
		}
	}

	// Объединяем в нужный формат с запятыми
	result := strings.Join(formatted, ", ")

	// Определяем имя файла для сохранения
	fileName := "output.txt"

	// Создаем файл или перезаписываем его, если он уже существует
	file, err = os.Create(fileName)
	if err != nil {
		fmt.Printf("Ошибка при создании файла: %v\n", err)
		return
	}
	defer file.Close() // Закрываем файл после завершения работы

	// Записываем результат в файл
	_, err = file.WriteString(result)
	if err != nil {
		fmt.Printf("Ошибка при записи в файл: %v\n", err)
		return
	}

	fmt.Printf("Результат сохранён в файл: %s\n", fileName)
}
