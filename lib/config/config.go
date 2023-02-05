package config

var Config *Configuration

type Configuration struct {
	GlobalVars    map[string]interface{}
	AdminPassword string
}

func Init() {
	Config = &Configuration{
		GlobalVars: map[string]interface{}{
			"ssh_user":     "test",
			"ssh_password": "test",
		},
		AdminPassword: "adminpassword",
	}
}
