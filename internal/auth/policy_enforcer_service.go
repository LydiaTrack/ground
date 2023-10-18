package auth

import (
	"context"
	"embed"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/persist"
	mongodbadapter "github.com/casbin/mongodb-adapter/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
	"lydia-track-base/internal/domain/auth"
	"lydia-track-base/internal/mongodb"
	"lydia-track-base/internal/utils"
	"os"
)

type PolicyEnforcer struct {
	adapter  persist.BatchAdapter
	enforcer casbin.Enforcer
}

var (
	initialised    = false
	policyEnforcer = &PolicyEnforcer{}
)

//go:embed configuration/rbac_model.conf
var content embed.FS

func InitializePolicyEnforcer() {
	ctx := context.Background()
	container := mongodb.GetContainer()

	host, err := container.Host(ctx)
	if err != nil {
		panic(err)
	}

	port, err := container.MappedPort(ctx, "27017")
	if err != nil {
		panic(err)
	}

	mongoClientOption := options.Client().ApplyURI("mongodb://" + host + ":" + port.Port())
	databaseName := os.Getenv("LYDIA_DB_NAME")

	a, err := mongodbadapter.NewAdapterWithCollectionName(mongoClientOption, databaseName, "casbin_rules")
	if err != nil {
		panic(err)
	}

	e, err := casbin.NewEnforcer("internal/auth/configuration/rbac_model.conf", a)
	if err != nil {
		panic(err)
	}

	// Load the policy from DB.
	err = e.LoadPolicy()
	if err != nil {
		panic(err)
	}

	policyEnforcer.adapter = a
	policyEnforcer.enforcer = *e
	initialised = true
	utils.Log("Policy enforcer initialized")
}

func GetPolicyEnforcer() *PolicyEnforcer {
	if !initialised {
		panic("Policy enforcer not initialized!")
	}
	return policyEnforcer
}

// Enforce decides whether a "subject" can access a "object" with the operation "action",
func (pe *PolicyEnforcer) Enforce(policy auth.SecurityPolicy) bool {
	result, err := pe.enforcer.Enforce(policy.Subject, policy.Object, policy.Action)
	if err != nil {
		panic(err)
	}
	return result
}

// AddPolicy adds a policy rule to the current policy.
func (pe *PolicyEnforcer) AddPolicy(policy auth.SecurityPolicy) bool {
	result, err := pe.enforcer.AddPolicy(policy.Subject, policy.Object, policy.Action)
	if err != nil {
		panic(err)
	}
	return result
}

// RemovePolicy removes a policy rule from the current policy.
func (pe *PolicyEnforcer) RemovePolicy(policy auth.SecurityPolicy) bool {
	result, err := pe.enforcer.RemovePolicy(policy.Subject, policy.Object, policy.Action)
	if err != nil {
		panic(err)
	}
	return result
}

// SavePolicy saves the current policy (usually after changed with Casbin API) back to file/database.
func (pe *PolicyEnforcer) SavePolicy() {
	err := pe.enforcer.SavePolicy()
	if err != nil {
		panic(err)
	}
}
