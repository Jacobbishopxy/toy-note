package persistence

import (
	"fmt"
	"testing"
	"time"
	"toy-note/api/entity"
	"toy-note/logger"

	"github.com/stretchr/testify/require"
)

// ============================================================================
// WARNING! All the test cases should be run sequentially.
// Otherwise, skip any previous test case can cause the next one to fail.
// ============================================================================

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

// Auto migrate the database and truncate all tables
func TestConnectionAndDataMigrationAndTruncateAll(t *testing.T) {

	r, err := newPgRepo()
	require.NoError(t, err)

	err = r.AutoMigrate()
	require.NoError(t, err)

	err = r.TruncateAll()
	require.NoError(t, err)
}

// Reset test case to initial state
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
// - TestUpdatePost
// - TestDeletePost
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

func TestUpdatePost(t *testing.T) {
	r, err := newPgRepo()
	require.NoError(t, err)

	// modify the Id 1 post.
	// - content changed
	// - affiliates #3 removed
	// - tag #1 removed
	newPost := entity.Post{
		UnitId:  entity.UnitId{Id: 2},
		Content: "updated content",
		Tags: []entity.Tag{
			{
				UnitId: entity.UnitId{Id: 2},
			},
		},
	}

	newPost, err = r.UpdatePost(newPost)

	fmt.Println(newPost)

	require.NoError(t, err)
	// Content should be updated
	require.Equal(t, newPost.Content, "updated content")
	// remain affiliates unchanged
	require.Len(t, newPost.Affiliates, 0)
	// unbound #1 tag
	require.Len(t, newPost.Tags, 1)
	require.Equal(t, newPost.Tags[0].Id, uint(2))
}

func TestDeletePost(t *testing.T) {
	r, err := newPgRepo()
	require.NoError(t, err)

	err = r.DeletePost(2)
	require.NoError(t, err)
}

// ============================================================================
// Test cases for Affiliates table
// - GetUnownedAffiliates
// - DeleteUnownedAffiliates
// ============================================================================

func TestGetUnownedAffiliates(t *testing.T) {
	r, err := newPgRepo()
	require.NoError(t, err)

	affiliates, err := r.GetUnownedAffiliates(entity.NewPagination(0, 10))
	require.NoError(t, err)

	require.Len(t, affiliates, 1)
	require.Equal(t, affiliates[0].Filename, "test3.txt")
}

func TestDeleteUnownedAffiliate(t *testing.T) {
	r, err := newPgRepo()
	require.NoError(t, err)

	// since #3 is unowned, it should be deleted;
	// but #1 is owned, so it should not be deleted.
	err = r.DeleteUnownedAffiliates([]uint{1, 3})
	require.NoError(t, err)
}
