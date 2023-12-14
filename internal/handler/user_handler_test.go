package handler

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestUserHandler_Get(t *testing.T) {
	req := &http.Request{
		Method: http.MethodGet,
	}
	rw := httptest.NewRecorder()

	handler := NewUserHandler(log.Default())
	handler.users = []*User{
		{
			Id:    uuid.MustParse("0d9496f4-acd4-4e7c-a92d-924edae4c908"),
			Email: "tester1@example.com",
		},
		{
			Id:    uuid.MustParse("755e5e2c-bbb2-4edc-91a8-959fd190f465"),
			Email: "tester2@example.com",
		},
	}

	handler.ServeHTTP(rw, req)

	expectedBody := `[{"id":"0d9496f4-acd4-4e7c-a92d-924edae4c908","email":"tester1@example.com"},{"id":"755e5e2c-bbb2-4edc-91a8-959fd190f465","email":"tester2@example.com"}]`
	expectedBody += "\n"

	if rw.Code != http.StatusOK {
		t.Errorf("invalid status code: got: %d, want: %d", rw.Code, http.StatusOK)
	}

	if rw.Body.String() != expectedBody {
		t.Errorf("invalid response body: got: %s, want: %s", rw.Body, expectedBody)
	}
}

func TestUserHandler_Post(t *testing.T) {
	req := &http.Request{
		Method: http.MethodPost,
		Body:   io.NopCloser(strings.NewReader(`{"id":"c3f17f85-7187-4d80-96d4-7d50614bf14a","email":"tester@example.com"}`)),
	}
	rw := httptest.NewRecorder()

	handler := NewUserHandler(log.Default())
	handler.ServeHTTP(rw, req)

	if rw.Code != http.StatusCreated {
		t.Errorf("invalid status code: got: %d, want: %d", rw.Code, http.StatusCreated)
	}

	if rw.Body.String() != "" {
		t.Errorf("invalid response body: got: %s, want: %s", rw.Body, "")
	}

	expectedUsers := []*User{
		{
			Id:    uuid.MustParse("c3f17f85-7187-4d80-96d4-7d50614bf14a"),
			Email: "tester@example.com",
		},
	}

	if diff := cmp.Diff(expectedUsers, handler.users); diff != "" {
		t.Errorf("users mismatch (-want +got):\n%s", diff)
	}
}

