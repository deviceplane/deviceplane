package info

import (
	//"github.com/cobaugh/osrelease"
	"github.com/deviceplane/deviceplane/pkg/models"
)

func getOSRelease() (*models.OSRelease, error) {
	//osRelease, err := osrelease.Read()
	//if err != nil {
	//	return nil, err
	//}

	return &models.OSRelease{
		PrettyName: "",//osRelease["PRETTY_NAME"],
		Name:       "",//osRelease["NAME"],
		VersionID:  "",//osRelease["VERSION_ID"],
		Version:    "",//osRelease["VERSION"],
		ID:         "",//osRelease["ID"],
		IDLike:     "",//osRelease["ID_LIKE"],
	}, nil
}
