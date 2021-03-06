package sessions

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
)

//RedisStore represents a session.Store backed by redis.
type RedisStore struct {
	//Redis client used to talk to redis server.
	Client *redis.Client
	//Used for key expiry time on redis.
	SessionDuration time.Duration
}

//NewRedisStore constructs and returns a new RedisStore
func NewRedisStore(client *redis.Client, sessionDuration time.Duration) *RedisStore {
	//initialize and return a new RedisStore struct
	store := RedisStore{
		Client:          client,
		SessionDuration: sessionDuration,
	}
	return &store
}

//Store implementation

//Save saves the provided `sessionState` and associated SessionID to the store.
//The `sessionState` parameter is typically a pointer to a struct containing
//all the data you want to associated with the given SessionID.
func (rs *RedisStore) Save(sid SessionID, sessionState interface{}) error {
	jsonSession, err := json.Marshal(sessionState)
	if err != nil {
		return err
	}
	rdb := rs.Client
	err2 := rdb.Set(sid.getRedisKey(), jsonSession, rs.SessionDuration).Err()
	if err2 != nil {
		return err2
	}
	return nil
}

//Get populates `sessionState` with the data previously saved
//for the given SessionID
func (rs *RedisStore) Get(sid SessionID, sessionState interface{}) error {
	session, err := rs.Client.Get(sid.getRedisKey()).Result()
	rs.Client.Get(sid.getRedisKey()).Err()
	if len(session) == 0 {
		return ErrStateNotFound
	}
	if err != nil {
		return err
	} else {
		json.Unmarshal([]byte(session), sessionState)
		rs.Client.Expire(sid.getRedisKey(), rs.SessionDuration)
	}
	return nil
}

//Delete deletes all state data associated with the SessionID from the store.
func (rs *RedisStore) Delete(sid SessionID) error {
	rs.Client.Del(sid.getRedisKey())
	return nil
}

//getRedisKey() returns the redis key to use for the SessionID
func (sid SessionID) getRedisKey() string {
	return "sid:" + sid.String()
}
