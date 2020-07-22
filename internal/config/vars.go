package config

var (
	AuthPort       = GetEnv("AUTH_PORT", "8002")
	NSQLookupdAddr = GetEnv("NSQLOOKUPD_ADDR", "127.0.0.1:4161")
	ChatPort       = GetEnv("CHAT_PORT", "8001")
	TokenKey       = GetEnv("SECRETE_KEY", "my_secrete")
	NSQAddr        = GetEnv("NSQ_ADDR", "127.0.0.1:4150")
	SendPort       = GetEnv("SEND_PORT", "8000")
	TokenName      = GetEnv("TOKEN_NAME", "token")
)
