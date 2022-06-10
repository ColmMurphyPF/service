package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/colmmurphy91/go-service/app/services/sales-api/handlers"
	"github.com/colmmurphy91/go-service/business/core/cafe"
	"github.com/colmmurphy91/go-service/business/data/dbtest"
	"github.com/colmmurphy91/go-service/business/sys/validate"
	v1Web "github.com/colmmurphy91/go-service/business/web/v1"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type CafeTests struct {
	app       http.Handler
	userToken string
}

func TestCafe(t *testing.T) {
	t.Parallel()

	test := dbtest.NewMongoIntegration(t, mongoC, "intercafes")
	t.Cleanup(test.Teardown)

	shutdown := make(chan os.Signal, 1)
	tests := CafeTests{
		app: handlers.APIMux(handlers.APIMuxConfig{
			Shutdown: shutdown,
			Log:      test.Log,
			Auth:     test.Auth,
			MDB:      test.DB,
		}),
		userToken: test.AdminToken(),
	}

	t.Run("postCafe400", tests.postCafe400)
	t.Run("postCafe401", tests.postCafe401)
	t.Run("postCafe404", tests.getCafe404)
	t.Run("deleteCafe404", tests.deleteCafeNotFound)
	t.Run("UpdateCafe404", tests.putCafe404)
	t.Run("CrudTests", tests.crudProduct)
}

func (ct *CafeTests) postCafe400(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/v1/cafes", strings.NewReader(`{}`))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+ct.userToken)
	ct.app.ServeHTTP(w, r)

	t.Log("Given the need to validate a new cafe can't be created with an invalid document.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using an incomplete cafe value.", testID)
		{
			if w.Code != http.StatusBadRequest {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 400 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 400 for the response.", dbtest.Success, testID)

			// Inspect the response.
			var got v1Web.ErrorResponse
			if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to unmarshal the response to an error type : %v", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to unmarshal the response to an error type.", dbtest.Success, testID)

			fields := validate.FieldErrors{
				{Field: "name", Error: "name is a required field"},
				{Field: "address", Error: "address is a required field"},
				{Field: "phone_number", Error: "phone_number is a required field"},
			}
			exp := v1Web.ErrorResponse{
				Error:  "data validation error",
				Fields: fields.Fields(),
			}

			// We can't rely on the order of the field errors so they have to be
			// sorted. Tell the cmp package how to sort them.
			sorter := cmpopts.SortSlices(func(a, b validate.FieldError) bool {
				return a.Field < b.Field
			})

			if diff := cmp.Diff(got, exp, sorter); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get the expected result. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get the expected result.", dbtest.Success, testID)
		}
	}
}

func (ct *CafeTests) postCafe401(t *testing.T) {
	nc := cafe.NewCafe{
		Name:        "Comic Books",
		Address:     "Testing",
		PhoneNumber: "Testing",
	}

	body, err := json.Marshal(&nc)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodPost, "/v1/cafes", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// Not setting an authorization header.
	ct.app.ServeHTTP(w, r)

	t.Log("Given the need to validate a new cafe can't be created with an invalid document.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using an incomplete cafe value.", testID)
		{
			if w.Code != http.StatusUnauthorized {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 401 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 401 for the response.", dbtest.Success, testID)
		}
	}
}

func (ct *CafeTests) getCafe404(t *testing.T) {
	id := "a224a8d6-3f9e-4b11-9900-e81a25d80702"

	r := httptest.NewRequest(http.MethodGet, "/v1/cafe/"+id, nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+ct.userToken)
	ct.app.ServeHTTP(w, r)

	t.Log("Given the need to validate getting a cafe with an unknown id.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the new cafe %s.", testID, id)
		{
			if w.Code != http.StatusNotFound {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 404 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 404 for the response.", dbtest.Success, testID)

			got := w.Body.String()
			fmt.Println(got)
			exp := "not found"
			if !strings.Contains(got, exp) {
				t.Logf("\t\tTest %d:\tGot : %v", testID, got)
				t.Logf("\t\tTest %d:\tExp: %v", testID, exp)
				t.Fatalf("\t%s\tTest %d:\tShould get the expected result.", dbtest.Failed, testID)
			}
			t.Logf("\t%s\tTest %d:\tShould get the expected result.", dbtest.Success, testID)
		}
	}
}

