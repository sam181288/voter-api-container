package tests

import (
	"log"
	"os"
	"strconv"
	"testing"

	"drexel.edu/todo/db"
	fake "github.com/brianvoe/gofakeit/v6" //aliasing package name
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

var (
	BASE_API = "http://localhost:1080"

	cli = resty.New()
)

func TestMain(m *testing.M) {

	//SETUP GOES FIRST
	rsp, err := cli.R().Delete(BASE_API + "/voters")

	if rsp.StatusCode() != 200 {
		log.Printf("error clearing database, %v", err)
		os.Exit(1)
	}

	code := m.Run()

	os.Exit(code)
}

func newRandVoter(id uint) db.Voter {
	return db.Voter{
		VoterId: id,
		Name:    fake.Name(),
		Email:   fake.Email(),
	}
}

func newRandPollForVoter(id uint) db.VoterHistory {
	return db.VoterHistory{
		PollId:   id,
		VoteId:   id,
		VoteDate: fake.Date(),
	}
}

func Test_LoadDB(t *testing.T) {
	numLoad := 3
	for i := 0; i < numLoad; i++ {

		item := newRandVoter(uint(i))
		rsp, err := cli.R().
			SetBody(item).
			Post(BASE_API + "/voters/" + strconv.Itoa(i))

		assert.Nil(t, err)
		assert.Equal(t, 200, rsp.StatusCode())
	}
}

func Test_GetAllVoters(t *testing.T) {
	var items []db.Voter

	rsp, err := cli.R().SetResult(&items).Get(BASE_API + "/voters")

	assert.Nil(t, err)
	assert.Equal(t, 200, rsp.StatusCode())

	assert.Equal(t, 3, len(items))
}

func Test_GetVoterByID(t *testing.T) {
	var item db.Voter

	rsp, err := cli.R().SetResult(&item).Get(BASE_API + "/voters/1")

	assert.Nil(t, err)
	assert.Equal(t, 200, rsp.StatusCode())

}

func Test_ApiHealth(t *testing.T) {
	rsp, err := cli.R().Get(BASE_API + "/voters/health")
	assert.Nil(t, err)
	assert.Equal(t, 200, rsp.StatusCode())
}

func Test_PostVoterPoll(t *testing.T) {
	//add a poll to voter 1
	numLoad := 3
	for i := 0; i < numLoad; i++ {
		pollHistory := newRandPollForVoter(uint(i))
		rsp, err := cli.R().
			SetBody(pollHistory).
			Post(BASE_API + "/voters/1/polls/" + strconv.Itoa(i))
		assert.Nil(t, err)
		assert.Equal(t, 200, rsp.StatusCode())
	}
}

func Test_GetVoterPolls(t *testing.T) {
	var items []db.VoterHistory

	rsp, err := cli.R().SetResult(&items).Get(BASE_API + "/voters/1/polls")

	assert.Nil(t, err)
	assert.Equal(t, 200, rsp.StatusCode())

	assert.Equal(t, 3, len(items))
}

func Test_GetVoterPollByID(t *testing.T) {
	var item db.VoterHistory

	rsp, err := cli.R().SetResult(&item).Get(BASE_API + "/voters/1/polls/1")

	assert.Nil(t, err)
	assert.Equal(t, 200, rsp.StatusCode())

}

// func Test_UpdateVoter(t *testing.T) {
// 	item := newRandVoter(1)
// 	rsp, err := cli.R().
// 		SetBody(item).
// 		Put(BASE_API + "/voters/1")

// 	assert.Nil(t, err)
// 	assert.Equal(t, 200, rsp.StatusCode())
// }

// func Test_UpdateVoterPollByID(t *testing.T) {
// 	item := newRandPollForVoter(1)
// 	rsp, err := cli.R().
// 		SetBody(item).
// 		Put(BASE_API + "/voters/1/polls/1")

// 	assert.Nil(t, err)
// 	assert.Equal(t, 200, rsp.StatusCode())
// }

// func Test_DeleteAllVoters(t *testing.T) {
// 	rsp, err := cli.R().Delete(BASE_API + "/voters")
// 	assert.Nil(t, err)
// 	assert.Equal(t, 200, rsp.StatusCode())

// 	rsp, err = cli.R().Get(BASE_API + "/voters")
// 	assert.Nil(t, err)
// 	assert.Equal(t, 200, rsp.StatusCode())

// 	var items []db.Voter
// 	assert.Equal(t, 0, len(items))
// }
