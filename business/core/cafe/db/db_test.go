package db

import (
	"context"
	"fmt"
	"github.com/colmmurphy91/go-service/business/data/dbtest"
	"github.com/colmmurphy91/go-service/foundation/docker"
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

	store := NewStore(log, db, context.Background())

	t.Log("Given the need to work with Cafe records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single Cafe.", testID)
		{

			c := Cafe{
				Name:        "Spinneys Kitchen",
				Address:     "Meydan",
				PhoneNumber: "0581234567",
			}

			cafe, err := store.Save(c)

			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create a cafe : %s.", dbtest.Failed, testID, err)
			}

			if cafe.ID.Hex() == "" {
				t.Fatalf("\t%s\tTest %d:\tShould have stored ID: %s.", dbtest.Failed, testID, err)
			}

			t.Logf("\t%s\tTest %d:\tShould be able to create a cafe.", dbtest.Success, testID)

			cafes, err := store.FindAll()

			if len(cafes) == 0 {
				t.Fatalf("\t%s\tTest %d:\tShould not have found cafe, should be deleted: %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to all cafes.", dbtest.Success, testID)

			savedID := cafe.ID.Hex()

			foundCafe, err := store.FindById(savedID)

			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to find a cafe by ID: %s.", dbtest.Failed, testID, err)
			}

			t.Logf("\t%s\tTest %d:\tShould be able to find a cafe.", dbtest.Success, testID)

			foundCafe.Name = "New Name"

			err = store.UpdateCafe(foundCafe)

			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to update a cafe : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to update a cafe.", dbtest.Success, testID)

			updatedCafe, err := store.FindById(savedID)

			if updatedCafe.Name != "New Name" {
				t.Fatalf("\t%s\tTest %d:\tShould be able to update a cafe should have new name: %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to update a cafe.", dbtest.Success, testID)

			err = store.DeleteByID(savedID)

			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to delete a cafe : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to delete a cafe.", dbtest.Success, testID)

			foundCafe, err = store.FindById(savedID)

			if err == nil {
				t.Fatalf("\t%s\tTest %d:\tShould not have found cafe, should be deleted: %s.", dbtest.Failed, testID, err)
			}

		}
	}

}
