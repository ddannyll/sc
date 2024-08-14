package folders

import (
	"fmt"

	"github.com/gofrs/uuid"
)

// NOTE: i have left the original function (almost) as it was and create a new one below it

/*
GetAllFolders takes in a *FetchFolderRequest and returns  (*FetchFolderResponse, error)
It uses the OrgId in the request to call FetchAllFoldersByOrgID, which does the actual
retrival of the Folders from, in this case, a set of sample data. The job of GetAllFolders
is simply to parse the request and wrap the response up nicely.
*/
func GetAllFolders(req *FetchFolderRequest) (*FetchFolderResponse, error) {
	// IMPROVEMENT: req can be nil, handle that case
	// IMPROVEMENT: err, f1, fs defined below are not used, we can get rid of them
	// NOTE: i have commented them out so the code compiles
	// var (
	// 	err error
	// 	f1  Folder
	// 	fs  []*Folder
	// )
	// IMPROVEMENT (style): f is a vague variable name, we can name this variable folders instead
	f := []Folder{}
	// IMPROVEMENT: here, (based on type signature)
	// there could potentially be an error from FetchAllFoldersByOrgID.
	// Although in the CURRENT actual implementation this can't happen,
	// will should still handle the error case.
	r, _ := FetchAllFoldersByOrgID(req.OrgID)
	// IMPROVEMENT: this function is hard to test since it is tightly coupled to
	// FetchAllFoldersByOrgID, we can instead pass (as an argument) a struct
	// which implements a FetchAllFoldersByOrgID method -> easier to test using DI
	for _, v := range r {
		// IMPROVEMENT: k is unused, can change to _ so code compiles
		f = append(f, *v)
	}
	var fp []*Folder
	for _, v1 := range f {
		// IMPROVEMENT: k1 is unused, can change to _ so code compiles
		fp = append(fp, &v1)
	}
	// IMPROVEMENT: the above code is fetching Folders as []*Folder, converting them to []Folder
	// and then converting them back to []*Folder, we can instead just direcly put the fetched folders
	// into FetchFolderResponse.
	var ffr *FetchFolderResponse
	ffr = &FetchFolderResponse{Folders: fp}
	return ffr, nil
}

// Improved version of the above function
func GetAllFoldersImproved(req *FetchFolderRequest, db Db) (*FetchFolderResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	folders, err := db.FetchAllFoldersByOrgID(req.OrgID)
	if err != nil {
		return nil, err
	}
	return &FetchFolderResponse{Folders: folders}, nil
}

/*
FetchAllFoldersByOrgID takes in an organisation UUID and fetches folders associated with
that organisation, returning them as: []*Folder.

Note in the function signature, FetchAllFoldersByOrdID "can" also return an error.
In this specific implementation, it appears that the actual retrieving of the folders is
mocked through GetSampleData(), and the function can never return an error.

However, if the implementation of fetching is changed to something like a database call
(which could potentially result in an error) the function type signature is setup so
we can pass an error back to the caller.
*/
func FetchAllFoldersByOrgID(orgID uuid.UUID) ([]*Folder, error) {
	folders := GetSampleData()

	resFolder := []*Folder{}
	for _, folder := range folders {
		if folder.OrgId == orgID {
			resFolder = append(resFolder, folder)
		}
	}
	return resFolder, nil
}
