package config

import "time"

//time.Duration int64

type Server struct {
	RunMode               string        `yaml:"RunMode"`               // gin 的运行模式（release 是生产模式）
	HttpPort              string        `yaml:"HttpPort"`              // 默认的 HTTP 监听端口
	ReadTimeout           time.Duration `yaml:"ReadTimeout"`           // 允许读取的最大持续时间
	WriteTimeout          time.Duration `yaml:"WriteTimeout"`          // 允许写入的最大持续时间
	DefaultContextTimeout time.Duration `yaml:"DefaultContextTimeout"` // 默认上下文超时
}

type AppConfig struct {
	Name      string `yaml:"Name"`      //App名
	Version   string `yaml:"Version"`   //版本号
	MachineID int64  `yaml:"MachineID"` //机器ID
	StartTime string `yaml:"StartTime"` //启动时间
}

type PublicConfig struct {
	Server Server      `yaml:"Server"`
	Log    LogConfig   `yaml:"Log"`
	App    AppConfig   `yaml:"App"`
	Page   PageConfig  `yaml:"Page"`
	Rules  RulesConfig `yaml:"Rules"`
	Auto   Auto        `yaml:"Auto"`
	Worker Worker      `yaml:"Worker"`
	Limit  Limit       `yaml:"Limit"`
}

type PrivateConfig struct {
	Mysql MysqlConfig `yaml:"Mysql"`
	Redis RedisConfig `yaml:"Redis"`
	Email Email       `yaml:"Email"`
	Token Token       `yaml:"Token"`
	//HuaWeiOBS  HuaWeiOBS        `yaml:"HuaWeiOBS"`
	//RocketMQ   RocketMQ         `yaml:"RocketMQ"`
}

type LogConfig struct {
	Level         string `yaml:"Level"`         // 日志级别
	LogSavePath   string `yaml:"LogSavePath"`   // 日志保存路径
	HighLevelFile string `yaml:"HighLevelFile"` // 高级别日志名
	LowLevelFile  string `yaml:"LowLevelFile"`  // 低级别日志名
	LogFileExt    string `yaml:"LogFileExt"`    // 日志文件后缀
	MaxSize       int    `yaml:"MaxSize"`       // 最大大小（MB）
	MaxAge        int    `yaml:"MaxAge"`        // 最大保存天数
	MaxBackups    int    `yaml:"MaxBackups"`    // 最大备份数
	Compress      bool   `yaml:"Compress"`      // 是否压缩
}

type PageConfig struct {
	DefaultPageSize int32  `yaml:"DefaultPageSize"`
	MaxPageSize     int32  `yaml:"MaxPageSize"`
	PageKey         string `yaml:"PageKey"`
	PageSizeKey     string `yaml:"PageSizeKey"`
}

// RulesConfig 注册相关规则
type RulesConfig struct {
	UsernameLenMax   int           `yaml:"UsernameLenMax"`   //用户名最大长度
	UsernameLenMin   int           `yaml:"UsernameLenMin"`   //用户名最小长度
	PasswordLenMax   int           `yaml:"PasswordLenMax"`   //密码最大长度
	PasswordLenMin   int           `yaml:"PasswordLenMin"`   //密码最小长度
	CodeLength       int           `yaml:"CodeLength"`       //验证码长度
	AccountNumMax    int64         `yaml:"AccountNumMax"`    //用户账号最大数量
	BiggestFileSize  int64         `yaml:"BiggestFileSize"`  //最大文件大小
	UserMarkDuration time.Duration `yaml:"UserMarkDuration"` //用户发送验证码间隔时间
	CodeMarkDuration time.Duration `yaml:"CodeMarkDuration"` //验证码有效时间
	DefaultAvatarURL string        `yaml:"DefaultAvatarURL"` //默认头像
}

type MysqlConfig struct {
	DirverName     string `yaml:"DirverName"`
	DataSourceName string `yaml:"DataSourceName"`
	MaxOpenConns   int    `yaml:"MaxOpenConns"`
	MaxIdleConns   int    `yaml:"MaxIdleConns"`
}

type RedisConfig struct {
	Address   string        `yaml:"Addrrss"`   //Redis 服务器地址
	Password  string        `yaml:"Password"`  //认证密码
	DB        int           `yaml:"DB"`        //Redis 数据库索引
	PoolSize  int           `yaml:"PoolSize"`  //Redis 连接池大小
	CacheTime time.Duration `yaml:"CacheTime"` //缓存时间
}

type Email struct {
	Username string   `yaml:"Username"` //登陆邮箱的用户名
	Password string   `yaml:"Password"`
	Host     string   `yaml:"Host"`  //邮箱服务器的主机地址
	From     string   `yaml:"From"`  //发件人邮箱
	To       []string `yaml:"To"`    //收件人邮箱
	Port     int      `yaml:"Port"`  //邮箱服务器端口号
	IsSSL    bool     `yaml:"IsSSL"` // 是否使用 SSL 加密
}

type Token struct {
	Key                  string        `yaml:"Key"`                  // 生成 token 的密钥
	AccessTokenExpire    time.Duration `yaml:"AccessTokenExpire"`    // 用户 token 的访问令牌
	RefreshTokenExpire   time.Duration `yaml:"RefreshTokenExpire"`   // 用户 token 的刷新令牌
	AccountTokenDuration time.Duration `yaml:"AccountTokenDuration"` // 账户 token 的有效期限
	AuthorizationKey     string        `yaml:"AuthorizationKey"`     // 授权密钥，用于进行授权验证
	AuthorizationType    string        `yaml:"AuthorizationType"`    // 授权类型，指定授权的具体方式或策略
}

// Limit 限流
type Limit struct {
	IPLimit  IPLimit  `json:"IPLimit"`
	APILimit APILimit `json:"APILimit"`
}

// IPLimit IP 限流
type IPLimit struct {
	Cap     int64 `yaml:"Cap"`     // 令牌桶容量
	GenNum  int64 `yaml:"GenNum"`  // 每次生成的令牌数量
	GenTime int64 `yaml:"GenTime"` // 生成令牌的时间间隔，即每个多长时间生成一次令牌
	Cost    int64 `yaml:"Cost"`    // 每次请求消耗的令牌数量
}

// APILimit API 限流
type APILimit struct {
	Count    int           `yaml:"Count"`    // 令牌桶容量
	Duration time.Duration `yaml:"Duration"` // 填充令牌桶的时间间隔，即每隔多长时间会填充一次令牌
	Burst    int           `yaml:"Burst"`    // 令牌桶的最大容忍峰值，即在某个时间点可以容忍的最大请求数量
}

// Retry 重试
type Retry struct {
	Duration time.Duration `yaml:"Duration"` // 重试的时间间隔
	MaxTimes int           `yaml:"MaxTimes"` // 最大重试次数
}

// Auto 自动任务配置
type Auto struct {
	Retry                     Retry         `yaml:"Retry"`
	DeleteExpiredFileDuration time.Duration `yaml:"DeleteExpiredFileDuration"` // 删除过期文件的时间
}

// Worker 工作池配置
type Worker struct {
	TaskChanCapacity   int `yaml:"TaskChanCapacity"`   // 任务队列容量
	WorkerChanCapacity int `yaml:"WorkerChanCapacity"` // 工作队列容量
	WorkerNum          int `yaml:"WorkerNum"`          // 工作池数
}
