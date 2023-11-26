package auth

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/persist"
)

type PolicyEnforcer struct {
	adapter  persist.BatchAdapter
	enforcer casbin.Enforcer
}

var (
	initialised    = false
	policyEnforcer = &PolicyEnforcer{}
)

func GetPolicyEnforcer() *PolicyEnforcer {
	if !initialised {
		panic("Policy enforcer not initialized!")
	}
	return policyEnforcer
}

//func InitializePolicyEnforcer() {
//	ctx := context.Background()
//
//	a := mongodb.GetCollection("casbin_rules", ctx)
//
//	e, err := casbin.NewEnforcer("internal/auth/configuration/rbac_model.conf", a)
//	if err != nil {
//		panic(err)
//	}
//
//	// Load the policy from DB.
//	err = e.LoadPolicy()
//	if err != nil {
//		panic(err)
//	}
//
//	policyEnforcer.adapter = a
//	policyEnforcer.enforcer = *e
//	utils.Log("Policy enforcer initialized")
//}

// Enforce decides whether a "subject" can access a "object" with the operation "action",
func (pe *PolicyEnforcer) Enforce(sub, obj, act string) bool {
	result, err := pe.enforcer.Enforce(sub, obj, act)
	if err != nil {
		panic(err)
	}
	return result
}

// AddPolicy adds a policy rule to the current policy.
func (pe *PolicyEnforcer) AddPolicy(sub, obj, act string) bool {
	result, err := pe.enforcer.AddPolicy(sub, obj, act)
	if err != nil {
		panic(err)
	}
	return result
}

// RemovePolicy removes a policy rule from the current policy.
func (pe *PolicyEnforcer) RemovePolicy(sub, obj, act string) bool {
	result, err := pe.enforcer.RemovePolicy(sub, obj, act)
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
