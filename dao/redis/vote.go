package redis

import (
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 //每一票值多少分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated   = errors.New("不允许重复投票")
)

func CreatPost(postID, communityID int64) error {
	// key := getRedisKey(KeyPostTimeZSet)
	// // 检查键的类型，确保它是有序集合
	// keyType, err := rdb.Type(key).Result()
	// if err != nil {
	// 	return fmt.Errorf("failed to get Redis key type: %v", err)
	// }
	// if keyType != "zset" && keyType != "none" { // 仅当键类型不是 zset 且键存在时报错
	// 	return fmt.Errorf("Redis key %s type is %s, expected zset", key, keyType)
	// }
	pipline := rdb.TxPipeline()
	//帖子时间
	pipline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()), // 当前时间戳作为分数
		Member: postID,                     // 帖子ID作为成员
	}).Result()

	// 帖子分数
	pipline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()), // 当前时间戳作为分数
		Member: postID,                     // 帖子ID作为成员
	}).Result()
	//把帖子id加到社区的set
	cKey := getRedisKey(KeyCommunitySetPrefix + strconv.Itoa(int(communityID)))
	pipline.SAdd(cKey, postID)
	_, err := pipline.Exec()
	return err
}

func VoteForPost(userID, postID string, value float64) error {
	//判断投票限制
	//去redis取帖子发布时间
	postTime := rdb.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	//更新分数
	//先查当前用户给当前帖子之前的投票记录
	ovalue := rdb.ZScore(getRedisKey(KeyPostVotedPrefix+postID), userID).Val()
	//如果这一次投票的值和之前保存的值一致，就提示不允许重复投票
	if value == ovalue {
		return ErrVoteRepeated
	}
	var op float64
	if value > ovalue {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ovalue - value) //计算两次投票的差值

	pipline := rdb.TxPipeline()
	pipline.ZIncrBy(getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID).Result()

	//记录用户为该帖子投票的数据
	if value == 0 {
		pipline.ZRem(getRedisKey(KeyPostVotedPrefix+postID), postID).Result()
	} else {
		pipline.ZAdd(getRedisKey(KeyPostVotedPrefix+postID), redis.Z{
			Score:  value,
			Member: userID,
		}).Result()
	}
	_, err := pipline.Exec()
	return err
}
