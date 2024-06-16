package api

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
