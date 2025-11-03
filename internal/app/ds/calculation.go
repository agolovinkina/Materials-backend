package ds

import (
	"math"
)

// CarbonDatingCalculation - модель для расчета радиоуглеродного возраста
type CarbonDatingCalculation struct {
	SampleActivity   float64 `json:"sample_activity"`
	ModernStandard   float64 `json:"modern_standard"`
	MeasurementError float64 `json:"measurement_error"`
	HalfLife         float64 `json:"half_life"`
}

// CarbonDatingResult - результат расчета
type CarbonDatingResult struct {
	AgeYears        float64 `json:"age_years"`        // Основной возраст в годах (float)
	AgeValue        int     `json:"age_value"`        // Целое значение возраста (3450)
	ErrorRange      float64 `json:"error_range"`      // Значение погрешности в годах (float)
	ErrorValue      int     `json:"error_value"`      // Целое значение погрешности (30)
	ConfidenceLevel float64 `json:"confidence_level"` // Уровень достоверности
}

// CalculateCarbonAge - расчет радиоуглеродного возраста
func (c *CarbonDatingCalculation) CalculateCarbonAge() CarbonDatingResult {
	halfLife := c.HalfLife
	if halfLife == 0 {
		halfLife = 5568.0
	}

	decayConstant := math.Ln2 / halfLife

	if c.ModernStandard == 0 {
		return CarbonDatingResult{
			AgeYears:        0,
			AgeValue:        0,
			ErrorRange:      0,
			ErrorValue:      0,
			ConfidenceLevel: 0,
		}
	}

	activityRatio := c.SampleActivity / c.ModernStandard

	var age float64
	if activityRatio > 0 {
		age = (1 / decayConstant) * math.Log(1/activityRatio)
	}

	errorRange := c.calculateError(age, activityRatio)
	confidenceLevel := c.calculateConfidenceLevel(activityRatio, c.MeasurementError)

	ageValue := int(math.Round(age))
	errorValue := int(math.Round(errorRange))

	// Корректируем ошибку в пределах разумного
	if errorValue < 10 {
		errorValue = 10
	}
	if errorValue > 100 {
		errorValue = 100
	}

	return CarbonDatingResult{
		AgeYears:        age,
		AgeValue:        ageValue,
		ErrorRange:      errorRange,
		ErrorValue:      errorValue,
		ConfidenceLevel: confidenceLevel,
	}
}

func (c *CarbonDatingCalculation) calculateError(age, activityRatio float64) float64 {
	if c.MeasurementError == 0 {
		return age * 0.01
	}
	errorPropagation := c.MeasurementError / activityRatio
	maxError := age * 0.1
	if errorPropagation > maxError {
		return maxError
	}
	return errorPropagation
}

func (c *CarbonDatingCalculation) calculateConfidenceLevel(activityRatio, measurementError float64) float64 {
	signalToNoise := activityRatio / measurementError
	confidence := 100.0 * (1 - math.Exp(-signalToNoise/5))
	if confidence > 99.9 {
		return 99.9
	}
	if confidence < 0.1 {
		return 0.1
	}
	return confidence
}
