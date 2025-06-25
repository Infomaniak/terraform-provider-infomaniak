package dbaas

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type ImportIds struct {
	PublicCloudId        int
	PublicCloudProjectId int
	DbaasId              int
	Id                   string
}

func parseBackupRestoreImport(req resource.ImportStateRequest) (*ImportIds, error) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 4 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" || idParts[3] == "" {
		return nil, fmt.Errorf("expected import identifier with format: public_cloud_id,public_cloud_project_id,dbaas_id,id. got: %q", req.ID)
	}

	var errorList error

	publicCloudId, err := strconv.ParseInt(idParts[0], 10, 64)
	errorList = errors.Join(errorList, err)
	publicCloudProjectId, err := strconv.ParseInt(idParts[1], 10, 64)
	errorList = errors.Join(errorList, err)
	dbaasId, err := strconv.ParseInt(idParts[2], 10, 64)
	errorList = errors.Join(errorList, err)
	id := idParts[3]

	if errorList != nil {
		return nil, fmt.Errorf("expected import identifier with format: public_cloud_id,public_cloud_project_id,dbaas_id,id. got: %q", req.ID)
	}

	return &ImportIds{
		PublicCloudId:        int(publicCloudId),
		PublicCloudProjectId: int(publicCloudProjectId),
		DbaasId:              int(dbaasId),
		Id:                   id,
	}, nil
}
