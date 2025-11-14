package ds

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
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

type AnalysisRequestDTO struct {
	AnalysisRequestID     int                `json:"AnalysisRequestID"`
	AnalysisRequestStatus string             `json:"AnalysisRequestStatus"`
	CreatedAt             time.Time          `json:"CreatedAt"`
	CreatorLogin          string             `json:"CreatorLogin"`
	FormedAt              *time.Time         `json:"FormedAt,omitempty"`
	CompletedAt           *time.Time         `json:"CompletedAt,omitempty"`
	ModeratorLogin        *string            `json:"ModeratorLogin,omitempty"`
	TextToAnalyse         string             `json:"TextToAnalyse"`
	Genres                []AnalysisGenreDTO `json:"Genres"`
}

type AnalysisGenreDTO struct {
	GenreID            int    `json:"GenreID"`
	GenreName          string `json:"GenreName"`
	GenreImageURL      string `json:"GenreImageURL"`
	CommentToRequest   string `json:"CommentToRequest"`
	ProbabilityPercent int    `json:"ProbabilityPercent"`
}

type UpdateAnalysisRequestDTO struct {
	TextToAnalyse string `json:"TextToAnalyse"`
}

type UpdateGenreRequestDTO struct {
	CommentToRequest   string `json:"CommentToRequest"`
	ProbabilityPercent int    `json:"ProbabilityPercent"`
}

type GenreDTO struct {
	GenreID       int    `json:"GenreID"`
	GenreName     string `json:"GenreName"`
	GenreImageURL string `json:"GenreImageURL"`
	GenreKeywords string `json:"GenreKeywords"`
}

type UpdateGenreDTO struct {
	GenreName     string `json:"GenreName"`
	GenreKeywords string `json:"GenreKeywords"`
}

type UserDTO struct {
	UserID int      `json:"UserID"`
	Login  string   `json:"Login"`
	Role   UserRole `json:"Role"`
}

type ChangeUserDTO struct {
	Login    string `json:"Login"`
	Password string `json:"Password"`
}

type CartIconDTO struct {
	RequestID uint `json:"RequestID"`
	Count     int  `json:"Count"`
}

type UserRole string

const (
	RoleGuest     UserRole = "guest"
	RoleCreator   UserRole = "creator"
	RoleModerator UserRole = "moderator"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID int      `json:"UserID"`
	Role   UserRole `json:"Role"`
}

type AuthResponseDTO struct {
	AccessToken string `json:"AccessToken"`
	TokenType   string `json:"TokenType"`
	ExpiresIn   int64  `json:"ExpiresIn"`
}
