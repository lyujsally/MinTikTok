package service

type User struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   string `json:"follow_count,omitempty"`
	FollowerCount string `json:"follower_count,omitempty"`
	IsFollow      string `json:"is_follow,omitempty"`
}

type FriendUser struct {
	User
	Message string `protobuf:"bytes,1,opt,name=message" json:"message,omitempty"`
	MsgType int64  `protobuf:"varint,2,opt,name=msgType" json:"msgType"`
}

// 定义用户关系接口以及用户关系中的各种方法
type FollowService interface {

	// IsFollow 根据当前用户id和目标用户id来判断当前用户是否关注了目标用户
	IsFollow(userId int64, targetId int64) (bool, error)
	// GetFolloweeCount 根据用户id来查询用户关注了多少其它用户
	GetFolloweeCount(userId int64) (int64, error)
	// GetFollowerCount 根据用户id来查询用户被多少其他用户关注
	GetFollowerCount(userId int64) (int64, error)

	// Follow 当前用户关注目标用户
	Follow(userId int64, targetId int64) (bool, error)
	// UnFollow 当前用户取消对目标用户的关注
	UnFollow(userId int64, targetId int64) (bool, error)
	// GetFolloweeList 获取当前用户的关注列表
	GetFolloweeList(userId int64) ([]User, error)
	// GetFollowerList 获取当前用户的粉丝列表
	GetFollowerList(userId int64) ([]User, error)
	// GetFriendList 获取当前用户的好友列表
	GetFriendList(userId int64) ([]User, error)
}
