package infrastructure

type Repository interface {
	GetCurrentBranch() string
}

type FileWriter interface {
	Write(path, content string) error
}

type FileReader interface {
	Read(path string) (string, error)
}

type FileAccessor interface {
	FileWriter
	FileReader
	Exist(path string) bool
}
