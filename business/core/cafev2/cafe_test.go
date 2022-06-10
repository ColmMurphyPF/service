package cafev2

import (
	"context"
	"fmt"
	"github.com/colmmurphy91/go-service/business/data/dbtest"
	"github.com/colmmurphy91/go-service/foundation/docker"
	"github.com/google/go-cmp/cmp"
	"testing"
)

var c *docker.Container

func TestMain(m *testing.M) {
	var err error
	c, err = dbtest.StartDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dbtest.StopDB(c)

	m.Run()
}

func TestCafeV2(t *testing.T) {
	log, db, teardown := dbtest.NewUnit(t, c, "testprod")
	t.Cleanup(teardown)

	core := NewCore(log, db)

	t.Log("Given the need to work with Product records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single Product.", testID)
		{
			ctx := context.Background()

			nc := NewCafe{
				Name:    "Comic Books",
				Address: "My Address",
				LogoURL: "www.blah.com",
			}

			cafe, err := core.Create(ctx, nc, "45b5fbd3-755f-4379-8f07-a58d4a30fa2e")
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create a cafe : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create a cafe.", dbtest.Success, testID)

			saved, err := core.QueryByOwnerID(ctx, "45b5fbd3-755f-4379-8f07-a58d4a30fa2e")
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve product by ID: %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve product by ID.", dbtest.Success, testID)

			if diff := cmp.Diff(cafe, saved); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same product. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same product.", dbtest.Success, testID)

			_, err = core.Create(ctx, nc, "45b5fbd3-755f-4379-8f07-a58d4a30fa2e")
			if err == nil {
				t.Fatalf("\t%s\tTest %d:\tShould be not be able to create 2 cafes : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould NOT be able to create a cafe.", dbtest.Success, testID)

			//upd := product.UpdateProduct{
			//	Name:     dbtest.StringPointer("Comics"),
			//	Cost:     dbtest.IntPointer(50),
			//	Quantity: dbtest.IntPointer(40),
			//}
			//updatedTime := time.Date(2019, time.January, 1, 1, 1, 1, 0, time.UTC)
			//
			//if err := core.Update(ctx, prd.ID, upd, updatedTime); err != nil {
			//	t.Fatalf("\t%s\tTest %d:\tShould be able to update product : %s.", dbtest.Failed, testID, err)
			//}
			//t.Logf("\t%s\tTest %d:\tShould be able to update product.", dbtest.Success, testID)

			//products, err := core.Query(ctx, 1, 3)
			//if err != nil {
			//	t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve updated product : %s.", dbtest.Failed, testID, err)
			//}
			//t.Logf("\t%s\tTest %d:\tShould be able to retrieve updated product.", dbtest.Success, testID)

			// Check specified fields were updated. Make a copy of the original product
			// and change just the fields we expect then diff it with what was saved.
			//want := prd
			//want.Name = *upd.Name
			//want.Cost = *upd.Cost
			//want.Quantity = *upd.Quantity
			//want.DateUpdated = updatedTime
			//
			//var idx int
			//for i, p := range products {
			//	if p.ID == want.ID {
			//		idx = i
			//	}
			//}
			//if diff := cmp.Diff(want, products[idx]); diff != "" {
			//	t.Fatalf("\t%s\tTest %d:\tShould get back the same product. Diff:\n%s", dbtest.Failed, testID, diff)
			//}
			//t.Logf("\t%s\tTest %d:\tShould get back the same product.", dbtest.Success, testID)
			//
			//upd = product.UpdateProduct{
			//	Name: dbtest.StringPointer("Graphic Novels"),
			//}
			//
			//if err := core.Update(ctx, prd.ID, upd, updatedTime); err != nil {
			//	t.Fatalf("\t%s\tTest %d:\tShould be able to update just some fields of product : %s.", dbtest.Failed, testID, err)
			//}
			//t.Logf("\t%s\tTest %d:\tShould be able to update just some fields of product.", dbtest.Success, testID)
			//
			//saved, err = core.QueryByID(ctx, prd.ID)
			//if err != nil {
			//	t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve updated product : %s.", dbtest.Failed, testID, err)
			//}
			//t.Logf("\t%s\tTest %d:\tShould be able to retrieve updated product.", dbtest.Success, testID)
			//
			//if saved.Name != *upd.Name {
			//	t.Fatalf("\t%s\tTest %d:\tShould be able to see updated Name field : got %q want %q.", dbtest.Failed, testID, saved.Name, *upd.Name)
			//} else {
			//	t.Logf("\t%s\tTest %d:\tShould be able to see updated Name field.", dbtest.Success, testID)
			//}
			//
			//if err := core.Delete(ctx, prd.ID); err != nil {
			//	t.Fatalf("\t%s\tTest %d:\tShould be able to delete product : %s.", dbtest.Failed, testID, err)
			//}
			//t.Logf("\t%s\tTest %d:\tShould be able to delete product.", dbtest.Success, testID)
			//
			//_, err = core.QueryByID(ctx, prd.ID)
			//if !errors.Is(err, product.ErrNotFound) {
			//	t.Fatalf("\t%s\tTest %d:\tShould NOT be able to retrieve deleted product : %s.", dbtest.Failed, testID, err)
			//}
			//t.Logf("\t%s\tTest %d:\tShould NOT be able to retrieve deleted product.", dbtest.Success, testID)
		}
	}
}
