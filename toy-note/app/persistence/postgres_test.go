package persistence

import (
	"fmt"
	"testing"
	"time"
	"toy-note/app/entity"
	"toy-note/logger"

	"github.com/stretchr/testify/require"
)

const logPath = "test.log"

var sqlConn = PgConn{
	host:     "localhost",
	port:     5432,
	user:     "root",
	password: "secret",
	dbname:   "dev",
	sslmode:  "disable",
}

func newPgRepo() (PgRepository, error) {
	if err := logger.Init("debug", logPath, true); err != nil {
		panic(err)
	}

	return NewPgRepository(logger.TNLogger, sqlConn)
}

func TestConnectionAndDataMigrationAndTruncateAll(t *testing.T) {

	r, err := newPgRepo()
	require.NoError(t, err)

	err = r.AutoMigrate()
	require.NoError(t, err)

	err = r.TruncateAll()
	require.NoError(t, err)
}

func TestTruncateAll(t *testing.T) {
	r, err := newPgRepo()
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

// ============================================================================
// Test cases for Posts table
// - TestCreatePost
// - TestGetAllPosts
// ============================================================================

func TestCreatePost(t *testing.T) {
	r, err := newPgRepo()
	require.NoError(t, err)

	post1 := entity.Post{
		Title:    "My first post",
		Subtitle: "1st subtitle",
		Content:  "1st post for test case",
		Date:     time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
		Affiliates: []entity.Affiliate{
			{
				Filename: "test1.txt",
			},
			{
				Filename: "test2.txt",
			},
		},
		Tags: []entity.Tag{
			{
				UnitId: entity.UnitId{Id: 1},
			},
		},
	}
	post2 := entity.Post{
		Title:    "The second post",
		Subtitle: "blah blah",
		Content:  "1st post for test case",
		Date:     time.Date(2021, time.April, 1, 0, 0, 0, 0, time.UTC),
		Affiliates: []entity.Affiliate{
			{
				Filename: "test3.txt",
			},
		},
		Tags: []entity.Tag{
			{
				UnitId: entity.UnitId{Id: 1},
			},
			{
				UnitId: entity.UnitId{Id: 2},
			},
		},
	}

	post1, err = r.CreatePost(post1)
	require.NoError(t, err)

	post2, err = r.CreatePost(post2)
	require.NoError(t, err)

}

func TestGetAllPosts(t *testing.T) {
	r, err := newPgRepo()
	require.NoError(t, err)

	posts, err := r.GetPosts(entity.NewPagination(0, 10))
	require.NoError(t, err)
	// since we have already created two posts, we should get two posts
	require.Len(t, posts, 2)
	// the 1st post should have two affiliates and one tag
	require.Equal(t, posts[0].Affiliates[0].Id, uint(1))
	require.Equal(t, posts[0].Affiliates[1].Id, uint(2))
	require.Equal(t, posts[0].Tags[0].Id, uint(1))
	require.Equal(t, posts[1].Affiliates[0].Id, uint(3))
	require.Equal(t, posts[1].Tags[0].Id, uint(1))
	require.Equal(t, posts[1].Tags[1].Id, uint(2))
}
