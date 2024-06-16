package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/harryet/supa-manager/database"
	"github.com/tjarratt/babble"
	"net/http"
)

type Status = string

const (
	StatusActiveHealthy     Status = "ACTIVE_HEALTHY"
	StatusActiveUnhealthy   Status = "ACTIVE_UNHEALTHY"
	StatusInitFailed        Status = "INIT_FAILED"
	StatusUnknown           Status = "UNKNOWN"
	StatusComingUp          Status = "COMING_UP"
	StatusGoingDown         Status = "GOING_DOWN"
	StatusInactive          Status = "INACTIVE"
	StatusPausing           Status = "PAUSING"
	StatusRemoved           Status = "REMOVED"
	StatusRestoring         Status = "RESTORING"
	StatusUpgrading         Status = "UPGRADING"
	StatusCreatingProject   Status = "CREATING_PROJECT"
	StatusRunningMigrations Status = "RUNNING_MIGRATIONS"
	StatusMigrationsFailed  Status = "MIGRATIONS_FAILED"
	StatusMigrationsPassed  Status = "MIGRATIONS_PASSED"
	StatusFunctionsDeployed Status = "FUNCTIONS_DEPLOYED"
	StatusFunctionsFailed   Status = "FUNCTIONS_FAILED"
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

func (a *Api) platformProject(c *gin.Context) {
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

type CreateProjects struct {
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

func (a *Api) platformCreateProject(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var createProject CreateProjects
	if err := c.BindJSON(&createProject); err != nil {
		c.JSON(400, gin.H{"error": "Bad Request"})
		return
	}

	proj, err := a.queries.CreateProject(c.Request.Context(), database.CreateProjectParams{
		ProjectRef:     babble.NewBabbler().Babble(),
		ProjectName:    createProject.Name,
		OrganizationID: createProject.OrgId,
		JwtSecret:      uuid.New().String(),
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
		Endpoint:                 fmt.Sprintf("https://%s.supamanager.io", proj.ProjectRef),
		AnonKey:                  "a.b.c",
		ServiceKey:               "a.b.c",
		IsBranchEnabled:          false,
		PreviewBranchRefs:        []interface{}{},
		IsPhysicalBackupsEnabled: false,
	})
}

type ProjectCreationResponse struct {
	Id                       int32         `json:"id"`
	Ref                      string        `json:"ref"`
	Name                     string        `json:"name"`
	Status                   string        `json:"status"`
	OrganizationId           int32         `json:"organization_id"`
	CloudProvider            string        `json:"cloud_provider"`
	Region                   string        `json:"region"`
	InsertedAt               string        `json:"inserted_at"`
	Endpoint                 string        `json:"endpoint"`
	AnonKey                  string        `json:"anon_key"`
	ServiceKey               string        `json:"service_key"`
	IsBranchEnabled          bool          `json:"is_branch_enabled"`
	PreviewBranchRefs        []interface{} `json:"preview_branch_refs"`
	IsPhysicalBackupsEnabled bool          `json:"is_physical_backups_enabled"`
}

func (a *Api) projectStatus(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{"status": project.Status})
}

func (a *Api) projectConnections(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"connections": []interface{}{}})
}

type ProjectAutoApiService struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	AppConfig struct {
		Endpoint string `json:"endpoint"`
		DbSchema string `json:"db_schema"`
	} `json:"app_config"`
	App struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"app"`
	ServiceApiKeys []struct {
		Tags string `json:"tags"`
		Name string `json:"name"`
	} `json:"service_api_keys"`
	Protocol string `json:"protocol"`
	Endpoint string `json:"endpoint"`
	RestUrl  string `json:"restUrl"`
	Project  struct {
		Ref string `json:"ref"`
	} `json:"project"`
	DefaultApiKey string `json:"defaultApiKey"`
	ServiceApiKey string `json:"serviceApiKey"`
}

func (a *Api) propsProjectApi(c *gin.Context) {
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
		"autoApiService": ProjectAutoApiService{
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
	})
}

func (a *Api) projectJwtSecretUpdateStatus(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	projectRef := c.Param("ref")
	_, err = a.queries.GetProjectByRef(c, projectRef)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jwtSecretUpdateStatus": nil,
	})
}

func (a *Api) projectSettings(c *gin.Context) {
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
