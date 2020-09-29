package infrastructure

type User struct {
	UserName string
	Email    string
}

type Repository interface {
	GetCurrentBranch() (string, error)
	GetUser() (User, error)
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
