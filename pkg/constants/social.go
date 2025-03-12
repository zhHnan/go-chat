package constants

// 好友申请处理结果
type HandlerResult int

const (
	NoHandlerResult HandlerResult = iota + 1
	PassHandlerResult
	RefuseHandlerResult
	CancelHandlerResult
)

// 群等级 1.创建者 2.管理员 3.普通成员
type GroupRoleLevel int

const (
	CreatorGroupRoleLevel GroupRoleLevel = iota + 1 //
	ManagerGroupRoleLevel
	AtLargeGroupRoleLevel
)

// 进群方式 1.邀请， 2.申请
type GroupJoinSource int

const (
	InviteGroupJoinSource GroupJoinSource = iota + 1
	PutInGroupJoinSource
)
