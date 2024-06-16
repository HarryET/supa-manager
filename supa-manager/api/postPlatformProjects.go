package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tjarratt/babble"
	"net/http"
	"strings"
	"supamanager.io/supa-manager/database"
)

type ProjectCreationBody struct {
	CloudProvider                  string `json:"cloud_provider"`
	OrgId                          int32  `json:"org_id"`
	Name                           string `json:"name"`
	DbPass                         string `json:"db_pass"`
	DbRegion                       string `json:"db_region"`
	CustomSupabaseInternalRequests struct {
		Ami struct {
			SearchTags struct {
				TagPostgresVersion string `json:"tag:postgresVersion"`
			} `json:"search_tags"`
		} `json:"ami"`
	} `json:"custom_supabase_internal_requests"`
	DesiredInstanceSize string `json:"desired_instance_size"`
}

type ProjectCreationResponse struct {
	Id                       int32    `json:"id"`
	Ref                      string   `json:"ref"`
	Name                     string   `json:"name"`
	Status                   string   `json:"status"`
	OrganizationId           int32    `json:"organization_id"`
	CloudProvider            string   `json:"cloud_provider"`
	Region                   string   `json:"region"`
	InsertedAt               string   `json:"inserted_at"`
	Endpoint                 string   `json:"endpoint"`
	AnonKey                  string   `json:"anon_key"`
	ServiceKey               string   `json:"service_key"`
	IsBranchEnabled          bool     `json:"is_branch_enabled"`
	PreviewBranchRefs        []string `json:"preview_branch_refs"`
	IsPhysicalBackupsEnabled bool     `json:"is_physical_backups_enabled"`
	IsReadReplicasEnabled    bool     `json:"is_read_replicas_enabled"`
	DiskVolumeSizeGb         int32    `json:"disk_volume_size_gb"`
	SubscriptionId           string   `json:"subscription_id"`
}

func (a *Api) postPlatformProjects(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var createProject ProjectCreationBody
	if err := c.BindJSON(&createProject); err != nil {
		c.JSON(400, gin.H{"error": "Bad Request"})
		return
	}

	proj, err := a.queries.CreateProject(c.Request.Context(), database.CreateProjectParams{
		ProjectRef:     strings.ToLower(babble.NewBabbler().Babble()),
		ProjectName:    createProject.Name,
		OrganizationID: createProject.OrgId,
		JwtSecret:      uuid.New().String(),
		CloudProvider:  strings.ToUpper(createProject.CloudProvider),
		Region:         strings.ToUpper(createProject.DbRegion),
	})

	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusCreated, ProjectCreationResponse{
		Id:                       proj.ID,
		Ref:                      proj.ProjectRef,
		Name:                     proj.ProjectName,
		Status:                   proj.Status,
		OrganizationId:           proj.OrganizationID,
		CloudProvider:            proj.CloudProvider,
		Region:                   proj.Region,
		InsertedAt:               proj.CreatedAt.Time.Format("2006-01-02T15:04:05.999Z"),
		Endpoint:                 fmt.Sprintf("https://%s.%s", proj.ProjectRef, a.config.Domain.Base),
		AnonKey:                  "a.b.c",
		ServiceKey:               "a.b.c",
		IsBranchEnabled:          false,
		PreviewBranchRefs:        []string{},
		IsPhysicalBackupsEnabled: false,
		IsReadReplicasEnabled:    false,
		DiskVolumeSizeGb:         0,
		SubscriptionId:           "wedontbill",
	})
}
