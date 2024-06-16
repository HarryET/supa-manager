package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Api) getPlatformProjectSettings(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	projectRef := c.Param("ref")
	proj, err := a.queries.GetProjectByRef(c, projectRef)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"project": Project{
			Id:                       proj.ID,
			Ref:                      proj.ProjectRef,
			Name:                     proj.ProjectName,
			Status:                   proj.Status,
			OrganizationId:           proj.OrganizationID,
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
		},
		"services": []interface{}{
			ProjectAutoApiService{
				Id:   0,
				Name: "Default API",
				AppConfig: struct {
					Endpoint string `json:"endpoint"`
					DbSchema string `json:"db_schema"`
				}{
					Endpoint: fmt.Sprintf("%s.supamanager.io", proj.ProjectRef),
					DbSchema: "public",
				},
				App: struct {
					Id   int    `json:"id"`
					Name string `json:"name"`
				}{
					Id:   1,
					Name: "Auto API",
				},
				ServiceApiKeys: []struct {
					Tags string `json:"tags"`
					Name string `json:"name"`
				}{
					{
						Tags: "anon",
						Name: "anon key",
					},
					{
						Tags: "service_role",
						Name: "service_role key",
					},
				},
				Protocol: "https",
				Endpoint: fmt.Sprintf("%s.supamanager.io", proj.ProjectRef),
				RestUrl:  fmt.Sprintf("https://%s.supamanager.io/rest/v1/", proj.ProjectRef),
				Project: struct {
					Ref string `json:"ref"`
				}{
					Ref: proj.ProjectRef,
				},
				DefaultApiKey: "a.b.c",
				ServiceApiKey: "a.b.c",
			},
		},
	})
}