func TestUserHandler_Patch(t *testing.T) {
	t.Run("update existing user", func(t *testing.T) {
		req := &http.Request{
			Method: http.MethodPatch,
			URL:    &url.URL{RawQuery: `id="755e5e2c-bbb2-4edc-91a8-959fd190f465"`},
			Body:   io.NopCloser(strings.NewReader(`{"email":"updated_tester@example.com"}`)),
		}
		rw := httptest.NewRecorder()

		handler := NewUserHandler(log.Default())
		handler.users = []*User{
			{
				Id:    uuid.MustParse("0d9496f4-acd4-4e7c-a92d-924edae4c908"),
				Email: "tester1@example.com",
			},
			{
				Id:    uuid.MustParse("755e5e2c-bbb2-4edc-91a8-959fd190f465"),
				Email: "tester2@example.com",
			},
		}

		handler.ServeHTTP(rw, req)

		if rw.Code != http.StatusOK {
			t.Errorf("invalid status code: got: %d, want: %d", rw.Code, http.StatusCreated)
		}

		if rw.Body.String() != "" {
			t.Errorf("invalid response body: got: %s, want: %s", rw.Body, "")
		}

		expectedUsers := []*User{
			{
				Id:    uuid.MustParse("0d9496f4-acd4-4e7c-a92d-924edae4c908"),
				Email: "tester1@example.com",
			},
			{
				Id:    uuid.MustParse("755e5e2c-bbb2-4edc-91a8-959fd190f465"),
				Email: "updated_tester@example.com",
			},
		}

		if diff := cmp.Diff(expectedUsers, handler.users); diff != "" {
			t.Errorf("users mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("update not-existing user", func(t *testing.T) {
		req := &http.Request{
			Method: http.MethodPatch,
			URL:    &url.URL{RawQuery: `id="ff19ca4b-3f98-4ea0-adca-852a8f201121"`},
			Body:   io.NopCloser(strings.NewReader(`{"email":"updated_tester@example.com"}`)),
		}
		rw := httptest.NewRecorder()

		handler := NewUserHandler(log.Default())
		handler.users = []*User{
			{
				Id:    uuid.MustParse("0d9496f4-acd4-4e7c-a92d-924edae4c908"),
				Email: "tester1@example.com",
			},
			{
				Id:    uuid.MustParse("755e5e2c-bbb2-4edc-91a8-959fd190f465"),
				Email: "tester2@example.com",
			},
		}

		handler.ServeHTTP(rw, req)

		expectedBody := `{"error": "resource not found"}`

		if rw.Code != http.StatusNotFound {
			t.Errorf("invalid status code: got: %d, want: %d", rw.Code, http.StatusCreated)
		}

		if rw.Body.String() != expectedBody {
			t.Errorf("invalid response body: got: %s, want: %s", rw.Body, expectedBody)
		}

		expectedUsers := []*User{
			{
				Id:    uuid.MustParse("0d9496f4-acd4-4e7c-a92d-924edae4c908"),
				Email: "tester1@example.com",
			},
			{
				Id:    uuid.MustParse("755e5e2c-bbb2-4edc-91a8-959fd190f465"),
				Email: "tester2@example.com",
			},
		}

		if diff := cmp.Diff(expectedUsers, handler.users); diff != "" {
			t.Errorf("users mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestUserHandler_Delete(t *testing.T) {
	t.Run("delete existing user", func(t *testing.T) {
		req := &http.Request{
			Method: http.MethodDelete,
			URL:    &url.URL{RawQuery: `id="755e5e2c-bbb2-4edc-91a8-959fd190f465"`},
		}
		rw := httptest.NewRecorder()

		handler := NewUserHandler(log.Default())
		handler.users = []*User{
			{
				Id:    uuid.MustParse("0d9496f4-acd4-4e7c-a92d-924edae4c908"),
				Email: "tester1@example.com",
			},
			{
				Id:    uuid.MustParse("755e5e2c-bbb2-4edc-91a8-959fd190f465"),
				Email: "tester2@example.com",
			},
		}

		handler.ServeHTTP(rw, req)

		if rw.Code != http.StatusNoContent {
			t.Errorf("invalid status code: got: %d, want: %d", rw.Code, http.StatusNoContent)
		}

		if rw.Body.String() != "" {
			t.Errorf("invalid response body: got: %s, want: %s", rw.Body, "")
		}

		expectedUsers := []*User{
			{
				Id:    uuid.MustParse("0d9496f4-acd4-4e7c-a92d-924edae4c908"),
				Email: "tester1@example.com",
			},
		}

		if diff := cmp.Diff(expectedUsers, handler.users); diff != "" {
			t.Errorf("users mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("delete not-existing user", func(t *testing.T) {
		req := &http.Request{
			Method: http.MethodDelete,
			URL:    &url.URL{RawQuery: `id="ff19ca4b-3f98-4ea0-adca-852a8f201121"`},
		}
		rw := httptest.NewRecorder()

		handler := NewUserHandler(log.Default())
		handler.users = []*User{
			{
				Id:    uuid.MustParse("0d9496f4-acd4-4e7c-a92d-924edae4c908"),
				Email: "tester1@example.com",
			},
			{
				Id:    uuid.MustParse("755e5e2c-bbb2-4edc-91a8-959fd190f465"),
				Email: "tester2@example.com",
			},
		}

		handler.ServeHTTP(rw, req)

		expectedBody := `{"error": "resource not found"}`

		if rw.Code != http.StatusNotFound {
			t.Errorf("invalid status code: got: %d, want: %d", rw.Code, http.StatusCreated)
		}

		if rw.Body.String() != expectedBody {
			t.Errorf("invalid response body: got: %s, want: %s", rw.Body, expectedBody)
		}

		expectedUsers := []*User{
			{
				Id:    uuid.MustParse("0d9496f4-acd4-4e7c-a92d-924edae4c908"),
				Email: "tester1@example.com",
			},
			{
				Id:    uuid.MustParse("755e5e2c-bbb2-4edc-91a8-959fd190f465"),
				Email: "tester2@example.com",
			},
		}

		if diff := cmp.Diff(expectedUsers, handler.users); diff != "" {
			t.Errorf("users mismatch (-want +got):\n%s", diff)
		}
	})
}
