package redlock

import (
	"fmt"
	"os"
	"time"

	redis "github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"

	"github.com/werunclub/baymax/v2/log"
)

var client *redis.Client

func DefaultClient() *redis.Client {
	return client
}

func Connect(address, password string, db int) {
	client = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		log.SourcedLogrus().WithError(err).WithField("Pong", pong).Warning("redis 连接失败")
		os.Exit(1)
	}
}

func Close() error {
	return client.Close()
}

func SetRedisClient(cli *redis.Client) {
	client = cli
}

// Deprecated: GetLock 使用 redis INCR 模拟分布式锁,
// :return (1, true) 时获得锁, 其它情况获取锁失败
func GetLock(key string, releaseTime time.Duration) (int64, bool) {
	rtn := client.Incr(key)
	if rtn.Err() != nil {
		logrus.WithError(rtn.Err()).Error("获取 redis 分布式锁发生错误")
		return 0, false
	}

	// 获得锁, 设置锁过期时间
	if rtn.Val() == 1 {
		client.Expire(key, releaseTime)
		return 1, true
	}

	return rtn.Val(), false
}

func ReleaseLock(key string) {
	client.Del(key)
}

type Lock struct {
	Key      string
	Value    int64
	Error    error
	LifeTime time.Duration
}

func (l *Lock) Assign() (bool, error) {
	rtn := client.Incr(l.Key)

	if err := rtn.Err(); err != nil {
		l.Error = err
		return false, l.Error
	}

	l.Value = rtn.Val()

	if l.Value == 1 {
		l.setExpire()
		return true, nil
	}

	if l.Value != 1 {
		l.Error = fmt.Errorf("assign lock faild with number [%d]", l.Value)
	}

	return false, l.Error
}

func (l *Lock) setExpire() {
	client.Expire(l.Key, l.LifeTime)
}

func (l *Lock) Release() {
	client.Del(l.Key)
}

func NewLock(key string, duration ...time.Duration) Lock {
	lifeTime := 1 * time.Minute
	if len(duration) > 1 {
		lifeTime = duration[0]
	}

	lock := Lock{Key: key, LifeTime: lifeTime}

	return lock
}

func AssignLock(key string, duration time.Duration) (bool, error) {
	lock := NewLock(key, duration)
	return lock.Assign()
}

// 通知处理时序
type Sequence struct {
	Key      string
	Value    int64
	Error    error
	LifeTime time.Duration
}

func NewSequence(key string) *Sequence {
	return &Sequence{Key: key, LifeTime: 48 * time.Hour}
}

func (s *Sequence) Get() int64 {
	val, err := client.Get(s.Key).Int64()
	s.Value = val
	s.Error = err

	if val == 0 && err == nil {
		s.Set(0)
		client.Expire(s.Key, s.LifeTime)
	}

	return val
}

func (s *Sequence) Set(val int64) {
	err := client.Set(s.Key, val, s.LifeTime).Err()
	s.Error = err
}
