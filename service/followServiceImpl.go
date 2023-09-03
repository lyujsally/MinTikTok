package service

import (
	"strconv"
	"strings"
	"sync"

	"github.com/lyujsally/MinTikTok-lyujsally/conf"
	"github.com/lyujsally/MinTikTok-lyujsally/dao/mysql"
	redisDao "github.com/lyujsally/MinTikTok-lyujsally/dao/redis"
	"github.com/lyujsally/MinTikTok-lyujsally/kitex_gen/relation"
	"github.com/lyujsally/MinTikTok-lyujsally/pkg/kafka"
)

type FollowServiceImp struct{}

var (
	followServiceImp  *FollowServiceImp //controller层通过该实例变量调用service的所有业务方法
	followServiceOnce sync.Once         //限定该service对象为单例，节约内存
)

// 生成并返回FollowServiceImp结构体单例变量
func NewFollowInstance() *FollowServiceImp {
	followServiceOnce.Do(
		func() {
			followServiceImp = &FollowServiceImp{}
		})
	return followServiceImp
}

// IsFollow 根据当前用户id和目标用户id来判断当前用户是否关注了目标用户
func (*FollowServiceImp) IsFollow(userId int64, entityId int64) (bool, error) {
	//先查Redis里面是否有此关系
	result, err := redisDao.RdbIsFollow.ZRank(strconv.FormatInt(userId, 10), strconv.FormatInt(entityId, 10)).Result()
	if err == nil && result != -1 {
		//重设过期时间
		redisDao.RdbIsFollow.Expire(strconv.FormatInt(userId, 10), conf.ExpireTime)
		return true, err
	}

	//SQL查询
	fr, err := mysql.NewRelationDaoInstance().FindRelation(userId, entityId)

	if err != nil {
		return false, err
	}
	if fr == nil {
		return false, nil
	}

	go redisDao.AddRelationToRedis(userId, entityId)

	return true, nil
}

// GetFolloweeCount 根据用户id来查询用户关注了多少其它用户
func (*FollowServiceImp) GetFolloweeCount(userId int64) (int64, error) {
	// 查Redis中是否已经存在。
	if cnt, err := redisDao.RdbFollowee.ZCard(strconv.FormatInt(userId, 10)).Result(); cnt > 0 {
		// 更新过期时间。
		redisDao.RdbFollowee.Expire(strconv.FormatInt(userId, 10), conf.ExpireTime)
		return cnt - 1, err
	}
	// SQL中查询。
	ids, err := mysql.NewRelationDaoInstance().GetFolloweeIds(userId)
	if nil != err {
		return 0, err
	}
	// 更新redis里的followers 和 followingPart
	go redisDao.AddFolloweesToRedis(userId, ids)

	return int64(len(ids)), err
}

// GetFollowerCount 根据用户id来查询用户被多少其他用户关注
func (*FollowServiceImp) GetFollowerCount(userId int64) (int64, error) {
	// 查Redis中是否已经存在。
	if cnt, err := redisDao.RdbFollower.ZCard(strconv.FormatInt(userId, 10)).Result(); cnt > 0 {
		// 更新过期时间。
		redisDao.RdbFollower.Expire(strconv.FormatInt(userId, 10), conf.ExpireTime)
		return cnt - 1, err
	}
	// SQL中查询。
	ids, err := mysql.NewRelationDaoInstance().GetFollowersIds(userId)
	if nil != err {
		return 0, err
	}
	// 更新redis里的followers 和 followingPart
	go redisDao.AddFollowersToRedis(userId, ids)

	return int64(len(ids)), err
}

// Follow 当前用户关注目标用户
func (*FollowServiceImp) Follow(userId int64, targetId int64) (bool, error) {
	//将关注消息加入消息队列
	msg := strings.Builder{}
	msg.WriteString(strconv.Itoa(int(userId)))
	msg.WriteString(",")
	msg.WriteString(strconv.Itoa(int(targetId)))
	kafka.KfkFollowAdd.FollowProducer(msg.String())
	return redisDao.UpdateRedisFollow(userId, targetId)
}

// UnFollow 当前用户取消对目标用户的关注
func (*FollowServiceImp) UnFollow(userId int64, targetId int64) (bool, error) {
	msg := strings.Builder{}
	msg.WriteString(strconv.Itoa(int(userId)))
	msg.WriteString(",")
	msg.WriteString(strconv.Itoa(int(targetId)))
	kafka.KfkFollowDel.FollowProducer(msg.String())
	return redisDao.UpdateRedisUnfollow(userId, targetId)
}

