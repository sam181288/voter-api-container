package api

import (
	"log"
	"net/http"
	"time"

	"drexel.edu/todo/db"
	"github.com/gofiber/fiber/v2"
)

// The api package creates and maintains a reference to the data handler
// this is a good design practice
type VoterAPI struct {
	db *db.VoterList
}

func New() (*VoterAPI, error) {
	dbHandler, err := db.NewVoterList()
	if err != nil {
		return nil, err
	}

	return &VoterAPI{db: dbHandler}, nil
}

//Below we implement the API functions.  Some of the framework
//things you will see include:
//   1) How to extract a parameter from the URL, for example
//	  the id parameter in /todo/:id
//   2) How to extract the body of a POST request
//   3) How to return JSON and a correctly formed HTTP status code
//	  for example, 200 for OK, 404 for not found, etc.  This is done
//	  using the c.JSON() function
//   4) How to return an error code and abort the request.  This is
//	  done using the c.AbortWithStatus() function

// implementation for GET /voters
// returns all todos
func (td *VoterAPI) ListAllVoters(c *fiber.Ctx) error {

	voterList, err := td.db.GetAllVoters()
	if err != nil {
		log.Println("Error Getting All Voters: ", err)
		return fiber.NewError(http.StatusNotFound,
			"Error Getting All Voters")
	}

	if voterList == nil {
		voterList = make([]db.Voter, 0)
		//todoList = make([]db.ToDoItem, 0)
	}

	return c.JSON(voterList)
}

// implementation for GET /voters/:id
// returns a single todo
func (td *VoterAPI) GetVoterByID(c *fiber.Ctx) error {

	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest)
	}
	//Convert the int to an uint
	idUint := uint(id)

	voter, err := td.db.GetVoter(idUint)
	if err != nil {
		log.Println("Item not found: ", err)
		return fiber.NewError(http.StatusNotFound)
	}

	return c.JSON(voter)
}

// implementation for POST /todo
// adds a new todo
func (td *VoterAPI) AddVoter(c *fiber.Ctx) error {
	var voter db.Voter

	if err := c.BodyParser(&voter); err != nil {
		log.Println("Error binding JSON: ", err)
		return fiber.NewError(http.StatusBadRequest)
	}

	if err := td.db.AddVoter(voter); err != nil {
		log.Println("Error adding item: ", err)
		return fiber.NewError(http.StatusInternalServerError)
	}

	return c.JSON(voter)
}

// immplementation for GET /voters/:id/polls - Gets the JUST the voter history for the voter with VoterID = :id
func (td *VoterAPI) GetVoterPolls(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest)
	}
	//Convert the int to an uint
	idUint := uint(id)

	voter, err := td.db.GetVoter(idUint)
	if err != nil {
		log.Println("Error getting voter: ", err)
		return fiber.NewError(http.StatusNotFound)
	}

	return c.JSON(voter.VoteHistory)
}

// implementation for GET /voters/:id/polls/:pollid - Gets JUST the single voter poll data with PollID = :id and VoterID = :id
func (td *VoterAPI) GetVoterPollById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest)
	}
	//Convert the int to an uint
	idUint := uint(id)

	pollid, err := c.ParamsInt("pollid")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest)
	}
	//Convert the int to an uint
	pollidUint := uint(pollid)

	voter, err := td.db.GetVoter(idUint)
	if err != nil {
		log.Println("Error getting voter: ", err)
		return fiber.NewError(http.StatusNotFound)
	}

	for _, poll := range voter.VoteHistory {
		if poll.PollId == pollidUint {
			return c.JSON(poll)
		}
	}

	return fiber.NewError(http.StatusNotFound)
}

// implementation for POST /voters/:id/polls/:pollid - adds one to the "database"
func (td *VoterAPI) AddVoterPollById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest)
	}
	//Convert the int to an uint
	idUint := uint(id)

	pollid, err := c.ParamsInt("pollid")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest)
	}
	//Convert the int to an uint
	pollidUint := uint(pollid)

	voter, err := td.db.GetVoter(idUint)
	if err != nil {
		log.Println("Error getting voter: ", err)
		return fiber.NewError(http.StatusNotFound)
	}

	voter.VoteHistory = append(voter.VoteHistory, db.VoterHistory{PollId: pollidUint, VoteId: idUint, VoteDate: time.Now()})

	if err := td.db.UpdateVoter(voter); err != nil {
		log.Println("Error updating voter: ", err)
		return fiber.NewError(http.StatusInternalServerError)
	}

	return c.JSON(voter)
}

// implementation for DELETE /voters
func (td *VoterAPI) DeleteAllVoters(c *fiber.Ctx) error {
	err := td.db.DeleteAllVoters()
	if err != nil {
		log.Println("Error deleting all voters: ", err)
		return fiber.NewError(http.StatusInternalServerError)
	}

	return c.SendStatus(http.StatusOK)
}

// implementation got PUT /voters/:id
func (td *VoterAPI) UpdateVoter(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest)
	}
	//Convert the int to an uint
	idUint := uint(id)
	voter, err := td.db.GetVoter(idUint)
	if err != nil {
		log.Println("Error getting voter: ", err)
		return fiber.NewError(http.StatusNotFound)
	}
	//update the voter
	if err := c.BodyParser(&voter); err != nil {
		log.Println("Error binding JSON: ", err)
		return fiber.NewError(http.StatusBadRequest)
	}
	return c.SendStatus(http.StatusOK)
}

// implementation got PUT /voters/:id/polls/:pollId
func (td *VoterAPI) UpdateVoterPollById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest)
	}
	pollId, err := c.ParamsInt("pollId")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest)
	}
	//Convert the int to an uint
	idUint := uint(id)
	pollIdUint := uint(pollId)
	voter, err := td.db.GetVoter(idUint)
	if err != nil {
		log.Println("Error getting voter: ", err)
		return fiber.NewError(http.StatusNotFound)
	}
	//update the voter poll
	for i, poll := range voter.VoteHistory {
		if poll.PollId == pollIdUint {
			if err := c.BodyParser(&voter.VoteHistory[i]); err != nil {
				log.Println("Error binding JSON: ", err)
				return fiber.NewError(http.StatusBadRequest)
			}
			break
		}
	}
	return c.SendStatus(http.StatusOK)
}

// implementation of GET /health. It is a good practice to build in a
// health check for your API.  Below the results are just hard coded
// but in a real API you can provide detailed information about the
// health of your API with a Health Check
func (td *VoterAPI) HealthCheck(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).
		JSON(fiber.Map{
			"status":             "ok",
			"version":            "1.0.0",
			"uptime":             time.Since(time.Now()),
			"users_processed":    1000,
			"errors_encountered": 10,
		})
}
