package ds

import (
	"fmt"
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
	CarbonAge       string  `json:"carbon_age"`
	AgeYears        float64 `json:"age_years"`
	ErrorRange      float64 `json:"error_range"`
	ConfidenceLevel float64 `json:"confidence_level"`
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
			CarbonAge:       "не определен",
			AgeYears:        0,
			ErrorRange:      0,
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

	return CarbonDatingResult{
		AgeYears:        age,
		ErrorRange:      errorRange,
		ConfidenceLevel: confidenceLevel,
		CarbonAge:       formatCarbonAge(age, errorRange),
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

func formatCarbonAge(age, errorRange float64) string {
	ageInt := int(math.Round(age))
	errorInt := int(math.Round(errorRange))
	if errorInt < 10 {
		errorInt = 10
	}
	if errorInt > 100 {
		errorInt = 100
	}
	return fmt.Sprintf("%d ± %d BP", ageInt, errorInt)
}
