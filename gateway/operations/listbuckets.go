package operations

import (
	"net/http"
	"treeverse-lake/gateway/errors"
	"treeverse-lake/gateway/permissions"
	"treeverse-lake/gateway/serde"
)

type ListBuckets struct{}

func (controller *ListBuckets) GetArn() string {
	return "arn:treeverse:repos:::*"
}

func (controller *ListBuckets) GetPermission() string {
	return permissions.PermissionReadRepo
}

func (controller *ListBuckets) Handle(o *AuthenticatedOperation) {
	repos, err := o.Index.ListRepos(o.ClientId)
	if err != nil {
		o.EncodeError(errors.Codes.ToAPIErr(errors.ErrInternalError))
		return
	}

	// assemble response
	buckets := make([]serde.Bucket, len(repos))
	for i, repo := range repos {
		buckets[i] = serde.Bucket{
			CreationDate: serde.Timestamp(repo.GetCreationDate()),
			Name:         repo.GetRepoId(),
		}
	}
	// get client
	client, err := o.Auth.GetClient(o.ClientId)
	if err != nil {
		o.EncodeError(errors.Codes.ToAPIErr(errors.ErrInternalError))
		return
	}

	// write response
	o.EncodeResponse(serde.ListBucketsOutput{
		Buckets: serde.Buckets{Bucket: buckets},
		Owner: serde.Owner{
			DisplayName: client.GetName(),
			ID:          client.GetId(),
		},
	}, http.StatusOK)

}