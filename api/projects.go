package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Project struct {
	Id                       int32         `json:"id"`
	Ref                      string        `json:"ref"`
	Name                     string        `json:"name"`
	Status                   string        `json:"status"`
	OrganizationId           int32         `json:"organization_id"`
	InsertedAt               string        `json:"inserted_at"`
	SubscriptionId           string        `json:"subscription_id"`
	CloudProvider            string        `json:"cloud_provider"`
	Region                   string        `json:"region"`
	DiskVolumeSizeGb         int32         `json:"disk_volume_size_gb"`
	Size                     string        `json:"size"`
	DbUserSupabase           string        `json:"db_user_supabase"`
	DbPassSupabase           string        `json:"db_pass_supabase"`
	DbDnsName                string        `json:"db_dns_name"`
	DbHost                   string        `json:"db_host"`
	DbPort                   int32         `json:"db_port"`
	DbName                   string        `json:"db_name"`
	SslEnforced              bool          `json:"ssl_enforced"`
	WalgEnabled              bool          `json:"walg_enabled"`
	InfraComputeSize         string        `json:"infra_compute_size"`
	PreviewBranchRefs        []interface{} `json:"preview_branch_refs"`
	IsBranchEnabled          bool          `json:"is_branch_enabled"`
	IsPhysicalBackupsEnabled bool          `json:"is_physical_backups_enabled"`
}

func (a *Api) platformProjects(c *gin.Context) {
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
