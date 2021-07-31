package gomap

import (
	"context"
	"errors"
	"testing"

	"github.com/mazrean/separated-webshell/domain"
	"github.com/mazrean/separated-webshell/domain/values"
	"github.com/mazrean/separated-webshell/store"
	"github.com/stretchr/testify/assert"
)

func TestWorkspace(t *testing.T) {
	t.Parallel()

	t.Run("Set", testSet)
	t.Run("Get", testGet)
}

func testSet(t *testing.T) {
	t.Parallel()
	t.Helper()

	w := NewWorkspace()

	testUserName, err := values.NewUserName("testUser")
	if err != nil {
		t.Errorf("Error creating test user name: %s", err)
	}

	testWorkspace := domain.NewWorkspace("test", "testWorkspace", testUserName)

	tests := []struct {
		description string
		userName    values.UserName
		workspace   *domain.Workspace
		isErr       bool
		err         error
	}{
		{
			description: "create workspace",
			userName:    testUserName,
			workspace:   testWorkspace,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctx := context.Background()

			err := w.Set(ctx, test.userName, test.workspace)

			if !test.isErr {
				assert.NoError(t, err)

				iWorkspace, ok := w.syncMap.Load(test.userName)

				assert.True(t, ok)
				assert.Equal(t, test.workspace, iWorkspace.(*domain.Workspace))
			} else {
				assert.Error(t, err)

				if test.err != nil && !errors.Is(err, test.err) {
					t.Errorf("expected error %+v, got %+v", test.err, err)
				}
			}
		})
	}
}

func testGet(t *testing.T) {
	t.Parallel()
	t.Helper()

	w := NewWorkspace()

	testUserName, err := values.NewUserName("testUser")
	if err != nil {
		t.Errorf("Error creating test user name: %s", err)
	}

	testNotSetUserName, err := values.NewUserName("testNotSetUser")
	if err != nil {
		t.Errorf("Error creating test user name: %s", err)
	}

	testWorkspace := domain.NewWorkspace("test", "testWorkspace", testUserName)

	tests := []struct {
		description string
		isSet       bool
		userName    values.UserName
		isErr       bool
		err         error
		workspace   *domain.Workspace
	}{
		{
			description: "workspace exists",
			isSet:       true,
			userName:    testUserName,
			workspace:   testWorkspace,
		},
		{
			description: "workspace does not exist",
			userName:    testNotSetUserName,
			isErr:       true,
			err:         store.ErrWorkspaceNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctx := context.Background()

			if test.isSet {
				w.syncMap.Store(test.userName, test.workspace)
			}

			workspace, err := w.Get(ctx, test.userName)

			if !test.isErr {
				assert.NoError(t, err)

				assert.Equal(t, test.workspace, workspace)
			} else {
				assert.Error(t, err)

				if test.err != nil && !errors.Is(err, test.err) {
					t.Errorf("expected error %+v, got %+v", test.err, err)
				}
			}
		})
	}
}
