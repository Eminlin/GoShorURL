package model

//App app.conf struct
type App struct {
	App struct {
		APIPort, Host             string
		SSL, DuplicateURL         bool
		IssueRateLimit, MurmurBit int16
		LogLevel, NotFoundPage    string
	}
	MySQL struct {
		HostPort, User, Password, Database string
		URLTable, ManageTable              string
	}
	Redis struct {
		HostPort, Password, Key string
		DB                      int
	}
}
