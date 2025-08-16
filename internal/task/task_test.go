package task

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTaskStatusConstants(t *testing.T) {
	assert.Equal(t, Status("pending"), StatusPending, "StatusPending should be 'pending'")
	assert.Equal(t, Status("in_progress"), StatusInProgress, "StatusInProgress should be 'in_progress'")
	assert.Equal(t, Status("completed"), StatusCompleted, "StatusCompleted should be 'completed'")
}

func TestTaskCreation(t *testing.T) {
	task := Task{
		ID:     1,
		Title:  "Test Task",
		Status: StatusPending,
	}

	assert.Equal(t, 1, task.ID, "Expected ID 1")
	assert.Equal(t, "Test Task", task.Title, "Expected title 'Test Task'")
	assert.Equal(t, StatusPending, task.Status, "Expected status 'pending'")
}
