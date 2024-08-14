package folders

import (
	"fmt"

	"github.com/gofrs/uuid"
)

/*
SOLUTION OVERVIEW:
I have used a "cursor" based pagination approach, wherein users of the function
supply a "startingAt" field to specify where they want folder results to start from.
They can also supply a Size field, whereby the size of the results will be bounded by
it.
NOTE: i have used the UUID of the folder as the "cursor" for simplicity

The function returns a Next field to the user -> which is the "cursor" location of the
next bit of data. Users can subsequently use this new "cusor" as the new "startingAt"
field for their next request.

Additionally, I supply a db interface to the function -> decouples the actual data fetching
so that I can more easily test edge cases using Dependency Injection.
*/

// GetAllFoldersPaginated retrieves all folders associated with a given
// organization ID from the database and paginates the results based on
// the request parameters.
//
// Parameters:
//   - req (*FetchFolderPaginatedRequest): A request object containing the necessary
//     parameters for pagination, such as OrgID, StartingAt, and Size.
//     The StartingAt parameter specifies the ID of the folder from which
//     the pagination should start.
//   - db (Db): The database interface that provides the method to fetch
//     all folders by organization ID.
//
// Returns:
//   - (*FetchFolderPaginatedResponse, error): A response object containing the paginated
//     folders and the ID of the next folder, if available. Returns an error
//     if any issues are encountered during execution.
func GetAllFoldersPaginated(
	req *FetchFolderPaginatedRequest, db Db,
) (*FetchFolderPaginatedResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	// NOTE: it may be better to handle pagination down in the DB layer
	// since in a large application pulling all the folders into memory
	// is costly. For the sake of this take home assessment, i haven't
	// done this.
	folders, err := db.FetchAllFoldersByOrgID(req.OrgID)
	if err != nil {
		return nil, err
	}

	startingAtIndex, err := getStartingAtIndex(folders, req.StartingAt)
	if err != nil {
		return nil, err
	}

	paginatedFolders := []*Folder{}
	i := startingAtIndex
	for i < min(len(folders), startingAtIndex+req.Size) {
		paginatedFolders = append(paginatedFolders, folders[i])
		i += 1
	}
	if i >= len(folders) {
		// no next, we are at the end of folders
		return &FetchFolderPaginatedResponse{
			Folders: paginatedFolders,
			Next:    nil,
		}, nil
	}
	return &FetchFolderPaginatedResponse{
		Folders: paginatedFolders,
		Next:    &folders[i].Id,
	}, nil
}

// getStartingAtIndex determines the index in the folders slice from which
// pagination should start based on the StartingAt UUID.
//
// Parameters:
// - folders ([]*Folder): A slice of folders retrieved from the database.
// - startingAt (*uuid.UUID): A pointer to the UUID of the folder to start at.
//
// Returns:
//   - (int, error): The index from which to start pagination, or an error if
//     the StartingAt UUID does not match any folder.
func getStartingAtIndex(folders []*Folder, startingAt *uuid.UUID) (int, error) {
	if startingAt == nil {
		return 0, nil
	}
	for i, folder := range folders {
		if folder.Id == *startingAt {
			return i, nil
		}
	}

	return 0, fmt.Errorf(
		"request StartingAt does not refer to a valid Folder",
	)
}
