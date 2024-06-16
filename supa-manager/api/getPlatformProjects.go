package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Api) getPlatformProjects(c *gin.Context) {
	account, err := a.GetAccountFromRequest(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	projects, err := a.queries.GetProjectsForAccountId(c, account.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	supaProjects := []Project{}
	for _, project := range projects {
		supaProjects = append(supaProjects, Project{
			Id:                       project.ID,
			Ref:                      project.ProjectRef,
			Name:                     project.ProjectName,
			Status:                   "INACTIVE",
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

	c.JSON(http.StatusOK, supaProjects)
}
