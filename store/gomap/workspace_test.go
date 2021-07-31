package gomap

import (
	"context"
	"errors"
	"testing"

	"github.com/mazrean/separated-webshell/domain"
	"github.com/mazrean/separated-webshell/domain/values"
	"github.com/stretchr/testify/assert"
)

func TestWorkspace(t *testing.T) {
	t.Parallel()

	t.Run("Set", testSet)
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
