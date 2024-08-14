package folders_test

import (
	"fmt"
	"testing"

	"github.com/georgechieng-sc/interns-2022/folders"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

type mockDb struct {
	mockData []*folders.Folder
	dbError  error
}

func (m mockDb) FetchAllFoldersByOrgID(orgID uuid.UUID) ([]*folders.Folder, error) {
	if m.dbError != nil {
		return nil, m.dbError
	}
	return m.mockData, nil
}
func newUUID() *uuid.UUID {
	id, _ := uuid.NewV4()
	return &id
}

func Test_GetAllFolders(t *testing.T) {
	t.Run("nil request", func(t *testing.T) {
		assert := assert.New(t)

		resp, err := folders.GetAllFoldersImproved(nil, mockDb{})
		assert.Nil(resp)
		assert.Error(err)
	})

	t.Run("basic request", func(t *testing.T) {
		var mockData []*folders.Folder = []*folders.Folder{{
			Id:   *newUUID(),
			Name: "Massive",
		}}
		assert := assert.New(t)
		resp, err := folders.GetAllFoldersImproved(
			&folders.FetchFolderRequest{OrgID: *newUUID()},
			mockDb{mockData: mockData},
		)
		assert.Nil(err)
		assert.Equal(mockData, resp.Folders)
	})

	t.Run("db error", func(t *testing.T) {
		dbError := fmt.Errorf("db test error")
		assert := assert.New(t)
		resp, err := folders.GetAllFoldersImproved(
			&folders.FetchFolderRequest{OrgID: *newUUID()},
			mockDb{dbError: dbError},
		)
		assert.Equal(dbError, err)
		assert.Nil(resp)
	})

	t.Run("lots of folders", func(t *testing.T) {
		var mockData []*folders.Folder = []*folders.Folder{
			{
				Id:   *newUUID(),
				Name: "1",
			},
			{
				Id:   *newUUID(),
				Name: "2",
			},
			{
				Id:   *newUUID(),
				Name: "3",
			},
			{
				Id:   *newUUID(),
				Name: "4",
			},
			{
				Id:   *newUUID(),
				Name: "5",
			},
		}
		assert := assert.New(t)
		resp, err := folders.GetAllFoldersImproved(
			&folders.FetchFolderRequest{OrgID: *newUUID()},
			mockDb{mockData: mockData},
		)
		assert.Nil(err)
		assert.Equal(mockData, resp.Folders)
	})
}

func Test_GetAllFoldersPaginated(t *testing.T) {
	t.Run("nil request", func(t *testing.T) {
		assert := assert.New(t)

		resp, err := folders.GetAllFoldersPaginated(nil, mockDb{})
		assert.Nil(resp)
		assert.Error(err)
	})

	t.Run("db error", func(t *testing.T) {
		dbError := fmt.Errorf("db test error")
		assert := assert.New(t)
		resp, err := folders.GetAllFoldersImproved(
			&folders.FetchFolderRequest{OrgID: *newUUID()},
			mockDb{dbError: dbError},
		)
		assert.Equal(dbError, err)
		assert.Nil(resp)
	})

	t.Run("4 folders, two pages, two items each", func(t *testing.T) {
		assert := assert.New(t)
		var mockData []*folders.Folder = []*folders.Folder{
			{
				Id:   *newUUID(),
				Name: "1",
			},
			{
				Id:   *newUUID(),
				Name: "2",
			},
			{
				Id:   *newUUID(),
				Name: "3",
			},
			{
				Id:   *newUUID(),
				Name: "4",
			},
		}
		resp, err := folders.GetAllFoldersPaginated(
			&folders.FetchFolderPaginatedRequest{OrgID: *newUUID(), Size: 2},
			mockDb{mockData: mockData},
		)
		assert.Nil(err)
		assert.Equal(*resp, folders.FetchFolderPaginatedResponse{
			Folders: []*folders.Folder{mockData[0], mockData[1]},
			Next:    &(mockData[2].Id),
		})

		resp, err = folders.GetAllFoldersPaginated(
			&folders.FetchFolderPaginatedRequest{
				OrgID:      *newUUID(),
				Size:       2,
				StartingAt: resp.Next,
			},
			mockDb{mockData: mockData},
		)

		assert.Nil(err)
		assert.Equal(*resp, folders.FetchFolderPaginatedResponse{
			Folders: []*folders.Folder{mockData[2], mockData[3]},
			Next:    nil,
		})
	})

	t.Run("3 folders, 2 pages, size 2", func(t *testing.T) {
		assert := assert.New(t)
		var mockData []*folders.Folder = []*folders.Folder{
			{
				Id:   *newUUID(),
				Name: "1",
			},
			{
				Id:   *newUUID(),
				Name: "2",
			},
			{
				Id:   *newUUID(),
				Name: "3",
			},
		}
		resp, err := folders.GetAllFoldersPaginated(
			&folders.FetchFolderPaginatedRequest{OrgID: *newUUID(), Size: 2},
			mockDb{mockData: mockData},
		)
		assert.Nil(err)
		assert.Equal(*resp, folders.FetchFolderPaginatedResponse{
			Folders: []*folders.Folder{mockData[0], mockData[1]},
			Next:    &(mockData[2].Id),
		})

		resp, err = folders.GetAllFoldersPaginated(
			&folders.FetchFolderPaginatedRequest{
				OrgID:      *newUUID(),
				Size:       2,
				StartingAt: resp.Next,
			},
			mockDb{mockData: mockData},
		)

		assert.Nil(err)
		assert.Equal(*resp, folders.FetchFolderPaginatedResponse{
			Folders: []*folders.Folder{mockData[2]},
			Next:    nil,
		})
	})

	t.Run("3 folders, size 4, 1 page", func(t *testing.T) {
		assert := assert.New(t)
		var mockData []*folders.Folder = []*folders.Folder{
			{
				Id:   *newUUID(),
				Name: "1",
			},
			{
				Id:   *newUUID(),
				Name: "2",
			},
			{
				Id:   *newUUID(),
				Name: "3",
			},
		}
		resp, err := folders.GetAllFoldersPaginated(
			&folders.FetchFolderPaginatedRequest{OrgID: *newUUID(), Size: 4},
			mockDb{mockData: mockData},
		)
		assert.Nil(err)
		assert.Equal(*resp, folders.FetchFolderPaginatedResponse{
			Folders: []*folders.Folder{mockData[0], mockData[1], mockData[2]},
			Next:    nil,
		})
	})

	t.Run("3 folders, size 0, 2 pages", func(t *testing.T) {
		assert := assert.New(t)
		var mockData []*folders.Folder = []*folders.Folder{
			{
				Id:   *newUUID(),
				Name: "1",
			},
			{
				Id:   *newUUID(),
				Name: "2",
			},
			{
				Id:   *newUUID(),
				Name: "3",
			},
		}
		resp, err := folders.GetAllFoldersPaginated(
			&folders.FetchFolderPaginatedRequest{OrgID: *newUUID(), Size: 0},
			mockDb{mockData: mockData},
		)
		assert.Nil(err)
		assert.Equal(*resp, folders.FetchFolderPaginatedResponse{
			Folders: []*folders.Folder{},
			Next:    &mockData[0].Id,
		})

		resp, err = folders.GetAllFoldersPaginated(
			&folders.FetchFolderPaginatedRequest{OrgID: *newUUID(), Size: 0},
			mockDb{mockData: mockData},
		)
		assert.Nil(err)
		assert.Equal(*resp, folders.FetchFolderPaginatedResponse{
			Folders: []*folders.Folder{},
			Next:    &mockData[0].Id,
		})
	})

	t.Run("invalid staringAt", func(t *testing.T) {
		assert := assert.New(t)
		var mockData []*folders.Folder = []*folders.Folder{}
		resp, err := folders.GetAllFoldersPaginated(
			&folders.FetchFolderPaginatedRequest{
				OrgID:      *newUUID(),
				Size:       4,
				StartingAt: newUUID(),
			},
			mockDb{mockData: mockData},
		)
		assert.NotNil(err)
		assert.Nil(resp)
	})
}