// deleteCafeNotFound validates deleting a cafe that does not exist is not a failure.
func (ct *CafeTests) deleteCafeNotFound(t *testing.T) {
	id := "112262f1-1a77-4374-9f22-39e575aa6348"

	r := httptest.NewRequest(http.MethodDelete, "/v1/cafes/"+id, nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+ct.userToken)
	ct.app.ServeHTTP(w, r)

	t.Log("Given the need to validate deleting a cafe that does not exist.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the new cafe %s.", testID, id)
		{
			if w.Code != http.StatusNoContent {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 204 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 204 for the response.", dbtest.Success, testID)
		}
	}
}

// putCafe404 validates updating a cafe that does not exist.
func (ct *CafeTests) putCafe404(t *testing.T) {
	id := primitive.NewObjectID().Hex()
	fmt.Println(id)

	up := cafe.UpdateCafe{
		Name: dbtest.StringPointer("Nonexistent"),
	}
	body, err := json.Marshal(&up)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodPut, "/v1/cafes/"+id, bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+ct.userToken)
	ct.app.ServeHTTP(w, r)

	t.Log("Given the need to validate updating a cafe that does not exist.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the new cafe %s.", testID, id)
		{
			if w.Code != http.StatusNotFound {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 404 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 404 for the response.", dbtest.Success, testID)

			got := w.Body.String()
			exp := "not found"
			if !strings.Contains(got, exp) {
				t.Logf("\t\tTest %d:\tGot : %v", testID, got)
				t.Logf("\t\tTest %d:\tExp: %v", testID, exp)
				t.Fatalf("\t%s\tTest %d:\tShould get the expected result.", dbtest.Failed, testID)
			}
			t.Logf("\t%s\tTest %d:\tShould get the expected result.", dbtest.Success, testID)
		}
	}
}

func (ct *CafeTests) crudProduct(t *testing.T) {
	p := ct.postCafe201(t)
	defer ct.deleteCafe204(t, p.ID)

	ct.getCafe200(t, p.ID)
	ct.getAllCafes(t)
	ct.putCafe204(t, p.ID)
}

func (ct *CafeTests) postCafe201(t *testing.T) cafe.Cafe {
	nc := cafe.NewCafe{
		Name:        "Comic Books",
		Address:     "Testing",
		PhoneNumber: "Testing",
	}

	body, err := json.Marshal(&nc)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodPost, "/v1/cafes", bytes.NewBuffer(body))
	r.Header.Set("Authorization", "Bearer "+ct.userToken)
	w := httptest.NewRecorder()

	// Not setting an authorization header.
	ct.app.ServeHTTP(w, r)

	var got cafe.Cafe

	t.Log("Given the need to validate a new cafe can't be created with an invalid document.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using an incomplete cafe value.", testID)
		{
			if w.Code != http.StatusCreated {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 201 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 201 for the response.", dbtest.Success, testID)

			if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to unmarshal the response : %v", dbtest.Failed, testID, err)
			}

			// Define what we wanted to receive. We will just trust the generated
			// fields like ID and Dates so we copy p.
			exp := got
			exp.Name = "Comic Books"
			exp.Address = "Testing"
			exp.PhoneNumber = "Testing"

			if diff := cmp.Diff(got, exp); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get the expected result. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get the expected result.", dbtest.Success, testID)
		}
	}
	return got
}

