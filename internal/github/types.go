package github

// Repository contains GitHub repository information
type Repository struct {
	Name            string
	Description     string
	Owner           string
	Language        string
	StargazersCount int
	ForksCount      int
}
