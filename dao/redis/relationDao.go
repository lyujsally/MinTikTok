package redis

import (
	"strconv"
	"time"

	"github.com/lyujsally/MinTikTok-lyujsally/conf"

	"github.com/go-redis/redis"
)

const (
	ACTION_TYPE_USER = 1
)

type FollowerValue struct {
	UserId   int64
	IsFollow bool
}

/*
func Follow(userId int64, entityType int, entityId int64) error {
	followeeKey := GetFolloweeKey(userId, entityType)
	followerKey := GetFollowerKey(entityId, entityType)
	now := time.Now().Unix()

	//执行Redis事务
	pipe := RDB.TxPipeline()

	//向事务添加命令
	pipe.ZAdd(followeeKey, redis.Z{
		Score:  float64(now),
		Member: strconv.FormatInt(entityId, 10),
	})
	isFo := IsFollow(entityId, entityType, userId)
	value := FollowerValue{
		UserId:   userId,
		IsFollow: isFo,
	}
	valueJSON, err := json.Marshal(value)
	if err != nil {
		log.Println("valueJSON Marshal failed")
		return err
	}
	pipe.ZAdd(followerKey, redis.Z{
		Score:  float64(now),
		Member: string(valueJSON), /*strconv.FormatInt(userId, 10),
	})

	//执行事务命令
	result, err := pipe.Exec()

	if err != nil {
		log.Println("follow pipe exec failed")
		return err
	}

	fmt.Println(result) // 打印事务执行结果，根据需要进行处理

	return nil
}

func UnFollow(userId int64, entityType int, entityId int64) error {
	followeeKey := GetFolloweeKey(userId, entityType)
	followerKey := GetFollowerKey(entityId, entityType)

	//执行Redis事务
	pipe := RDB.TxPipeline()

	//向事务添加命令
	pipe.ZRem(followeeKey, entityId)
	isFo := IsFollow(entityId, entityType, userId)
	value := FollowerValue{
		UserId:   userId,
		IsFollow: isFo,
	}
	valueJSON, err := json.Marshal(value)
	if err != nil {
		log.Println("valueJSON Marshal failed")
		return err
	}
	pipe.ZRem(followerKey, string(valueJSON))

	//执行事务命令
	result, err := pipe.Exec()

	if err != nil {
		log.Println("unfollow pipe exec failed")
		return err
	}

	fmt.Println(result) // 打印事务执行结果，根据需要进行处理

	return nil
}*/

// 将isfollow信息注入redis
func AddRelationToRedis(userId int64, targetId int64) {

	now := time.Now().Unix()
	// 第一次存入时，给该key添加一个-1为key，防止脏数据的写入。当然set可以去重，直接加，便于CPU。
	RdbIsFollow.ZAdd(strconv.FormatInt(userId, 10), redis.Z{
		Score:  float64(now),
		Member: "-1",
	})
	// 将查询到的关注关系注入Redis
	RdbIsFollow.ZAdd(strconv.FormatInt(userId, 10), redis.Z{
		Score:  float64(now),
		Member: strconv.FormatInt(targetId, 10),
	})
	// 更新过期时间。
	RdbIsFollow.Expire(strconv.FormatInt(userId, 10), conf.ExpireTime)
}

// 将Followers信息注入redis
func AddFollowersToRedis(userId int64, ids []int64) {
	now := time.Now().Unix()
	RdbFollower.ZAdd(strconv.FormatInt(userId, 10), redis.Z{
		Score:  float64(now),
		Member: "-1",
	})
	for i, id := range ids {
		RdbFollower.ZAdd(strconv.FormatInt(userId, 10), redis.Z{
			Score:  float64(now),
			Member: strconv.FormatInt(id, 10),
		})
		RdbIsFollow.ZAdd(strconv.FormatInt(id, 10), redis.Z{
			Score:  float64(now),
			Member: strconv.FormatInt(userId, 10),
		})
		RdbIsFollow.ZAdd(strconv.FormatInt(id, 10), redis.Z{
			Score:  float64(now),
			Member: "-1",
		})
		// 更新部分关注者的时间
		RdbIsFollow.Expire(strconv.FormatInt(id, 10), conf.ExpireTime+time.Duration((i%10)<<8))
	}
	// 更新followers的过期时间。
	RdbFollower.Expire(strconv.FormatInt(userId, 10), conf.ExpireTime)

}

