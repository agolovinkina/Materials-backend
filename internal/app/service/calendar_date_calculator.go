package service

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
)

type CalendarDateCalculator struct {
	regionalData map[string]struct {
		DeltaR      float64
		SigmaDeltaR float64
	}
}

type CalendarDateResult struct {
	FormattedDate string `json:"formatted_date"`
	Year          int    `json:"year"`
	ErrorRange    int    `json:"error_range"`
	Era           string `json:"era"`
}

func NewCalendarDateCalculator() *CalendarDateCalculator {
	calculator := &CalendarDateCalculator{}
	calculator.initializeRegionalData()
	return calculator
}

func (c *CalendarDateCalculator) initializeRegionalData() {
	c.regionalData = map[string]struct {
		DeltaR      float64
		SigmaDeltaR float64
	}{
		"Северная Атлантика (открытый океан)":  {400, 20},
		"Северная Атлантика (прибрежные зоны)": {200, 35},
		"Северное море":                 {150, 30},
		"Норвежское море":               {350, 30},
		"Баренцево море":                {400, 40},
		"Балтийское море (южная часть)": {300, 40},
		"Балтийское море (северная часть, Ботнический залив)": {500, 50},
		"Средиземное море (западная часть)":                   {50, 20},
		"Средиземное море (восточная часть)":                  {200, 30},
		"Эгейское море":                               {180, 30},
		"Чёрное море":                                 {525, 40},
		"Карибское море":                              {65, 25},
		"Мексиканский залив":                          {200, 30},
		"Тихий океан (тропики, запад)":                {200, 30},
		"Тихий океан (северо-запад, Охотское море)":   {700, 50},
		"Тихий океан (северо-восток, Берингово море)": {700, 40},
		"Тихий океан (южная часть)":                   {400, 40},
		"Индийский океан (север)":                     {150, 30},
		"Индийский океан (юг)":                        {400, 30},
		"Южный океан (Антарктика)":                    {1150, 60},
	}
}

func (c *CalendarDateCalculator) CalculateCalendarDate(carbonAgeStr string, region string) (*CalendarDateResult, error) {
	R, sigmaR, err := c.parseCarbonAge(carbonAgeStr)
	if err != nil {
		return nil, err
	}

	regionalParams, exists := c.regionalData[region]
	if !exists {
		return nil, fmt.Errorf("регион '%s' не найден в базе данных", region)
	}

	DeltaR := regionalParams.DeltaR
	sigmaDeltaR := regionalParams.SigmaDeltaR

	lambda, sigmaLambda := c.getCalibrationParams(R)

	t_cal, totalError := c.calculateCalendarDateFormula(R, DeltaR, lambda, sigmaR, sigmaDeltaR, sigmaLambda)

	calendarYear, era := c.convertToCalendarYear(t_cal)

	formattedDate := fmt.Sprintf("%d %s ± %d лет", calendarYear, era, int(math.Round(totalError)))
	if calendarYear < 0 {
		formattedDate = fmt.Sprintf("%d %s ± %d лет", -calendarYear, era, int(math.Round(totalError)))
	}

	return &CalendarDateResult{
		FormattedDate: formattedDate,
		Year:          calendarYear,
		ErrorRange:    int(math.Round(totalError)),
		Era:           era,
	}, nil
}

func (c *CalendarDateCalculator) calculateCalendarDateFormula(R, DeltaR, lambda, sigmaR, sigmaDeltaR, sigmaLambda float64) (float64, float64) {
	t_cal := 1950 - (R-DeltaR)/lambda

	errorTerm1 := math.Pow(sigmaR, 2)
	errorTerm2 := math.Pow(sigmaDeltaR, 2)
	errorTerm3 := math.Pow((R-DeltaR)/math.Pow(lambda, 2), 2) * math.Pow(sigmaLambda, 2)
	totalError := math.Sqrt(errorTerm1 + errorTerm2 + errorTerm3)

	return t_cal, totalError
}

func (c *CalendarDateCalculator) getCalibrationParams(R float64) (float64, float64) {
	switch {
	case R <= 150:
		return 1.00, 0.02
	case R <= 1000:
		return 1.005, 0.005
	case R <= 2000:
		return 1.025, 0.005
	case R <= 2500:
		return 1.035, 0.005
	case R <= 3500:
		return 1.045, 0.005
	case R <= 6000:
		return 1.055, 0.005
	default:
		return 1.06, 0.01
	}
}

func (c *CalendarDateCalculator) parseCarbonAge(carbonAge string) (float64, float64, error) {
	re := regexp.MustCompile(`(\d+)\s*±\s*(\d+)\s*BP`)
	matches := re.FindStringSubmatch(carbonAge)

	if len(matches) != 3 {
		return 0, 0, fmt.Errorf("неверный формат радиоуглеродного возраста: %s", carbonAge)
	}

	years, err1 := strconv.ParseFloat(matches[1], 64)
	uncertainty, err2 := strconv.ParseFloat(matches[2], 64)

	if err1 != nil || err2 != nil {
		return 0, 0, fmt.Errorf("ошибка парсинга чисел в радиоуглеродном возрасте")
	}

	return years, uncertainty, nil
}

func (c *CalendarDateCalculator) convertToCalendarYear(t_cal float64) (int, string) {
	year := int(math.Round(t_cal))

	if year >= 0 {
		return year, "AD"
	} else {
		return year, "BC"
	}
}

func (c *CalendarDateCalculator) CalculateProbability(carbonAgeStr string, region string) (*CalendarDateResult, error) {
	return c.CalculateCalendarDate(carbonAgeStr, region)
}
