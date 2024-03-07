package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	RedisNilError        = "redis: nil"
	RedisDefaultLocation = "0.0.0.0:6379"
	RedisKeyPrefix       = "voter:"
)

type cache struct {
	client  *redis.Client
	context context.Context
}

type VoterHistory struct {
	PollId   uint
	VoteId   uint
	VoteDate time.Time
}

type Voter struct {
	VoterId     uint           `json:"id"`
	Name        string         `json:"name"`
	Email       string         `json:"email"`
	VoteHistory []VoterHistory `json:"history"`
}

type VoterMap map[uint]Voter

type VoterList struct {
	Voters VoterMap //A map of VoterIDs as keys and Voter structs as values

	cache
}

func NewVoterList() (*VoterList, error) {
	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		redisUrl = RedisDefaultLocation
	}

	return NewWithCacheInstance(redisUrl)
}

func NewWithCacheInstance(location string) (*VoterList, error) {

	client := redis.NewClient(&redis.Options{
		Addr: location,
	})

	//We use this context to coordinate betwen our go code and
	//the redis operaitons
	ctx := context.TODO()

	//This is the reccomended way to ensure that our redis connection
	//is working
	err := client.Ping(ctx).Err()
	if err != nil {
		log.Println("Error connecting to redis" + err.Error())
		return nil, err
	}

	//Return a pointer to a new ToDo struct
	return &VoterList{
		cache: cache{
			client:  client,
			context: ctx,
		},
	}, nil
}

func isRedisNilError(err error) bool {
	return errors.Is(err, redis.Nil) || err.Error() == RedisNilError
}

func redisKeyFromId(id int) string {
	return fmt.Sprintf("%s%d", RedisKeyPrefix, id)
}

func (vl *VoterList) doesKeyExist(id int) bool {
	kc, _ := vl.client.Exists(vl.context, redisKeyFromId(id)).Result()
	return kc > 0
}

func fromJsonString(s string, voter *Voter) error {
	err := json.Unmarshal([]byte(s), &voter)
	if err != nil {
		return err
	}
	return nil
}

func (vl *VoterList) AddVoter(voter Voter) error {
	//check if the voter exists
	if vl.doesKeyExist(int(voter.VoterId)) {
		return fmt.Errorf("Voter with id %d already exists", voter.VoterId)
	}
	//add voter to redis
	return vl.client.JSONSet(vl.context, redisKeyFromId(int(voter.VoterId)), ".", voter).Err()
}

func (vl *VoterList) UpdateVoter(voter Voter) error {
	//check if the voter exists
	if !vl.doesKeyExist(int(voter.VoterId)) {
		return fmt.Errorf("Voter with id %d doen't exists", voter.VoterId)
	}
	//update voter in redis
	return vl.client.JSONSet(vl.context, redisKeyFromId(int(voter.VoterId)), ".", voter).Err()
}

func (vl *VoterList) GetVoter(voterId uint) (*Voter, error) {
	newVoter := &Voter{}
	err := vl.getVoterFromRedis(redisKeyFromId(int(voterId)), newVoter)
	if err != nil {
		return nil, err
	}
	return newVoter, nil
}

func (vl *VoterList) GetAllVoters() ([]Voter, error) {
	//get all voters from redis cache
	voters := make([]Voter, 0)
	keys, err := vl.client.Keys(vl.context, RedisKeyPrefix+"*").Result()
	if err != nil {
		return nil, err
	}
	for _, key := range keys {
		newVoter := &Voter{}
		err := vl.getVoterFromRedis(key, newVoter)
		if err != nil {
			return nil, err
		}
		voters = append(voters, *newVoter)
	}
	return voters, nil
}

func (vl *VoterList) getVoterFromRedis(key string, voter *Voter) error {
	voterJson, err := vl.client.JSONGet(vl.context, key, ".").Result()
	if err != nil {
		return err
	}
	return fromJsonString(voterJson, voter)
}

func (vl *VoterList) DeleteVoter(voterId uint) error {
	//delete voter from redis
	voter := Voter{}
	err := vl.client.Get(vl.context, redisKeyFromId(int(voterId))).Scan(&voter)
	if err != nil {
		if isRedisNilError(err) {
			return fmt.Errorf("Voter with id %d doesn't exist", voterId)
		}
		return err
	}
	//delete voter from redis
	vl.client.Del(vl.context, redisKeyFromId(int(voterId)))
	return nil
}

func (vl *VoterList) DeleteAllVoters() error {
	//delete all voters from redis cache
	vl.client.FlushAll(vl.context)
	return nil
}
