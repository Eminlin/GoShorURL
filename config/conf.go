package config

var (
	Domain        = "hatchurl.com"
	SSL           = true
	SubPath       = ""
	MySQLURL      = "127.0.0.1:3306"
	MySQLUser     = ""
	MySQLPassword = ""
	MySQLDatabase = ""
	RedisURL      = "127.0.0.1:6639"

	//ShortURL generation rate limit per minute
	IssueRateLimit = 30

	//MurmurHash Bit options：32/64
	MurmurBit = 32

	//
	DuplicateURL = false
)
