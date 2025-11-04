package helperpackage

import "math/rand"

// GetRandomPhrase возвращает случайную фразу для комментария
func GetRandomPhrase() string {
	phrases := []string{
		"Требуется радиоуглеродный анализ",
		"Необходимо определить возраст образца",
		"Анализ для археологических исследований",
		"Датаирование для научной работы",
		"Определение календарного возраста",
	}
	return phrases[rand.Intn(len(phrases))]
}

// GetRandomSampleDescription возвращает случайное описание образца
func GetRandomSampleDescription() string {
	descriptions := []string{
		"Фрагмент археологической находки",
		"Образец с раскопок",
		"Лабораторный образец",
		"Полевой образец",
		"Исследуемый материал",
		"Артефакт с исторического объекта",
	}
	return descriptions[rand.Intn(len(descriptions))]
}
