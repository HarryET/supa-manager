package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/harryet/supa-manager/conf"
	"github.com/harryet/supa-manager/database"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matthewhartstonge/argon2"
	"log/slog"
	"net/http"
	"time"
)

type Api struct {
	isHealthy bool
	logger    *slog.Logger
	config    *conf.Config
	queries   *database.Queries
	pgPool    *pgxpool.Pool
	argon     argon2.Config
}

func CreateApi(logger *slog.Logger, config *conf.Config) (*Api, error) {
	conn, err := pgxpool.New(context.Background(), config.DatabaseUrl)
	if err != nil {
		logger.Error(fmt.Sprintf("Unable to connect to database: %v", err))
		return nil, err
	}

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
		pgPool:  conn,
		argon:   argon2.DefaultConfig(),
	}, nil
}

func (a *Api) GetAccountIdFromRequest(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("missing Authorization header")
	}

	tokenString := authHeader[len("Bearer "):]
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.config.JwtSecret), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	return claims.Subject, nil
}

func (a *Api) GetAccountFromRequest(c *gin.Context) (*database.Account, error) {
	id, err := a.GetAccountIdFromRequest(c)
	if err != nil {
		return nil, err
	}

	if id == "" {
		return nil, errors.New("missing account ID")
	}

	account, err := a.queries.GetAccountByGoTrueID(c.Request.Context(), id)
	if err != nil {
		return nil, err
	}

	return &account, nil
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

func (a *Api) telemetry(c *gin.Context) {
	c.AbortWithStatus(http.StatusNoContent)
}

func (a *Api) Router() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/", a.index)
	r.GET("/status", a.status)
	r.GET("/telemetry/:event", a.telemetry)

	r.GET("/organizations", a.getOrganizations)
	r.GET("/profile", a.profile)
	r.GET("/profile/permissions", a.profilePermissions)
	r.POST("/profile/password-check", a.profilePasswordCheck)
	r.GET("/organizations/:slug/members/reached-free-project-limit", a.platformReachedFreeProjectLimit)
	r.GET("/projects/:ref/status", a.projectStatus)
	r.GET("/props/project/:ref/jwt-secret-update-status", a.projectJwtSecretUpdateStatus)

	gotrue := r.Group("/auth")
	{
		gotrue.POST("/token", a.gotrueToken)
	}

	platform := r.Group("/platform")
	{
		platform.POST("/signup", a.platformSignup)
		platform.GET("/notifications", a.platformNotifications)
		platform.GET("/stripe/invoices/overdue", a.platformOverdueInvoices)
		platform.GET("/projects", a.platformProjects)
		platform.POST("/projects", a.platformCreateProject)
		platform.GET("/projects/:ref", a.platformProject)
		platform.POST("/organizations", a.platformCreateOrganization)
		platform.GET("/organizations/:slug/billing/subscription", a.platformOrganizationBillingSubscription)
		// TODO rename
		platform.GET("/integrations/:integration/connections", a.projectConnections)
		platform.GET("/props/project/:ref/settings", a.projectSettings)
	}

	configcat := r.Group("/configcat")
	{
		configcat.GET("/configuration-files/:key/config_v5.json", a.getConfigCatConfiguration)
	}

	return r
}
