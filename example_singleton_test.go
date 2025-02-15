package anchor_test

import (
	"fmt"
	"math/rand/v2"

	"github.com/kyuff/anchor"
)

func ExampleSingleton() {
	type Service struct {
		ID       int
		HostName string
	}

	service := anchor.Singleton(func() (*Service, error) {
		fmt.Println("Creating service")
		return &Service{
			ID:       rand.Int(),
			HostName: "localhost",
		}, nil
	})

	s := service()

	fmt.Printf("host: %s\n", s.HostName)
	fmt.Printf("id equal: %v\n", service().ID == s.ID)

	// Output:
	// Creating service
	// host: localhost
	// id equal: true
}
