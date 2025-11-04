package ds

import (
	"time"
)

type MaterialAnalysisRequestDTO struct {
	RequestID        int                  `json:"RequestID"`
	RequestStatus    string               `json:"RequestStatus"`
	CreatedAt        time.Time            `json:"CreatedAt"`
	CreatorLogin     string               `json:"CreatorLogin"`
	FormedAt         *time.Time           `json:"FormedAt,omitempty"`
	CompletedAt      *time.Time           `json:"CompletedAt,omitempty"`
	ModeratorLogin   *string              `json:"ModeratorLogin,omitempty"`
	Region           string               `json:"Region"`
	ExpeditionDate   time.Time            `json:"ExpeditionDate"`
	CarbonAgeValue   int                  `json:"CarbonAgeValue"`
	CarbonAgeError   int                  `json:"CarbonAgeError"`
	TotalPrice       float64              `json:"TotalPrice"`
	RequestMaterials []RequestMaterialDTO `json:"RequestMaterials"`
}

type RequestMaterialDTO struct {
	MaterialID        uint    `json:"MaterialID"`
	MaterialName      string  `json:"MaterialName"`
	MaterialImageURL  string  `json:"MaterialImageURL"`
	Comment           string  `json:"Comment"`
	CalendarYear      int     `json:"CalendarYear"`
	CalendarError     int     `json:"CalendarError"`
	FormattedDate     string  `json:"FormattedDate"`
	SampleDescription string  `json:"SampleDescription"`
	SampleWeight      float64 `json:"SampleWeight"`
	IsPrimary         bool    `json:"IsPrimary"`
}

type MaterialDTO struct {
	MaterialID          uint   `json:"MaterialID"`
	MaterialName        string `json:"MaterialName"`
	MaterialImageURL    string `json:"MaterialImageURL"`
	MaterialDescription string `json:"MaterialDescription"`
	Isotopes            string `json:"Isotopes"`
	SampleRequirements  string `json:"SampleRequirements"`
}

type UpdateMaterialDTO struct {
	MaterialName        string `json:"MaterialName"`
	MaterialDescription string `json:"MaterialDescription"`
	Isotopes            string `json:"Isotopes"`
	SampleRequirements  string `json:"SampleRequirements"`
}

type UpdateMaterialRequestDTO struct {
	Region         string    `json:"Region"`
	ExpeditionDate time.Time `json:"ExpeditionDate"`
	CarbonAgeValue int       `json:"CarbonAgeValue"`
	CarbonAgeError int       `json:"CarbonAgeError"`
}

type UpdateRequestMaterialDTO struct {
	Comment           string  `json:"Comment"`
	SampleDescription string  `json:"SampleDescription"`
	SampleWeight      float64 `json:"SampleWeight"`
}

type UserDTO struct {
	UserID      uint   `json:"UserID"`
	Login       string `json:"Login"`
	IsModerator bool   `json:"IsModerator"`
}

type ChangeUserDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type CartIconDTO struct {
	RequestID uint `json:"RequestID"`
	Count     int  `json:"Count"`
}
