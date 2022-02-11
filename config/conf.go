package config

var (
	//Storage mode options:MySQL/gob
	StoreMode = ""

	//Optional if you choose MySQL StoreMode
	MySQLURL      = "127.0.0.1:3306"
	MySQLUser     = ""
	MySQLPassword = ""
	MySQLDatabase = ""

	//New shortURL rate limit per minute
	IssueRateLimit = 30

	//MurmurHash Bit options：32/64
	MurmurBit = 32
)
