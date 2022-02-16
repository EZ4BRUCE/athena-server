package setting

// RPC服务器配置结构体
type RPCSettingS struct {
	Port            string // gRPC端口
	AggregationTime int32  // 聚合时间（单位时间的个数）
	MaxConn         int    // 初始化用户注册表的大小
}

// 日志写入配置结构体
type LOGSettingS struct {
	LogSavePath string // 日志文件的位置
	LogFileName string // 日志文件名
	LogFileExt  string // 日志文件后缀
	MaxSize     int    // 在进行切割之前，日志文件的最大大小（以MB为单位）
	MaxBackups  int    // 保留旧文件的最大个数
	MaxAge      int    // 保留旧文件的最大天数
	Compress    bool   // 是否压缩/归档旧文件
}

// 邮件配置结构体
type EmailSettingS struct {
	Host     string // 发送服务器ip
	Port     int    // 端口
	UserName string // 发送用户名
	Password string // smtp口令
	IsSSL    bool
	From     string
	To       []string
}

// 告警数据库配置结构体
type ReportDBSettingS struct {
	DBType string
	Host   string
	DBName string
}

// 规则数据库配置结构体
type RuleDBSettingS struct {
	DBType       string
	UserName     string
	Password     string
	Host         string
	DBName       string
	Charset      string
	ParseTime    bool
	MaxIdleConns int
	MaxOpenConns int
}

// 配置读取函数
func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}
	return nil
}
