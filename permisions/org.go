package permisions

import (
	"encoding/json"
	"strconv"
	"strings"
)

func ConstructPermissions(orgIds []int32) interface{} {
	permissionStrings := []string{}
	for _, orgId := range orgIds {
		permissionStrings = append(permissionStrings, strings.ReplaceAll(strings.ReplaceAll(OrgString, "<ORG_ID>", strconv.Itoa(int(orgId))), "<ROLE_ID>", "1"))
	}
	permissionsJson := "[" + strings.Join(permissionStrings, ",") + "]"

	var permissions interface{}
	err := json.Unmarshal([]byte(permissionsJson), &permissions)
	if err != nil {
		return nil
	}
	return permissions
}

const OrgString = `{
  "organization_id": <ORG_ID>,
  "resources": ["projects", "integrations.vercel_connections", "integrations.github_connections"],
  "actions": ["write:Create"],
  "condition": null,
  "restrictive": false,
  "project_ids": []
}, {
  "organization_id": <ORG_ID>,
  "resources": ["projects"],
  "actions": ["write:Update"],
  "condition": {
    "and": [{
      "var": "resource.project_id"
    }]
  },
  "restrictive": false,
  "project_ids": []
}, {
  "organization_id": <ORG_ID>,
  "resources": ["preview_branches", "approved_oauth_apps", "third_party_auth"],
  "actions": ["write:Create", "write:Update", "write:Delete"],
  "condition": null,
  "restrictive": false,
  "project_ids": []
}, {
  "organization_id": <ORG_ID>,
  "resources": ["projects.pgsodium_root_key_encrypted"],
  "actions": ["read:Read", "write:Update"],
  "condition": null,
  "restrictive": false,
  "project_ids": []
}, {
  "organization_id": <ORG_ID>,
  "resources": ["%"],
  "actions": ["billing:Write", "infra:Execute"],
  "condition": null,
  "restrictive": false,
  "project_ids": []
}, {
  "organization_id": <ORG_ID>,
  "resources": ["notifications", "integrations.vercel_connections", "integrations.github_connections"],
  "actions": ["write:Update", "write:Delete"],
  "condition": null,
  "restrictive": false,
  "project_ids": []
}, {
  "organization_id": <ORG_ID>,
  "resources": ["preview_branches"],
  "actions": ["write:Create", "write:Update", "write:Delete"],
  "condition": {
    "!": {
      "var": "resource.is_default"
    }
  },
  "restrictive": false,
  "project_ids": []
}, {
  "organization_id": <ORG_ID>,
  "resources": ["back_ups", "events"],
  "actions": ["write:Create"],
  "condition": null,
  "restrictive": false,
  "project_ids": []
}, {
  "organization_id": <ORG_ID>,
  "resources": ["%"],
  "actions": ["infra:Execute"],
  "condition": {
    "and": [{
      "var": "resource_name"
    }, {
      "!": {
        "in": [{
          "var": "resource_name"
        }, ["queue_jobs.projects.initialize_or_resume", "queue_jobs.projects.pause"]]
      }
    }]
  },
  "restrictive": false,
  "project_ids": []
}, {
  "organization_id": <ORG_ID>,
  "resources": ["user_content"],
  "actions": ["write:Update", "write:Delete"],
  "condition": {
    "and": [{
      "var": "resource.visibility"
    }, {
      "var": "resource.owner_id"
    }, {
      "var": "subject.id"
    }, {
      "or": [{
        "!=": [{
          "var": "resource.visibility"
        }, "user"]
      }, {
        "==": [{
          "var": "resource.owner_id"
        }, {
          "var": "subject.id"
        }]
      }]
    }]
  },
  "restrictive": false,
  "project_ids": []
}, {
  "organization_id": <ORG_ID>,
  "resources": ["user_content"],
  "actions": ["write:Create"],
  "condition": {
    "and": [{
      "var": "resource.owner_id"
    }, {
      "var": "subject.id"
    }, {
      "==": [{
        "var": "resource.owner_id"
      }, {
        "var": "subject.id"
      }]
    }]
  },
  "restrictive": false,
  "project_ids": []
}, {
  "organization_id": <ORG_ID>,
  "resources": ["field.jwt_secret", "service_api_keys.service_role_key"],
  "actions": ["read:Read"],
  "condition": null,
  "restrictive": false,
  "project_ids": []
}, {
  "organization_id": <ORG_ID>,
  "resources": ["%"],
  "actions": ["auth:Execute", "functions:Write", "storage:Admin:Write", "tenant:Sql:Admin:Write", "tenant:Sql:CreateTable", "tenant:Sql:Query", "tenant:Sql:Write:%"],
  "condition": null,
  "restrictive": false,
  "project_ids": []
}, {
  "organization_id": <ORG_ID>,
  "resources": ["custom_config_gotrue", "custom_config_postgrest", "owner_reassign", "services"],
  "actions": ["write:Create", "write:Update"],
  "condition": null,
  "restrictive": false,
  "project_ids": []
}, {
  "organization_id": <ORG_ID>,
  "resources": ["members", "organizations", "auth.subject_roles", "users", "user_invites", "auth.permissions", "auth.roles"],
  "actions": ["read:Read"],
  "condition": null,
  "restrictive": false,
  "project_ids": []
}, {
  "organization_id": <ORG_ID>,
  "resources": ["auth.subject_roles"],
  "actions": ["write:Delete"],
  "condition": {
    "and": [{
      "var": "resource.subject_id"
    }, {
      "var": "subject.gotrue_id"
    }, {
      "==": [{
        "var": "resource.subject_id"
      }, {
        "var": "subject.gotrue_id"
      }]
    }]
  },
  "restrictive": false,
  "project_ids": []
}, {
  "organization_id": <ORG_ID>,
  "resources": ["user_invites", "auth.permissions", "auth.roles"],
  "actions": ["write:Create", "write:Update", "write:Delete"],
  "condition": null,
  "restrictive": false,
  "project_ids": []
}, {
  "organization_id": <ORG_ID>,
  "resources": ["auth.subject_roles"],
  "actions": ["write:Create", "write:Delete"],
  "condition": null,
  "restrictive": false,
  "project_ids": []
}, {
  "organization_id": <ORG_ID>,
  "resources": ["organizations"],
  "actions": ["write:Update"],
  "condition": null,
  "restrictive": false,
  "project_ids": []
}, {
  "organization_id": <ORG_ID>,
  "resources": ["back_ups", "custom_config_gotrue", "custom_config_postgrest", "customers", "events", "gotrue_config", "infrastructure", "invoices", "member_active_free_projects", "notifications", "organizations", "owner_reassign", "physical_backups", "postgrest_config", "service_api_keys", "services", "stats_daily_projects", "subscriptions", "subscription_items", "preview_branches", "resource_exhaustion_notifications", "approved_oauth_apps", "third_party_auth", "integrations.vercel_connections", "integrations.github_connections"],
  "actions": ["read:Read"],
  "condition": null,
  "restrictive": false,
  "project_ids": []
}, {
  "organization_id": <ORG_ID>,
  "resources": ["projects", "third_party_auth"],
  "actions": ["read:Read"],
  "condition": {
    "and": [{
      "var": "resource.project_id"
    }]
  },
  "restrictive": false,
  "project_ids": []
}, {
  "organization_id": <ORG_ID>,
  "resources": ["user_content_folders"],
  "actions": ["read:Read", "write:Create", "write:Update", "write:Delete"],
  "condition": {
    "and": [{
      "var": "resource.owner_id"
    }, {
      "var": "subject.id"
    }, {
      "==": [{
        "var": "resource.owner_id"
      }, {
        "var": "subject.id"
      }]
    }]
  },
  "restrictive": false,
  "project_ids": []
}, {
  "organization_id": <ORG_ID>,
  "resources": ["user_content"],
  "actions": ["write:Update", "write:Delete"],
  "condition": {
    "and": [{
      "var": "resource.owner_id"
    }, {
      "var": "resource.type"
    }, {
      "var": "resource.visibility"
    }, {
      "var": "subject.id"
    }, {
      "and": [{
        "!==": [{
          "var": "resource.type"
        }, "report"]
      }, {
        "===": [{
          "var": "resource.owner_id"
        }, {
          "var": "subject.id"
        }]
      }]
    }]
  },
  "restrictive": false,
  "project_ids": []
}, {
  "organization_id": <ORG_ID>,
  "resources": ["user_content"],
  "actions": ["write:Create"],
  "condition": {
    "and": [{
      "var": "resource.owner_id"
    }, {
      "var": "resource.type"
    }, {
      "var": "subject.id"
    }, {
      "and": [{
        "!==": [{
          "var": "resource.type"
        }, "report"]
      }, {
        "===": [{
          "var": "resource.owner_id"
        }, {
          "var": "subject.id"
        }]
      }]
    }]
  },
  "restrictive": false,
  "project_ids": []
}, {
  "organization_id": <ORG_ID>,
  "resources": ["user_content"],
  "actions": ["read:Read"],
  "condition": {
    "and": [{
      "var": "resource.visibility"
    }, {
      "var": "resource.owner_id"
    }, {
      "var": "subject.id"
    }, {
      "or": [{
        "!=": [{
          "var": "resource.visibility"
        }, "user"]
      }, {
        "==": [{
          "var": "resource.owner_id"
        }, {
          "var": "subject.id"
        }]
      }]
    }]
  },
  "restrictive": false,
  "project_ids": []
}, {
  "organization_id": <ORG_ID>,
  "resources": ["%"],
  "actions": ["analytics:Read", "billing:Read", "functions:Read", "storage:Admin:Read", "tenant:Sql:Admin:Read", "tenant:Sql:Read:Select"],
  "condition": null,
  "restrictive": false,
  "project_ids": []
}, {
  "organization_id": <ORG_ID>,
  "resources": ["auth.subject_roles"],
  "actions": ["write:Create", "write:Delete"],
  "condition": {
    "and": [{
      "var": "resource.role_id"
    }, {
      "!==": [{
        "var": "resource.role_id"
      }, <ROLE_ID>]
    }]
  },
  "restrictive": false,
  "project_ids": null
}, {
  "organization_id": <ORG_ID>,
  "resources": ["user_invites"],
  "actions": ["write:Create", "write:Update", "write:Delete"],
  "condition": {
    "and": [{
      "var": "resource.role_id"
    }, {
      "!==": [{
        "var": "resource.role_id"
      }, <ROLE_ID>]
    }]
  },
  "restrictive": false,
  "project_ids": null
}`
