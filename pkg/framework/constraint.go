// ConstraintChecker to check the constraint of a cloud connector

package framework

import (
	"fmt"

	"github.com/s3studio/cloud-bench-checker/pkg/auth"
	"github.com/s3studio/cloud-bench-checker/pkg/connector"
	def "github.com/s3studio/cloud-bench-checker/pkg/definition"

	"github.com/Masterminds/semver/v3"
)

// ConstraintChecker: Used to check the constraint of a cloud connector
type ConstraintChecker struct {
	conf *def.ConfConstraint
}

// NewConstraintChecker: Constructor of ConstraintChecker
// @param: conf: Definition of Listor
func NewConstraintChecker(conf *def.ConfConstraint) *ConstraintChecker {
	constraintChecker := ConstraintChecker{conf: conf}
	return &constraintChecker
}

// Check: Check the constraint
// @param: authProvider: IAuthProvider to provide profile of auth
// @param: cloudType: Type of cloud that the constraint is associated with
// @return: Empty string if the constraint is satisfied, or description if not satisfied
// @return: Error
func (c *ConstraintChecker) Check(authProvider auth.IAuthProvider, cloudType string) (string, error) {
	switch cloudType {
	case string(def.K8S):
		if c.conf.ConstraintK8s.Version == "" {
			// constraint not set
			return "", nil
		}

		serverVersion, err := connector.GetK8sVersion(authProvider)
		if err != nil {
			return "", err
		}

		target, err := semver.NewVersion(serverVersion)
		if err != nil {
			return "", fmt.Errorf("failed to parse version: %w", err)
		}

		constraint, err := semver.NewConstraint(c.conf.ConstraintK8s.Version)
		if err != nil {
			return "", fmt.Errorf("failed to parse version constraint: %w", err)
		}

		if constraint.Check(target) {
			return "", nil
		} else {
			return fmt.Sprintf("constraint not satisfied, need %s, got %s", c.conf.ConstraintK8s.Version, serverVersion), nil
		}
	default:
		// Treated as satisfied if the cloudType has no constraint implement
		return "", nil
	}
}
