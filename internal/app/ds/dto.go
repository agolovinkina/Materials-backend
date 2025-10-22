package ds

import (
	"time"
)

// MaterialAnalysisRequestDTO - DTO для заявки на анализ
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
	CarbonAge        string               `json:"CarbonAge"`
	TotalPrice       float64              `json:"TotalPrice"`
	RequestMaterials []RequestMaterialDTO `json:"RequestMaterials"`
}

// RequestMaterialDTO - DTO для связи заявка-материал
type RequestMaterialDTO struct {
	MaterialID         uint    `json:"MaterialID"`
	MaterialName       string  `json:"MaterialName"`
	MaterialImageURL   string  `json:"MaterialImageURL"`
	Comment            string  `json:"Comment"`
	ProbabilityPercent int     `json:"ProbabilityPercent"`
	CalendarDate       string  `json:"CalendarDate"`
	SampleDescription  string  `json:"SampleDescription"`
	SampleWeight       float64 `json:"SampleWeight"`
	IsPrimary          bool    `json:"IsPrimary"`
}

// MaterialDTO - DTO для материала
type MaterialDTO struct {
	MaterialID          uint   `json:"MaterialID"`
	MaterialName        string `json:"MaterialName"`
	MaterialImageURL    string `json:"MaterialImageURL"`
	MaterialDescription string `json:"MaterialDescription"`
	Isotopes            string `json:"Isotopes"`
	SampleRequirements  string `json:"SampleRequirements"`
}

// UpdateMaterialDTO - DTO для обновления материала
type UpdateMaterialDTO struct {
	MaterialName        string `json:"MaterialName"`
	MaterialDescription string `json:"MaterialDescription"`
	Isotopes            string `json:"Isotopes"`
	SampleRequirements  string `json:"SampleRequirements"`
}

// UpdateMaterialRequestDTO - DTO для обновления заявки
type UpdateMaterialRequestDTO struct {
	Region         string    `json:"Region"`
	ExpeditionDate time.Time `json:"ExpeditionDate"`
	CarbonAge      string    `json:"CarbonAge"`
}

// UpdateRequestMaterialDTO - DTO для обновления связи заявка-материал
type UpdateRequestMaterialDTO struct {
	Comment            string  `json:"Comment"`
	ProbabilityPercent int     `json:"ProbabilityPercent"`
	SampleDescription  string  `json:"SampleDescription"`
	SampleWeight       float64 `json:"SampleWeight"`
}

// UserDTO - DTO для пользователя
type UserDTO struct {
	UserID      uint   `json:"UserID"`
	Login       string `json:"Login"`
	IsModerator bool   `json:"IsModerator"`
}

// ChangeUserDTO - DTO для изменения пользователя
type ChangeUserDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// CartIconDTO - DTO для иконки корзины
type CartIconDTO struct {
	RequestID uint `json:"RequestID"`
	Count     int  `json:"Count"`
}
