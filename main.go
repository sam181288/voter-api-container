package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"drexel.edu/todo/api"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var (
	hostFlag string
	portFlag uint
)

func processCmdLineFlags() {

	flag.StringVar(&hostFlag, "h", "0.0.0.0", "Listen on all interfaces")
	flag.UintVar(&portFlag, "p", 1080, "Default Port")

	flag.Parse()
}

func main() {
	processCmdLineFlags()

	app := fiber.New()
	app.Use(cors.New())
	app.Use(recover.New())

	apiHandler, err := api.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	app.Get("/voters", apiHandler.ListAllVoters)
	app.Post("/voters/:id<int>", apiHandler.AddVoter)
	app.Get("/voters/:id<int>", apiHandler.GetVoterByID)
	app.Get("/voters/:id<int>/polls", apiHandler.GetVoterPolls)
	app.Get("/voters/:id<int>/polls/:pollid<int>", apiHandler.GetVoterPollById)
	app.Post("/voters/:id<int>/polls/:pollid<int>", apiHandler.AddVoterPollById)
	app.Put("/voters/:id<int>", apiHandler.UpdateVoter)
	app.Put("/voters/:id<int>/polls/:pollid<int>", apiHandler.UpdateVoterPollById)
	app.Delete("/voters", apiHandler.DeleteAllVoters)
	app.Get("/voters/health", apiHandler.HealthCheck)

	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	log.Println("Starting server on ", serverPath)
	app.Listen(serverPath)
}
