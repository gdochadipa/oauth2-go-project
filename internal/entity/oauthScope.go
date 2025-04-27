package entity

// OAuthScope represents an authorization scope.
type OAuthScope struct {
	Id   string `gorm:"primaryKey" json:"id"`
	Name string `gorm:"not null" json:"name"`
}
