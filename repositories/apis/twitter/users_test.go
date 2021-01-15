package twitter_test

import (
	"testing"

	"github.com/jake-hansen/followrs/repositories/apis/twitter"
	"github.com/stretchr/testify/assert"
)

// TestUserService_Show tests the Show function in UserService.
func TestUserService_Show(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mux, server := NewTestServer()

		defer server.Close()

		testUser := &twitter.User{
			ID:       "1",
			Name:     "testname",
			Username: "testusername",
		}

		wrapper := &twitter.DataWrapper{Data: testUser}

		StandardHandler(t, mux, "/users/by/username/test", wrapper)

		client, _ := twitter.NewTwitterAPI(server.URL, "", "", "")

		user, err := client.UserService.Show("test")
		assert.NoError(t, err)
		assert.Equal(t, testUser, user)
	})

	t.Run("user-not-found", func(t *testing.T) {
		mux, server := NewTestServer()

		defer server.Close()

		wrapper := &twitter.DataWrapper{
			Errors: new([]twitter.Error),
		}

		notFoundError := twitter.Error{
			Title: "Not Found Error",
		}

		newErrorSlice := append(*wrapper.Errors, notFoundError)
		wrapper.Errors = &newErrorSlice

		StandardHandler(t, mux, "/users/by/username/notfound", wrapper)

		client, _ := twitter.NewTwitterAPI(server.URL, "", "", "")

		user, err := client.UserService.Show("notfound")
		assert.Nil(t, user)
		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
	})
}
