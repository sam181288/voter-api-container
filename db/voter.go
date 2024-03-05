package db

import (
	"errors"
	"time"
)

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
}

func NewVoterList() (*VoterList, error) {

	voterList := &VoterList{
		Voters: make(map[uint]Voter),
	}

	return voterList, nil
}

func (vl *VoterList) AddVoter(voter Voter) error {

	_, ok := vl.Voters[voter.VoterId]
	if ok {
		return errors.New("voter already exists")
	}

	vl.Voters[voter.VoterId] = voter
	return nil
}

func (vl *VoterList) UpdateVoter(voter Voter) error {

	_, ok := vl.Voters[voter.VoterId]
	if !ok {
		return errors.New("voter does not exist")
	}

	vl.Voters[voter.VoterId] = voter
	return nil
}

func (vl *VoterList) GetVoter(voterId uint) (Voter, error) {

	voter, ok := vl.Voters[voterId]
	if !ok {
		return Voter{}, errors.New("voter does not exist")
	}

	return voter, nil
}

func (vl *VoterList) DeleteVoter(voterId uint) error {

	_, ok := vl.Voters[voterId]
	if !ok {
		return errors.New("voter does not exist")
	}

	delete(vl.Voters, voterId)
	return nil
}

func (vl *VoterList) GetAllVoters() ([]Voter, error) {

	voters := make([]Voter, 0, len(vl.Voters))
	for _, voter := range vl.Voters {
		voters = append(voters, voter)
	}

	return voters, nil
}

func (vl *VoterList) DeleteAllVoters() error {
	vl.Voters = make(map[uint]Voter)
	return nil
}
