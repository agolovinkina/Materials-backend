package service

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type CalendarDateCalculator struct{}

func NewCalendarDateCalculator() *CalendarDateCalculator {
	return &CalendarDateCalculator{}
}

// CalculateCalendarDate рассчитывает календарную дату на основе радиоуглеродного возраста
// carbonAge в формате "3450 ± 30 BP"
// region для учета региональных вариаций
func (c *CalendarDateCalculator) CalculateCalendarDate(carbonAge string, region string) (string, error) {
	// Парсим радиоуглеродный возраст
	years, uncertainty, err := c.parseCarbonAge(carbonAge)
	if err != nil {
		return "", err
	}

	// Применяем калибровочную кривую
	calibratedYears := c.applyCalibrationCurve(years)

	// Учитываем неопределенность
	finalYears := c.applyUncertainty(calibratedYears, uncertainty)

	// Конвертируем в календарную дату (BC/AD)
	calendarDate := c.convertToCalendarDate(finalYears)

	// Форматируем результат
	return c.formatResult(calendarDate, uncertainty, region), nil
}

func (c *CalendarDateCalculator) parseCarbonAge(carbonAge string) (int, int, error) {
	// Регулярное выражение для парсинга "3450 ± 30 BP"
	re := regexp.MustCompile(`(\d+)\s*±\s*(\d+)\s*BP`)
	matches := re.FindStringSubmatch(carbonAge)

	if len(matches) != 3 {
		return 0, 0, fmt.Errorf("неверный формат радиоуглеродного возраста: %s", carbonAge)
	}

	years, err1 := strconv.Atoi(matches[1])
	uncertainty, err2 := strconv.Atoi(matches[2])

	if err1 != nil || err2 != nil {
		return 0, 0, fmt.Errorf("ошибка парсинга чисел в радиоуглеродном возрасте")
	}

	return years, uncertainty, nil
}

func (c *CalendarDateCalculator) applyCalibrationCurve(years int) float64 {
	// Упрощенная калибровочная кривая на основе INTCAL20
	// Для реального применения нужно использовать библиотеку калибровки
	baseYears := float64(years)

	// Нелинейная калибровка (упрощенная модель)
	if baseYears < 1000 {
		return baseYears * 0.95
	} else if baseYears < 5000 {
		return baseYears * 1.1
	} else {
		return baseYears * 1.25
	}
}

func (c *CalendarDateCalculator) applyUncertainty(calibratedYears float64, uncertainty int) float64 {
	// Учитываем неопределенность как стандартное отклонение
	uncertaintyFactor := float64(uncertainty) * 1.5
	return calibratedYears + uncertaintyFactor
}

func (c *CalendarDateCalculator) convertToCalendarDate(years float64) int {
	// Конвертируем радиоуглеродные годы в календарные годы BC/AD
	// 1950 год - это "present" в радиоуглеродном датировании
	const referenceYear = 1950
	calendarYear := referenceYear - int(math.Round(years))

	return calendarYear
}

func (c *CalendarDateCalculator) formatResult(calendarYear int, uncertainty int, region string) string {
	var era string
	var displayYear int

	if calendarYear >= 0 {
		era = "AD"
		displayYear = calendarYear
	} else {
		era = "BC"
		displayYear = -calendarYear
	}

	// Учитываем региональные особенности
	regionalAdjustment := c.getRegionalAdjustment(region)
	if regionalAdjustment != 0 {
		displayYear += regionalAdjustment
		if displayYear < 0 {
			era = "BC"
			displayYear = -displayYear
		}
	}

	return fmt.Sprintf("%d %s ± %d лет", displayYear, era, uncertainty)
}

func (c *CalendarDateCalculator) getRegionalAdjustment(region string) int {
	// Региональные поправки для разных географических зон
	adjustments := map[string]int{
		"Балтийское море":  -25,
		"Северная Европа":  -20,
		"Средиземноморье":  15,
		"Ближний Восток":   25,
		"Центральная Азия": 10,
		"Восточная Азия":   5,
		"Северная Америка": -15,
		"Южная Америка":    0,
	}

	if adjustment, exists := adjustments[region]; exists {
		return adjustment
	}
	return 0
}

// CalculateProbability рассчитывает вероятность соответствия материала
func (c *CalendarDateCalculator) CalculateProbability(sampleWeight float64, isotopes string, requirements string) int {
	// Проверяем соответствие веса требованиям
	weightScore := c.calculateWeightScore(sampleWeight)

	// Проверяем наличие необходимых изотопов
	isotopeScore := c.calculateIsotopeScore(isotopes)

	// Общая вероятность
	totalScore := (weightScore + isotopeScore) / 2

	return int(math.Round(totalScore * 100))
}

func (c *CalendarDateCalculator) calculateWeightScore(weight float64) float64 {
	// Идеальный вес образца: 1-5 грамм
	if weight >= 1.0 && weight <= 5.0 {
		return 1.0
	} else if weight >= 0.5 && weight < 1.0 {
		return 0.7
	} else if weight > 5.0 && weight <= 10.0 {
		return 0.8
	} else if weight > 0.1 && weight < 0.5 {
		return 0.4
	} else {
		return 0.1
	}
}

func (c *CalendarDateCalculator) calculateIsotopeScore(isotopes string) float64 {
	// Проверяем наличие ключевых изотопов для радиоуглеродного анализа
	requiredIsotopes := []string{"C14", "C13", "C12"}
	foundCount := 0

	for _, isotope := range requiredIsotopes {
		if strings.Contains(strings.ToUpper(isotopes), isotope) {
			foundCount++
		}
	}

	return float64(foundCount) / float64(len(requiredIsotopes))
}
