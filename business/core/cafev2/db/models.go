package db

type Cafe struct {
	ID      string `db:"cafe_id"`
	OwnerID string `db:"owner_id"`
	Name    string `db:"cafe_name"`
	Address string `db:"address"`
	LogoURL string `db:"logo_url"`
}
