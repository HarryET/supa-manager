package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

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

func (a *Api) getProjectApi(c *gin.Context) {
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
