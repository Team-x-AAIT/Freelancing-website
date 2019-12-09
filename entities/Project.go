package entities

// Project is a stuct that defines the type Project.
type Project struct {
	ID          string
	Title       string
	Description string
	Budget      float64
	Category    string
	Subcategory string
	WorkType    string
}

// NewProject is a function that returns a new Project type from provided arguments.
func NewProject(id, title, description, catagory, subcatagroy, worktype string, budget float64) *Project {
	project := Project{
		ID:          id,
		Title:       title,
		Description: description,
		Budget:      budget,
		Category:    catagory,
		Subcategory: subcatagroy,
		WorkType:    worktype}

	return &project

}
