package redis

const (
	KeyPrefix             = "bluebell:"
	KeyPostTimeZSet       = "post:time"   //zset;帖子及发帖时间
	KeyPostScoreZSet      = "post:score"  //zset;帖子及投票分数
	KeyPostVotedPrefix    = "post:voted:" //zset;记录用户及投票的类型；参数是post_id
	KeyCommunitySetPrefix = "community:"  //set:保存每个分区下帖子的id
)

func getRedisKey(key string) string {
	return KeyPrefix + key
}
