package conf

var DatabaseCfg = &database{}

var ServerCfg = &server{
	Port: 8080,
}

var RedisCfg = &redis{
	Network: "tcp",
	Port:    6379,
	DB:      "0",
}

var CORSCfg = &cors{
	AllowOrigins:     []string{},
	AllowMethods:     []string{"PUT", "POST", "GET", "OPTIONS", "DELETE"},
	AllowHeaders:     []string{"Cookie", "Authorization", "Content-Length", "Content-Type"},
	AllowCredentials: false,
	ExposeHeaders:    nil,
}
