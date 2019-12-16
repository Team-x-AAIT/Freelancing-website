package entities

// Project is a stuct that defines the type Project.
type Project struct {
	ID           string
	Title        string
	Description  string
	Details      string
	AttachedFile string
	Category     string
	Subcategory  string
	Budget       float64
	WorkType     int64
}

// NewProject is a function that returns a new Project type from provided arguments.
func NewProject(title, description, details, attachedFile, catagory, subcatagroy string, budget float64, worktype int64) *Project {
	project := Project{
		Title:        title,
		Description:  description,
		Details:      details,
		AttachedFile: attachedFile,
		Category:     catagory,
		Subcategory:  subcatagroy,
		Budget:       budget,
		WorkType:     worktype}

	return &project

}
