package service

import (
	"os"
	"testing"
	"toy-note/api/entity"
	"toy-note/api/persistence"
	"toy-note/logger"

	"github.com/stretchr/testify/require"
)

const logPath = "test.log"

var sqlConn = persistence.PgConn{
	Host:    "localhost",
	Port:    5432,
	User:    "root",
	Pass:    "secret",
	Db:      "dev",
	Sslmode: "disable",
}

var mongoConn = persistence.MongoConn{
	Host: "localhost",
	Port: 27017,
	User: "root",
	Pass: "secret",
}

func newService() (*ToyNoteService, error) {
	if err := logger.Init("debug", logPath, true); err != nil {
		panic(err)
	}

	return NewToyNoteService(logger.TNLogger, sqlConn, mongoConn)
}

func TestNewService(t *testing.T) {

	_, err := newService()

	require.NoError(t, err)
}

/*
In this test case, we don't need to test methods that already tested in persistence package.
Hence, we only need to care about the compositional methods:

- DownloadAffiliate
- RebindAffiliate
- DeleteUnownedAffiliates
*/

func TestDownloadAffiliate(t *testing.T) {
	s, err := newService()
	require.NoError(t, err)

	filenameUsedForSaving := "test.log.bak"

	// Step 1 & 2 are the prerequisites for the following test cases

	// 1. Open a file
	reader, err := os.Open("test.log")
	require.NoError(t, err)

	// 2. Upload the file to mongo, and get the file ObjectId
	id, err := s.UploadAffiliate(reader, filenameUsedForSaving)
	require.NoError(t, err)

	// 3. Create a post with affiliate (with the file ObjectId)
	post := entity.Post{
		Title:   "test",
		Content: "test note service",
		Affiliates: []entity.Affiliate{
			{
				ObjectId: id,
				Filename: filenameUsedForSaving,
			},
		},
	}

	post, err = s.SavePost(post)
	require.NoError(t, err)
	require.NotEmpty(t, post.Id)
	require.NotEmpty(t, post.Affiliates[0].Id)

	fo, err := s.DownloadAffiliate(post.Affiliates[0].Id)
	require.NoError(t, err)
	require.Equal(t, filenameUsedForSaving, fo.Filename)
	require.NotEmpty(t, fo.Size)
	require.NotEmpty(t, fo.Content)
}

func TestRebindAndDeleteUnownedAffiliate(t *testing.T) {
	s, err := newService()
	require.NoError(t, err)

	filename1 := "test.txt"
	filename2 := "dev.txt"

	// create a post with an affiliate
	post := entity.Post{
		Title:   "test",
		Content: "test note service",
		Affiliates: []entity.Affiliate{
			{
				ObjectId: "000000000000000000000000",
				Filename: filename1,
			},
			{
				ObjectId: "999999999999999999999999",
				Filename: filename2,
			},
		},
	}

	// save the post
	post, err = s.SavePost(post)
	require.NoError(t, err)

	// affiliates has been created and its id is given by Pg
	affiliateId1 := post.Affiliates[0].Id
	require.NotEmpty(t, affiliateId1)
	affiliateId2 := post.Affiliates[2].Id
	require.NotEmpty(t, affiliateId2)

	// remove all affiliates from post and save the post again
	post.Affiliates = []entity.Affiliate{}
	post, err = s.SavePost(post)
	require.NoError(t, err)
	require.Empty(t, post.Affiliates)

	// check if the affiliates are still there, and affiliate1's post_refer is now empty
	check1, err := s.pg.GetAffiliate(affiliateId1)
	require.NoError(t, err)
	require.Empty(t, check1.PostRefer)
	// affiliate2 is still there, and it's post_refer is empty as well
	check2, err := s.pg.GetAffiliate(affiliateId1)
	require.NoError(t, err)
	require.Empty(t, check2.PostRefer)

	// now rebind affiliate1 to the post
	err = s.RebindAffiliate(post.Id, affiliateId1)
	require.NoError(t, err)

	// check again the affiliate has been rebound to the post
	check1, err = s.pg.GetAffiliate(affiliateId1)
	require.NoError(t, err)
	require.Equal(t, check1.PostRefer, post.Id)

	// now remove these affiliates
	err = s.DeleteUnownedAffiliates([]uint{affiliateId1, affiliateId2})
	require.NoError(t, err)

	// make sure only affiliate2 is deleted, since it's not referred by any post
	check1, err = s.pg.GetAffiliate(affiliateId1)
	require.NoError(t, err)
	require.Empty(t, check1.PostRefer)
	check2, err = s.pg.GetAffiliate(affiliateId2)
	require.NoError(t, err)
}