// GetFolloweeList 获取当前用户的关注列表
func (*FollowServiceImp) GetFolloweeList(userId int64) ([]*relation.User, error) {

	users := make([]*relation.User, 1)
	// 查询出错。
	if err := mysql.DB.Raw("select id,`name`,"+
		"\ncount(if(tag = 'follower' and cancel is not null,1,null)) follower_count,"+
		"\ncount(if(tag = 'follow' and cancel is not null,1,null)) follow_count,"+
		"\n 'true' is_follow\nfrom\n("+
		"\nselect f1.follower_id fid,u.id,`name`,f2.cancel,'follower' tag"+
		"\nfrom follows f1 join users u on f1.user_id = u.id and f1.cancel = 0"+
		"\nleft join follows f2 on u.id = f2.user_id and f2.cancel = 0\n\tunion all"+
		"\nselect f1.follower_id fid,u.id,`name`,f2.cancel,'follow' tag"+
		"\nfrom follows f1 join users u on f1.user_id = u.id and f1.cancel = 0"+
		"\nleft join follows f2 on u.id = f2.follower_id and f2.cancel = 0\n) T"+
		"\nwhere fid = ? group by fid,id,`name`", userId).Scan(&users).Error; nil != err {
		return nil, err
	}
	// 返回关注对象列表。
	return users, nil

}

// GetFollowerList 获取当前用户的粉丝列表
func (*FollowServiceImp) GetFollowerList(userId int64) ([]*relation.User, error) {
	users := make([]*relation.User, 1)

	err := mysql.DB.Raw("select T.id,T.name,T.follow_cnt follow_count,T.follower_cnt follower_count,if(f.cancel is null,'false','true') is_follow"+
		"\nfrom follows f right join"+
		"\n(select fid,id,`name`,"+
		"\ncount(if(tag = 'follower' and cancel is not null,1,null)) follower_cnt,"+
		"\ncount(if(tag = 'follow' and cancel is not null,1,null)) follow_cnt"+
		"\nfrom("+
		"\nselect f1.user_id fid,u.id,`name`,f2.cancel,'follower' tag"+
		"\nfrom follows f1 join users u on f1.follower_id = u.id and f1.cancel = 0"+
		"\nleft join follows f2 on u.id = f2.user_id and f2.cancel = 0"+
		"\nunion all"+
		"\nselect f1.user_id fid,u.id,`name`,f2.cancel,'follow' tag"+
		"\nfrom follows f1 join users u on f1.follower_id = u.id and f1.cancel = 0"+
		"\nleft join follows f2 on u.id = f2.follower_id and f2.cancel = 0"+
		"\n) T group by fid,id,`name`"+
		"\n) T on f.user_id = T.id and f.follower_id = T.fid and f.cancel = 0 where fid = ?", userId).
		Scan(&users).Error
	if err != nil {
		// 查询出错。
		return nil, err
	}

	// 查询成功。
	return users, nil
}

// GetFriendList 获取当前用户的好友列表
func (*FollowServiceImp) GetFriendList(userId int64) ([]*relation.User, error) {
	fu := make([]*relation.User, 1)
	//friends := make([]*relation.FriendUser, 1)

	err := mysql.DB.Raw("select T.id,T.name,T.follow_cnt follow_count,T.follower_cnt follower_count,if(f.cancel is null,'false','true') is_follow"+
		"\nfrom follows f right join"+
		"\n(select fid,id,`name`,"+
		"\ncount(if(tag = 'follower' and cancel is not null,1,null)) follower_cnt,"+
		"\ncount(if(tag = 'follow' and cancel is not null,1,null)) follow_cnt"+
		"\nfrom("+
		"\nselect f1.user_id fid,u.id,`name`,f2.cancel,'follower' tag"+
		"\nfrom follows f1 join users u on f1.follower_id = u.id and f1.cancel = 0"+
		"\nleft join follows f2 on u.id = f2.user_id and f2.cancel = 0"+
		"\nunion all"+
		"\nselect f1.user_id fid,u.id,`name`,f2.cancel,'follow' tag"+
		"\nfrom follows f1 join users u on f1.follower_id = u.id and f1.cancel = 0"+
		"\nleft join follows f2 on u.id = f2.follower_id and f2.cancel = 0"+
		"\n) T group by fid,id,`name`"+
		"\n) T on f.user_id = T.id and f.follower_id = T.fid and f.cancel = 0 where fid = ? and f.cancel is not null", userId).
		Scan(&fu).Error
	if err != nil {
		// 查询出错。
		return nil, err
	}
	/*
		for _, friendUser := range fu {
			friend := relation.FriendUser{
				User:    friendUser,
				Message: "",
				MsgType: 0,
			}
			friends = append(friends, &friend)
		}
	*/
	// 查询成功。
	return fu, nil
}
