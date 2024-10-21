package redis

import (
	"bluebell/models"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func GetIDsFromKey(key string, page, size int64) ([]string, error) {
	//确定查询的索引起始点
	start := (page - 1) * size
	end := start + size - 1
	//按分数从大到小查询指定数量的元素
	return rdb.ZRevRange(key, start, end).Result()
}

func GetPostIDsByOrder(p *models.ParamPostList) ([]string, error) {
	//从redis获取id
	//根据用户请求中携带的order参数确定要查询的key
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	return GetIDsFromKey(key, p.Page, p.Size)
}

// 根据ids查询每篇帖子投赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	// data = make([]int64, 0, len(ids))
	// for _, id := range ids {
	// 	key := getRedisKey(KeyPostVotedPrefix + id)
	// 	//查找key中分数是1的元素的数量(统计每篇帖子的赞成票的数量)
	// 	v := rdb.ZCount(key, "1", "1").Val()
	// 	data = append(data, v)
	// }
	pipline := rdb.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedPrefix + id)
		pipline.ZCount(key, "1", "1")
	}
	cmders, err := pipline.Exec()
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

// 按社区查询ids
func GetCommunityPostIDsByOrder(p *models.ParamPostList) ([]string, error) {
	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}

	//使用zinterstore把分社区的帖子set和帖子分数的zset取交集生成一个新的zset
	cKey := getRedisKey(KeyCommunitySetPrefix + strconv.Itoa(int(p.CommunityID)))
	//利用缓存key减少zinterstore执行的次数
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	if rdb.Exists(orderKey).Val() < 1 {
		//不存在，需要计算
		pipline := rdb.Pipeline()
		pipline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, cKey, orderKey)
		pipline.Expire(key, 60*time.Second) //设置超时时间
		_, err := pipline.Exec()
		if err != nil {
			return nil, err
		}
	}
	//存在的话就根据key查询ids
	return GetIDsFromKey(key, p.Page, p.Size)
}
