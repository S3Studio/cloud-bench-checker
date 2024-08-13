// ConstraintChecker to check the constraint of a cloud connector

package framework

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/s3studio/cloud-bench-checker/pkg/auth"
	"github.com/s3studio/cloud-bench-checker/pkg/connector"
	def "github.com/s3studio/cloud-bench-checker/pkg/definition"
)

func TestConstraintChecker_Check(t *testing.T) {
	patches := gomonkey.ApplyFunc(connector.GetK8sVersion,
		func(authProvider auth.IAuthProvider) (string, error) {
			return "1.29", nil
		})
	defer patches.Reset()

	type args struct {
		authProvider auth.IAuthProvider
		cloudType    string
	}
	tests := []struct {
		name            string
		c               *ConstraintChecker
		args            args
		wantEmptyString bool
		wantErr         bool
	}{
		{
			"Valid result",
			&ConstraintChecker{
				&def.ConfConstraint{ConstraintK8s: def.ConfConstraintK8s{Version: "1.29"}},
			},
			args{nil, string(def.K8S)},
			true,
			false,
		},
		{
			"Valid result with wildcard",
			&ConstraintChecker{
				&def.ConfConstraint{ConstraintK8s: def.ConfConstraintK8s{Version: "1.*"}},
			},
			args{nil, string(def.K8S)},
			true,
			false,
		},
		{
			"Valid result with tilde",
			&ConstraintChecker{
				&def.ConfConstraint{ConstraintK8s: def.ConfConstraintK8s{Version: "~1.x"}},
			},
			args{nil, string(def.K8S)},
			true,
			false,
		},
		{
			"Valid result with caret",
			&ConstraintChecker{
				&def.ConfConstraint{ConstraintK8s: def.ConfConstraintK8s{Version: "^1.27"}},
			},
			args{nil, string(def.K8S)},
			true,
			false,
		},
		{
			"Valid result with compare",
			&ConstraintChecker{
				&def.ConfConstraint{ConstraintK8s: def.ConfConstraintK8s{Version: ">=1.27, <=1.29"}},
			},
			args{nil, string(def.K8S)},
			true,
			false,
		},
		{
			"Valid result with different version",
			&ConstraintChecker{
				&def.ConfConstraint{ConstraintK8s: def.ConfConstraintK8s{Version: "1.27"}},
			},
			args{nil, string(def.K8S)},
			false,
			false,
		},
		{
			"Valid result with no constraint",
			&ConstraintChecker{
				&def.ConfConstraint{},
			},
			args{nil, string(def.K8S)},
			true,
			false,
		},
		{
			"Valid result with no constraint for cloud",
			&ConstraintChecker{
				&def.ConfConstraint{},
			},
			args{nil, string(def.TENCENT_CLOUD)},
			true,
			false,
		},
		{
			"failed to parse version constraint",
			&ConstraintChecker{
				&def.ConfConstraint{ConstraintK8s: def.ConfConstraintK8s{Version: "invalid"}},
			},
			args{nil, string(def.K8S)},
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Check(tt.args.authProvider, tt.args.cloudType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConstraintChecker.Check() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == "") != tt.wantEmptyString {
				t.Errorf("ConstraintChecker.Check() = %v, wantEmptyString %v", got, tt.wantEmptyString)
			}
		})
	}
}
