package folders

import "github.com/gofrs/uuid"

type FetchFolderRequest struct {
	OrgID uuid.UUID
}

type FetchFolderResponse struct {
	Folders []*Folder
}

type Db interface {
	FetchAllFoldersByOrgID(orgID uuid.UUID) ([]*Folder, error)
}

type FetchFolderPaginatedRequest struct {
	OrgID      uuid.UUID
	Size       int
	StartingAt *uuid.UUID // nil UUID will represent a "first" request
}

type FetchFolderPaginatedResponse struct {
	Folders []*Folder
	Next    *uuid.UUID // nil UUID will indicate that there is no more data
}
