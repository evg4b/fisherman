package context

// AppInfo is application info structure
//
type AppInfo struct {
	AppPath            string
	IsRegisteredInPath bool
	GlobalConfigPath   string
	LocalConfigPath    string
	Cwd                string
}
