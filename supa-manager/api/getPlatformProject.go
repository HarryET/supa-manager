package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Api) getPlatformProject(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	projectRef := c.Param("ref")
	project, err := a.queries.GetProjectByRef(c, projectRef)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, Project{
		Id:                       project.ID,
		Ref:                      project.ProjectRef,
		Name:                     project.ProjectName,
		Status:                   project.Status,
		OrganizationId:           project.OrganizationID,
		InsertedAt:               "",
		SubscriptionId:           "-",
		CloudProvider:            "k8s",
		Region:                   "mars-1",
		DiskVolumeSizeGb:         0,
		Size:                     "",
		DbUserSupabase:           "",
		DbPassSupabase:           "",
		DbDnsName:                "",
		DbHost:                   "",
		DbPort:                   0,
		DbName:                   "",
		SslEnforced:              false,
		WalgEnabled:              false,
		InfraComputeSize:         "",
		PreviewBranchRefs:        []interface{}{},
		IsBranchEnabled:          false,
		IsPhysicalBackupsEnabled: false,
	})
}