func (ct *CafeTests) getCafe200(t *testing.T, id string) {
	r := httptest.NewRequest(http.MethodGet, "/v1/cafes/"+id, nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+ct.userToken)
	ct.app.ServeHTTP(w, r)

	t.Log("Given the need to validate getting a cafe that exists.")
	{
		testID := 0
		t.Logf("\tTest : %d\tWhen using the new cafe %s.", testID, id)
		{
			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\tTest : %d\tShould receive a status code of 200 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest : %d\tShould receive a status code of 200 for the response.", dbtest.Success, testID)

			var got cafe.Cafe
			if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
				t.Fatalf("\t%s\tTest : %d\tShould be able to unmarshal the response : %v", dbtest.Failed, testID, err)
			}

			// Define what we wanted to receive. We will just trust the generated
			// fields like Dates so we copy p.
			exp := got
			exp.ID = id
			exp.Name = "Comic Books"
			exp.Address = "Testing"
			exp.PhoneNumber = "Testing"

			if diff := cmp.Diff(got, exp); diff != "" {
				t.Fatalf("\t%s\tTest : %d\tShould get the expected result. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest : %d\tShould get the expected result.", dbtest.Success, testID)
		}
	}
}

func (ct *CafeTests) putCafe204(t *testing.T, id string) {
	body := `{"name": "Graphic Novels", "address": "100, and 2"}`
	r := httptest.NewRequest(http.MethodPut, "/v1/cafes/"+id, strings.NewReader(body))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+ct.userToken)
	ct.app.ServeHTTP(w, r)

	t.Log("Given the need to update a cafe with the products endpoint.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the modified cafe value.", testID)
		{
			if w.Code != http.StatusNoContent {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 204 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 204 for the response.", dbtest.Success, testID)

			r = httptest.NewRequest(http.MethodGet, "/v1/cafes/"+id, nil)
			w = httptest.NewRecorder()

			r.Header.Set("Authorization", "Bearer "+ct.userToken)
			ct.app.ServeHTTP(w, r)

			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 200 for the retrieve : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 200 for the retrieve.", dbtest.Success, testID)

			var ru cafe.Cafe
			if err := json.NewDecoder(w.Body).Decode(&ru); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to unmarshal the response : %v", dbtest.Failed, testID, err)
			}

			if ru.Name != "Graphic Novels" {
				t.Fatalf("\t%s\tTest %d:\tShould see an updated Name : got %q want %q", dbtest.Failed, testID, ru.Name, "Graphic Novels")
			}
			t.Logf("\t%s\tTest %d:\tShould see an updated Name.", dbtest.Success, testID)
		}
	}
}

func (ct *CafeTests) deleteCafe204(t *testing.T, id string) {
	r := httptest.NewRequest(http.MethodDelete, "/v1/cafes/"+id, nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+ct.userToken)
	ct.app.ServeHTTP(w, r)

	t.Log("Given the need to validate deleting a cafe that does exist.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the new cafe %s.", testID, id)
		{
			if w.Code != http.StatusNoContent {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 204 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 204 for the response.", dbtest.Success, testID)
		}
	}
}

func (ct *CafeTests) getAllCafes(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/v1/cafes", nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+ct.userToken)
	ct.app.ServeHTTP(w, r)

	t.Log("Given the need to validate getting a list cafe that exists.")
	{
		testID := 0
		t.Logf("\tTest : %d\tWhen using get all cafes.", testID)
		{
			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\tTest : %d\tShould receive a status code of 200 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest : %d\tShould receive a status code of 200 for the response.", dbtest.Success, testID)

			var got []cafe.Cafe
			if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
				t.Fatalf("\t%s\tTest : %d\tShould be able to unmarshal the response : %v", dbtest.Failed, testID, err)
			}

			// Define what we wanted to receive. We will just trust the generated
			// fields like Dates so we copy p.
			exp := got
			exp[0].Name = "Comic Books"
			exp[0].Address = "Testing"
			exp[0].PhoneNumber = "Testing"

			if diff := cmp.Diff(got, exp); diff != "" {
				t.Fatalf("\t%s\tTest : %d\tShould get the expected result. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest : %d\tShould get the expected result.", dbtest.Success, testID)
		}
	}
}
