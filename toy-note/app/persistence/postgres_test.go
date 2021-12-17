package persistence

import (
	"fmt"
	"testing"
	"toy-note/app/entity"
	"toy-note/logger"

	"github.com/stretchr/testify/require"
)

const logPath = "test.log"
const sqlUri = `host=localhost
				user=root
				password=secret
				dbname=dev
				port=5432
				sslmode=disable
				TimeZone=Asia/Shanghai`

func newPgRepo() (PgRepository, error) {
	if err := logger.Init("debug", logPath, true); err != nil {
		panic(err)
	}

	return NewPgRepository(logger.TNLogger, sqlUri)
}

func TestConnectionAndDataMigrationAndTruncateAll(t *testing.T) {

	r, err := newPgRepo()
	require.NoError(t, err)

	err = r.AutoMigrate()
	require.NoError(t, err)

	err = r.TruncateAll()
	require.NoError(t, err)
}

// ============================================================================
// Test cases for Tags table
// - TestCreateTagAndGetAll
// - TestUpdateTag
// - TestDeleteTag
// ============================================================================

func TestCreateTagAndGetAll(t *testing.T) {
	r, err := newPgRepo()
	require.NoError(t, err)

	tag1 := entity.Tag{
		Name:        "test",
		Description: "1st test tag",
	}
	tag2 := entity.Tag{
		Name:        "dev",
		Description: "#2",
	}

	tag1, err = r.CreateTag(tag1)
	require.NoError(t, err)
	require.Equal(t, tag1.Id, uint(1))

	tag2, err = r.CreateTag(tag2)
	require.NoError(t, err)
	require.Equal(t, tag2.Id, uint(2))

	tags, err := r.GetTags()
	require.NoError(t, err)
	require.Len(t, tags, 2)

	fmt.Printf("all tags: %v", tags)
}

func TestUpdateTag(t *testing.T) {
	r, err := newPgRepo()
	require.NoError(t, err)

	tag := entity.Tag{
		UnitId:      entity.UnitId{Id: 2},
		Name:        "dev+",
		Description: "#2 ,edited",
	}

	updatedTag, err := r.UpdateTag(tag)
	require.NoError(t, err)

	fmt.Printf("updated tag: %v", updatedTag)
}

func TestDeleteTag(t *testing.T) {
	r, err := newPgRepo()
	require.NoError(t, err)

	err = r.DeleteTag(1)
	require.NoError(t, err)
}