// 将Followees信息注入redis
func AddFolloweesToRedis(userId int64, ids []int64) {
	now := time.Now().Unix()
	RdbFollowee.ZAdd(strconv.FormatInt(userId, 10), redis.Z{
		Score:  float64(now),
		Member: "-1",
	})
	for i, id := range ids {
		RdbFollowee.ZAdd(strconv.FormatInt(userId, 10), redis.Z{
			Score:  float64(now),
			Member: strconv.FormatInt(id, 10),
		})
		RdbIsFollow.ZAdd(strconv.FormatInt(userId, 10), redis.Z{
			Score:  float64(now),
			Member: strconv.FormatInt(id, 10),
		})
		RdbIsFollow.ZAdd(strconv.FormatInt(userId, 10), redis.Z{
			Score:  float64(now),
			Member: "-1",
		})
		// 更新部分关注者的时间
		RdbIsFollow.Expire(strconv.FormatInt(userId, 10), conf.ExpireTime+time.Duration((i%10)<<8))
	}
	// 更新followees的过期时间。
	RdbFollowee.Expire(strconv.FormatInt(userId, 10), conf.ExpireTime)

}

// 添加Redis里当前用户关注目标用户的信息
func UpdateRedisFollow(userId int64, targetId int64) (bool, error) {

	userIdStr := strconv.FormatInt(userId, 10)
	targetIdStr := strconv.FormatInt(targetId, 10)
	now := time.Now().Unix()
	// 当前targetid是否作为键在redis粉丝列表存在
	if cnt, _ := RdbFollower.ZCard(targetIdStr).Result(); cnt != 0 {
		RdbFollower.ZAdd(targetIdStr, redis.Z{
			Score:  float64(now),
			Member: userIdStr,
		})
		RdbFollower.Expire(targetIdStr, conf.ExpireTime)
	}

	// 当前userId是否作为键在redis关注列表存在
	if cnt, _ := RdbFollowee.ZCard(userIdStr).Result(); cnt != 0 {
		RdbFollowee.ZAdd(userIdStr, redis.Z{
			Score:  float64(now),
			Member: targetIdStr,
		})
		RdbFollowee.Expire(userIdStr, conf.ExpireTime)
	}

	// 加入IsFollow列表
	RdbIsFollow.ZAdd(userIdStr, redis.Z{
		Score:  float64(now),
		Member: targetIdStr,
	})
	RdbIsFollow.ZAdd(userIdStr, redis.Z{
		Score:  float64(now),
		Member: "-1",
	})
	RdbIsFollow.Expire(userIdStr, conf.ExpireTime)

	return true, nil

}

func UpdateRedisUnfollow(userId int64, targetId int64) (bool, error) {

	userIdStr := strconv.FormatInt(userId, 10)
	targetIdStr := strconv.FormatInt(targetId, 10)
	// 当前targetid是否作为键在redis粉丝列表存在
	if cnt, _ := RdbFollower.ZCard(targetIdStr).Result(); cnt != 0 {
		RdbFollower.ZRem(targetIdStr, userIdStr)
		RdbFollower.Expire(targetIdStr, conf.ExpireTime)
	}

	// 当前userId是否作为键在redis关注列表存在
	if cnt, _ := RdbFollowee.ZCard(userIdStr).Result(); cnt != 0 {
		RdbFollowee.ZRem(userIdStr, targetIdStr)
		RdbFollowee.Expire(userIdStr, conf.ExpireTime)
	}

	// 加入IsFollow列表
	if cnt, _ := RdbIsFollow.ZCard(userIdStr).Result(); cnt != 0 {
		RdbIsFollow.ZRem(userIdStr, targetIdStr)
		RdbIsFollow.Expire(userIdStr, conf.ExpireTime)
	}

	return true, nil

}

// 再次删除Redis里的信息，防止脏数据，保证最终一致性
func ReUpdateRedisUnfollow(userId int64, targetId int64) {

	userIdStr := strconv.FormatInt(userId, 10)
	targetIdStr := strconv.FormatInt(targetId, 10)
	// 当前targetid是否作为键在redis粉丝列表存在
	if cnt, _ := RdbFollower.ZCard(targetIdStr).Result(); cnt != 0 {
		RdbFollower.ZRem(targetIdStr, userIdStr)
	}

	// 当前userId是否作为键在redis关注列表存在
	if cnt, _ := RdbFollowee.ZCard(userIdStr).Result(); cnt != 0 {
		RdbFollowee.ZRem(userIdStr, targetIdStr)
	}

	// 加入IsFollow列表
	if cnt, _ := RdbIsFollow.ZCard(userIdStr).Result(); cnt != 0 {
		RdbIsFollow.ZRem(userIdStr, targetIdStr)
	}
}
