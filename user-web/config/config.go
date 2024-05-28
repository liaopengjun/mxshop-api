package config

type UserSrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Expire   int    `mapstructure:"expire" json:"expire"`
	DB       int    `mapstructure:"db" json:"db"`
	Password string `mapstructure:"password" json:"password"`
	PoolSize int    `mapstructure:"pool_size" json:"pool_size"`
}

type LogConfig struct {
	Level         string `mapstructure:"level"`          // 级别
	Format        string `mapstructure:"format"`         // 输出
	Prefix        string `mapstructure:"prefix"`         // 日志前缀
	Director      string `mapstructure:"director"`       // 日志文件夹
	ShowLine      bool   `mapstructure:"show-line"`      // 显示行
	EncodeLevel   string `mapstructure:"encode-level"`   // 编码级
	StacktraceKey string `mapstructure:"stacktrace-key"` // 栈名
	LogInConsole  bool   `mapstructure:"log-in-console"` // 输出控制台
}

type ServerConfig struct {
	Name          string        `mapstructure:"name" json:"name"`
	Host          string        `mapstructure:"host" json:"host"`
	Tags          []string      `mapstructure:"tags" json:"tags"`
	Port          int           `mapstructure:"port" json:"port"`
	UserSrvInfo   UserSrvConfig `mapstructure:"user_srv" json:"user_srv"`
	JWTInfo       JWTConfig     `mapstructure:"jwt" json:"jwt"`
	RedisInfo     RedisConfig   `mapstructure:"redis" json:"redis"`
	ConsulInfo    ConsulConfig  `mapstructure:"consul" json:"consul"`
	LogConfigInfo LogConfig     `mapstructure:"log" json:"log"`
}
