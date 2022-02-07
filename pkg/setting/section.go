package setting

// RPC服务器配置结构体
type RPCSettingS struct {
	Port            string
	AggregationTime int
	MaxConn         int
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

// 邮件配置结构体
type EmailSettingS struct {
	Host     string
	Port     int
	UserName string
	Password string
	IsSSL    bool
	From     string
	To       []string
}

// 配置读取函数
func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}
	return nil
}
