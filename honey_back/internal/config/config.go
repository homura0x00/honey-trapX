package config

// Config 应用配置：settings.yaml 解析结果。无包级全局变量，全部由 main 装配注入。
type Config struct {
	Server  Server  `yaml:"server"`
	Mysql   MySQL   `yaml:"mysql"`
	Redis   Redis   `yaml:"redis"`
	Session Session `yaml:"session"`
	Admin   Admin   `yaml:"admin"`
	Log     Log     `yaml:"log"`
}

type Server struct {
	Port        int    `yaml:"port"`
	Mode        string `yaml:"mode"`
	Timeout     int    `yaml:"timeout"`     // 单个外部操作超时（秒）
	FrontendURL string `yaml:"frontendURL"` // CORS 允许来源
}

type Log struct {
	Level  string `yaml:"level"`  // debug/info/warn/error
	Output string `yaml:"output"` // stdout/file/both
}

type MySQL struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Config   string `yaml:"config"` // 追加到 DSN 的高级参数
	DBName   string `yaml:"dbname"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	LogLevel string `yaml:"logLevel"`
}

type Redis struct {
	IP       string `yaml:"ip"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	PoolSize int    `yaml:"poolSize"`
}

// Session redis session + Cookie 会话配置
type Session struct {
	CookieName   string `yaml:"cookieName"`
	RedisPrefix  string `yaml:"redisPrefix"`
	ExpireHours  int    `yaml:"expireHours"`
	CookieSecure bool   `yaml:"cookieSecure"` // 生产 HTTPS 下应开
}

// Admin 无用户时引导的默认管理员
type Admin struct {
	Account  string `yaml:"account"`
	Password string `yaml:"password"`
	UserName string `yaml:"userName"`
}
