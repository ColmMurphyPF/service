package cafe

type Cafe struct {
	ID          string `json:"id"`
	Name        string `json:"name" `
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
}

type NewCafe struct {
	Name        string `json:"name" validate:"required"`
	Address     string `json:"address" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}

type UpdateCafe struct {
	ID          string  `json:"id"`
	Name        *string `json:"name" validate:"omitempty"`
	Address     *string `json:"address" validate:"omitempty"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty"`
}
