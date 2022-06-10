package db

//
//import (
//	"context"
//	"database/sql"
//	"fmt"
//	"github.com/colmmurphy91/go-service/business/data/dbtest"
//	"github.com/colmmurphy91/go-service/foundation/docker"
//	"github.com/google/go-cmp/cmp"
//	"github.com/google/uuid"
//	"testing"
//	"time"
//)
//
//var c *docker.Container
//
//func TestMain(m *testing.M) {
//	var err error
//	c, err = dbtest.StartDB()
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer dbtest.StopDB(c)
//
//	m.Run()
//}
//
//func TestCafe(t *testing.T) {
//	log, db, teardown := dbtest.NewUnit(t, c, "testuser")
//	t.Cleanup(teardown)
//
//	store := NewStore(log, db)
//
//	ctx := context.Background()
//
//	t.Log("Given the need to work with Product records.")
//	{
//		testID := 0
//		t.Logf("\tTest %d:\tWhen handling a single Product.", testID)
//		{
//			nu := User{
//				ID:           uuid.NewString(),
//				Name:         "Bill Kennedy",
//				Email:        "bill@ardanlabs.com",
//				Roles:        []string{"admin"},
//				PasswordHash: []byte("gophers"),
//				DateCreated:  time.Now().UTC(),
//				DateUpdated:  time.Now().UTC(),
//				Confirmed:    false,
//				ConfirmHash: sql.NullInt64{
//					Valid: true,
//					Int64: 123,
//				},
//			}
//
//			err := store.Create(ctx, nu)
//			if err != nil {
//				t.Fatalf("\t%s\tTest %d:\tShould be able to create a user : %s.", dbtest.Failed, testID, err)
//			}
//
//			dbUser, err := store.QueryByEmail(ctx, nu.Email)
//
//			if err != nil {
//				t.Fatalf("\t%s\tTest %d:\tShould be able to find a user : %s.", dbtest.Failed, testID, err)
//			}
//
//			if diff := cmp.Diff(nu, dbUser); diff != "" {
//				t.Fatalf("\t%s\tTest %d:\tShould get back the same user. Diff:\n%s", dbtest.Failed, testID, diff)
//			}
//
//			if dbUser.Confirmed {
//				t.Fatalf("\t%s\tTest %d:\tShould not be confirmed : %s.", dbtest.Failed, testID, err)
//			}
//
//			fmt.Println(dbUser.ConfirmHash)
//
//			dbUser.Confirmed = true
//			dbUser.ConfirmHash = sql.NullInt64{
//				Valid: false,
//			}
//
//			err = store.Update(ctx, dbUser)
//			if err != nil {
//				t.Fatalf("\t%s\tTest %d:\tShould be able to find a user : %s.", dbtest.Failed, testID, err)
//			}
//
//			sUser, err := store.QueryByEmail(ctx, nu.Email)
//
//			fmt.Println(sUser)
//
//		}
//	}
//
//}
