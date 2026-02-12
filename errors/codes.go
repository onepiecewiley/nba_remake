package myErrors

// ErrorCode 错误码类型
type ErrorCode int

// 系统级错误码 (1000-1999)
const (
	CodeSuccess            ErrorCode = 0    // 成功
	CodeInternalError      ErrorCode = 1000 // 内部服务器错误
	CodeDatabaseError      ErrorCode = 1001 // 数据库错误
	CodeCacheError         ErrorCode = 1002 // 缓存错误
	CodeNetworkError       ErrorCode = 1003 // 网络错误
	CodeTimeout            ErrorCode = 1004 // 请求超时
	CodeServiceUnavailable ErrorCode = 1005 // 服务不可用
)

// 业务级错误码 (2000-2999)
const (
	CodeInvalidParam     ErrorCode = 2000 // 参数错误
	CodeMissingParam     ErrorCode = 2001 // 缺少参数
	CodeInvalidFormat    ErrorCode = 2002 // 格式错误
	CodeValidationFailed ErrorCode = 2003 // 验证失败
	CodeDuplicateData    ErrorCode = 2004 // 数据重复
	CodeDataNotFound     ErrorCode = 2005 // 数据不存在
	CodeDataExpired      ErrorCode = 2006 // 数据已过期
)

// 用户相关错误码 (3000-3999)
const (
	CodeUserNotFound     ErrorCode = 3000 // 用户不存在
	CodeUserExists       ErrorCode = 3001 // 用户已存在
	CodePasswordError    ErrorCode = 3002 // 密码错误
	CodePermissionDenied ErrorCode = 3003 // 权限不足
	CodeTokenExpired     ErrorCode = 3004 // Token过期
	CodeTokenInvalid     ErrorCode = 3005 // Token无效
)

// 球员相关错误码 (4000-4099)
const (
	CodePlayerNotFound    ErrorCode = 4000 // 球员不存在
	CodePlayerExists      ErrorCode = 4001 // 球员已存在
	CodeInvalidPlayerData ErrorCode = 4002 // 球员数据无效
	CodePlayerInUse       ErrorCode = 4003 // 球员正在使用中
)

// 球队相关错误码 (4100-4199)
const (
	CodeTeamNotFound    ErrorCode = 4100 // 球队不存在
	CodeTeamExists      ErrorCode = 4101 // 球队已存在
	CodeInvalidTeamData ErrorCode = 4102 // 球队数据无效
	CodeTeamFull        ErrorCode = 4103 // 球队人数已满
)

// 比赛相关错误码 (4200-4299)
const (
	CodeMatchNotFound    ErrorCode = 4200 // 比赛不存在
	CodeMatchExists      ErrorCode = 4201 // 比赛已存在
	CodeInvalidMatchData ErrorCode = 4202 // 比赛数据无效
	CodeMatchInProgress  ErrorCode = 4203 // 比赛进行中
	CodeMatchFinished    ErrorCode = 4204 // 比赛已结束
	CodeInvalidMatchTime ErrorCode = 4205 // 比赛时间无效
)

// 第三方服务错误码 (5000-5999)
const (
	CodeKafkaError         ErrorCode = 5000 // Kafka错误
	CodeKafkaSendFailed    ErrorCode = 5001 // Kafka发送失败
	CodeKafkaConsumeFailed ErrorCode = 5002 // Kafka消费失败
	CodeRedisError         ErrorCode = 5003 // Redis错误
	CodeMongoError         ErrorCode = 5004 // MongoDB错误
	CodeESError            ErrorCode = 5005 // Elasticsearch错误
)
