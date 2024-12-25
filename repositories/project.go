package repositories

import (
	"blue-owl-service/models"

	"github.com/jmoiron/sqlx"
)

// GetListProjects retrieves all projects from the database
func GetListProjects(db *sqlx.DB) ([]models.Project, error) {
	var projects []models.Project
	query := `SELECT id, name, description, status FROM api_projects`

	err := db.Select(&projects, query)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

// GetProjectByID fetches a single project by ID from the database
func GetProjectByID(db *sqlx.DB, projectID int) (*models.Project, error) {
	var project models.Project
	query := `SELECT id, name, description, status FROM api_projects WHERE id = $1`

	err := db.Get(&project, query, projectID)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// CreateProject inserts a new project into the database
func CreateProject(db *sqlx.DB, project *models.Project) (int, error) {
	query := `INSERT INTO api_projects (name, description, status) 
              VALUES ($1, $2, $3) RETURNING id`

	var projectID int
	err := db.QueryRow(query, project.Name, project.Description, project.Status).Scan(&projectID)
	if err != nil {
		return 0, err
	}
	return projectID, nil
}

// UpdateProject updates an existing project in the database
func UpdateProject(db *sqlx.DB, project *models.Project) error {
	query := `UPDATE api_projects SET name = $1, description = $2, status = $3 WHERE id = $4`

	_, err := db.Exec(query, project.Name, project.Description, project.Status, project.ID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteProject deletes a project from the database by ID
func DeleteProject(db *sqlx.DB, projectID int) error {
	query := `DELETE FROM api_projects WHERE id = $1`

	_, err := db.Exec(query, projectID)
	if err != nil {
		return err
	}
	return nil
}

// CheckProjectExistsByName checks if a project with the given name exists (case-insensitive)
func CheckProjectExistsByName(db *sqlx.DB, name string) (bool, error) {
	query := `SELECT COUNT(*) FROM api_projects WHERE LOWER(name) = LOWER($1)`

	var count int
	err := db.Get(&count, query, name)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
