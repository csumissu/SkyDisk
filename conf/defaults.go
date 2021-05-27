package conf

var DatabaseCfg = &database{}

var ServerCfg = &server{
	Port: 8080,
}

var CORSCfg = &cors{
	AllowOrigins:     []string{},
	AllowMethods:     []string{"PUT", "POST", "GET", "OPTIONS", "DELETE"},
	AllowHeaders:     []string{"Cookie", "Authorization", "Content-Length", "Content-Type"},
	AllowCredentials: false,
	ExposeHeaders:    nil,
}
