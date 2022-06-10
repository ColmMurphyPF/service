package cafev2

import (
	"unsafe"

	"github.com/colmmurphy91/go-service/business/core/cafev2/db"
)

type Cafe struct {
	ID      string `json:"cafe_id"`
	OwnerID string `json:"-"`
	Name    string `json:"cafe_name"`
	Address string `json:"address"`
	LogoURL string `json:"logo_url"`
}

// NewCafe is what we require from clients when adding a Product.
type NewCafe struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
	LogoURL string `json:"logo_url"`
}

// UpdateProduct defines what information may be provided to modify an
// existing Product. All fields are optional so clients can send just the
// fields they want changed. It uses pointer fields so we can differentiate
// between a field that was not provided and a field that was provided as
// explicitly blank. Normally we do not want to use pointers to basic types but
// we make exceptions around marshalling/unmarshalling.
//type UpdateProduct struct {
//	Name     *string `json:"name"`
//	Cost     *int    `json:"cost" validate:"omitempty,gte=0"`
//	Quantity *int    `json:"quantity" validate:"omitempty,gte=1"`
//}

// =============================================================================

func toCafe(dbCaf db.Cafe) Cafe {
	cu := (*Cafe)(unsafe.Pointer(&dbCaf))
	return *cu
}
