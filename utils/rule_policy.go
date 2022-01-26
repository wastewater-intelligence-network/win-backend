package utils

import "github.com/shaj13/go-guardian/v2/auth"

const (
	PolicyOpen     = "$open"
	PolicyAllUsers = "$allUsers"
)

type PolicyRules []string

type Policy struct {
	EndpointRoleMap map[string]PolicyRules
}

func NewPolicy() *Policy {
	return &Policy{
		EndpointRoleMap: make(map[string]PolicyRules),
	}
}

func (pr Policy) AddRule(endpoint string, roles PolicyRules) {
	pr.EndpointRoleMap[endpoint] = roles
}

func (pr Policy) Check(requestUri string, user auth.Info) bool {
	userRoles := user.GetGroups()
	policyRules := pr.EndpointRoleMap[requestUri]

	if len(policyRules) == 0 {
		return true
	}

	for _, pRule := range policyRules {
		if pRule == PolicyOpen || pRule == PolicyAllUsers {
			return true
		}
		for _, uRule := range userRoles {
			if uRule == pRule {
				return true
			}
		}
	}
	return false
}

func (pr Policy) IsOpen(requestUri string) bool {
	for _, pRule := range pr.EndpointRoleMap[requestUri] {
		if pRule == PolicyOpen {
			return true
		}
	}
	return false
}
