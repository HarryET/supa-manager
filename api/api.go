package api

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/harryet/supa-manager/conf"
	"github.com/harryet/supa-manager/database"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"net/http"
	"time"
)

type Api struct {
	isHealthy bool
	logger    *slog.Logger
	config    *conf.Config
	queries   *database.Queries
	pg        *pgx.Conn
}

func CreateApi(logger *slog.Logger, config *conf.Config) (*Api, error) {
	conn, err := pgx.Connect(context.Background(), config.DatabaseUrl)
	if err != nil {
		logger.Error(fmt.Sprintf("Unable to connect to database: %v", err))
		return nil, err
	}
	defer conn.Close(context.Background())

	if err := conf.EnsureMigrationsTableExists(conn); err != nil {
		logger.Error(fmt.Sprintf("Failed to ensure migrations table: %v", err))
		return nil, err
	}

	queries := database.New(conn)

	if success, err := conf.EnsureMigrations(conn, queries); err != nil || !success {
		logger.Error(fmt.Sprintf("Failed to run migrations: %v", err))
		return nil, err
	}

	return &Api{
		logger:  logger,
		config:  config,
		queries: queries,
		pg:      conn,
	}, nil
}

func (a *Api) ListenAddress() string {
	return ":8080"
}

func (a *Api) index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (a *Api) status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"is_healthy": a.isHealthy})
}

func (a *Api) gotrueAuthorize(c *gin.Context) {
	// TODO: Implement OAuth Support Optionally
	// Currently we get the referer and redirect back to it without doing any auth
	referer := c.GetHeader("Referer")
	c.Redirect(http.StatusFound, referer)
}

func (a *Api) gotrueToken(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (a *Api) Router() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/", a.index)
	r.GET("/status", a.status)

	gotrue := r.Group("/auth/v1")
	{
		gotrue.POST("/authorize", a.gotrueAuthorize)
		gotrue.POST("/token", a.gotrueToken)
	}

	return r
}
