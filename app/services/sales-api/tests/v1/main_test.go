package tests

import (
	"fmt"
	"testing"

	"github.com/colmmurphy91/go-service/business/data/dbtest"
	"github.com/colmmurphy91/go-service/foundation/docker"
)

var sqlC *docker.Container
var mongoC *docker.Container

func TestMain(m *testing.M) {
	var err error
	sqlC, err = dbtest.StartDB()
	mongoC, err = dbtest.StartMongo()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dbtest.StopDB(sqlC)
	defer dbtest.StopDB(mongoC)

	m.Run()
}
