package pgembed

const (
	connStr        = "host=10.0.0.213 port=5555 user=rupesh dbname=vectordb sslmode=disable"
	modelURL       = "http://localhost:11434/api/embeddings" // replace with your url
	EmbedderUrl    = "http://localhost:11434/api/embed"
	embedModel     = "nomic-embed-text"
	modelDimension = 768
	QueryLimit     = 2
)
