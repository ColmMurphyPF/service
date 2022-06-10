package cafe

import (
	"errors"
	"fmt"
	"github.com/colmmurphy91/go-service/business/data/dbtest"
	"github.com/colmmurphy91/go-service/foundation/docker"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

var c *docker.Container

func TestMain(m *testing.M) {
	var err error
	c, err = dbtest.StartMongo()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dbtest.StopDB(c)

	m.Run()
}

func TestCafe(t *testing.T) {
	log, db, teardown := dbtest.NewMongoUnit(t, c, "testcafe")
	t.Cleanup(teardown)

	core := NewCore(log, db)

	t.Log("Given the need to work with Product records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single Product.", testID)
		{

			nc := NewCafe{
				Name:        "Spinneys Cafe",
				Address:     "Location",
				PhoneNumber: "PhoneNumber",
			}

			caf, err := core.CreateCafe(nc)

			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create a cafe : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create a cafe.", dbtest.Success, testID)

			bc := NewCafe{}
			_, err = core.CreateCafe(bc)

			if err == nil {
				t.Fatalf("\t%s\tTest %d:\tShould not be able to have empty fields: %s.", dbtest.Failed, testID, err)
			}

			t.Logf("\t%s\tTest %d:\tShould not  be able to create a bad cafe.", dbtest.Success, testID)

			foundCafe, err := core.FindCafe(caf.ID)

			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to find a cafe : %s.", dbtest.Failed, testID, err)
			}

			if foundCafe.Name != nc.Name {
				t.Fatalf("\t%s\tTest %d:\tnames should be equal: %s.", dbtest.Failed, testID, err)
			}

			t.Logf("\t%s\tTest %d:\tShould be able to find a cafe.", dbtest.Success, testID)

			cafes, err := core.FindAll()

			if len(cafes) == 0 {
				t.Fatalf("\t%s\tTest %d:\tShould not have found cafe, should be deleted: %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to all cafes.", dbtest.Success, testID)

			_, err = core.FindCafe("2")

			if !errors.Is(err, ErrNotFound) {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create find a cafe : %s.", dbtest.Failed, testID, err)
			}

			t.Logf("\t%s\tTest %d:\tShould throw error when cafe does not exist.", dbtest.Success, testID)

			caf.Name = "new name"
			updateCafe := UpdateCafe{
				ID:   caf.ID,
				Name: &caf.Name,
			}
			err = core.UpdateCafe(updateCafe)

			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to update a cafes name : %s.", dbtest.Failed, testID, err)
			}

			t.Logf("\t%s\tTest %d:\tShould throw error when cafe does not exist.", dbtest.Success, testID)

			caf2, err := core.FindCafe(caf.ID)

			if caf2.Name != "new name" {
				t.Fatalf("\t%s\tTest %d:\tSname hsould be updated : %s.", dbtest.Failed, testID, err)
			}

			caf.ID = primitive.NewObjectID().Hex()
			err = core.UpdateCafe(UpdateCafe{ID: primitive.NewObjectID().Hex()})

			if !errors.Is(err, ErrNotFound) {
				t.Fatalf("\t%s\tTest %d:\tShould get error not found : %s.", dbtest.Failed, testID, err)
			}

			t.Logf("\t%s\tTest %d:\tShould be able to update a cafe that does not exist.", dbtest.Success, testID)

			err = core.DeleteCafe(caf.ID)

			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to delete a cafe : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to delete a cafe.", dbtest.Success, testID)

		}
	}

}
