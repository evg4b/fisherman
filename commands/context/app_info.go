package context

// AppInfo is application info structure
//
type AppInfo struct {
	Cwd                string
	AppPath            string
	IsRegisteredInPath bool
	GlobalConfigPath   string
	LocalConfigPath    string
}
