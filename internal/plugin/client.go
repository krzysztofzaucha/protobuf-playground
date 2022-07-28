package main

import (
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/krzysztofzaucha/protobuf-playground/internal"
	"github.com/krzysztofzaucha/protobuf-playground/internal/model"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"time"
)

// Client is a client plugin symbol name.
var Client client

var errClientPlugin = errors.New("client")

type client struct {
	config *internal.Config
}

func (c *client) WithConfig(config *internal.Config) {
	c.config = config
}

// Execute method executes plugin logic.
func (c *client) Execute() error {
	fmt.Println("starting client")

	cc, err := grpc.Dial(fmt.Sprintf("%s:%d", c.config.Server.Host, c.config.Server.Port), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	// Close connection before exiting application
	defer cc.Close()

	cl := model.NewPersonServiceClient(cc)

	var cres *model.CreatePersonResponse

	for i := 0; i < 3; i++ {
		cres, err = cl.CreatePerson(context.Background(), &model.CreatePersonRequest{
			Person: &model.Person{
				Name: "John Smith",
				Email: []string{
					"john.smith@example.com",
				},
				Address: []*model.Address{
					{
						Town:    "London",
						City:    "London",
						Country: "United Kingdom",
					},
				},
			},
		})

		if i < 3 {
			time.Sleep(time.Duration(1) * time.Second)
			continue
		}

		if err != nil {
			panic(err)
		}
	}

	person, err := json.Marshal(cres.GetPerson())
	if err != nil {
		panic(err)
	}

	fmt.Println(string(person))

	rres, err := cl.ReadPerson(context.Background(), &model.ReadPersonRequest{
		ID: 1,
	})

	person, err = json.Marshal(rres.GetPerson())
	if err != nil {
		panic(err)
	}

	fmt.Println(string(person))

	return nil
}
