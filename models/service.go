package models

// APIService represents the api_services table
type APIService struct {
	ID          int    `json:"id" db:"id"`                         // Ensure `db` tag matches the column name in the database
	ProjectID   int    `json:"api_project_id" db:"api_project_id"` // Use `db` tag to map the column name
	Name        string `json:"name" db:"name"`
	Endpoint    string `json:"endpoint" db:"endpoint"`
	Description string `json:"description,omitempty" db:"description"`
	Status      string `json:"status" db:"status"`
}
