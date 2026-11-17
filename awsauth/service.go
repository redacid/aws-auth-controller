/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package awsauth

import (
	"errors"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/client-go/kubernetes"
)

// ServiceConfig is the configuration for a Service object.
type ServiceConfig struct {
	KubeClient    kubernetes.Interface
	Log           logr.Logger
	MaxRetryCount int
	MaxRetryTime  time.Duration
	MinRetryTime  time.Duration
	WithRetries   bool
}

// Service provides aws-auth configmap management behavior.
type Service interface {
	// UpsertMapRole upserts a MapRole into the configmap keyed by username.
	UpsertMapRole(mapRole MapRole) error

	// RemoveMapRole removes a MapRole from the configmap by keyed by username
	RemoveMapRole(mapRole MapRole) error

	// UpsertMapUser upserts a MapUser into the configmap keyed by username.
	UpsertMapUser(mapUser MapUser) error

	// RemoveMapUser removes a MapUser from the configmap keyed by username
	RemoveMapUser(mapUser MapUser) error

	// CheckMapAccountExists checks if an MapAccount with the specified account ID exists.
	CheckMapAccountExists(mapAccount MapAccount) error

	// CheckMapUserExists checks if an MapUser with the specified username and userarn ID exists.
	CheckMapUserExists(mapUser MapUser) error

	// CheckMapRoleExists checks if an MapRole with the specified username and rolearn ID exists.
	CheckMapRoleExists(mapRole MapRole) error

	// UpsertMapAccount upserts a mapAccount into the configmap keyed by username.
	UpsertMapAccount(mapAccount MapAccount) error

	// RemoveMapAccount removes a mapAccount from the configmap keyed by username
	RemoveMapAccount(mapAccount MapAccount) error
}

// NewService returns an implementation of the Service interface.
func NewService(cfg *ServiceConfig) (Service, error) {
	if cfg.WithRetries {
		if cfg.MaxRetryCount < 1 {
			return nil, errors.New("retry max count config must be greater than zero")
		}
	}
	return impl{cfg: *cfg}, nil
}

type impl struct {
	cfg ServiceConfig
}

func (svc impl) CheckMapRoleExists(mapRole MapRole) error {
	svc.cfg.Log.V(1).Info("CheckMapRoleExists", "username", mapRole.Username, "rolearn", mapRole.RoleARN)
	mapper := NewMapper(svc.cfg.KubeClient)
	err := mapper.CheckExists(&Arguments{
		DataType:      MapRoleData,
		Username:      mapRole.Username,
		UserARN:       mapRole.RoleARN,
		WithRetries:   svc.cfg.WithRetries,
		MaxRetryCount: svc.cfg.MaxRetryCount,
		MaxRetryTime:  svc.cfg.MaxRetryTime,
		MinRetryTime:  svc.cfg.MinRetryTime,
	})
	if err != nil {
		svc.cfg.Log.Info("username or rolearn already exists", "username", mapRole.Username, "rolearn", mapRole.RoleARN)
	}
	return err
}

// UpsertMapRole upserts a MapRole into the configmap keyed by username.
func (svc impl) UpsertMapRole(mapRole MapRole) error {
	svc.cfg.Log.V(1).Info("UpsertMapRole", "username", mapRole.Username, "rolearn", mapRole.RoleARN)
	mapper := NewMapper(svc.cfg.KubeClient)
	err := mapper.Upsert(&Arguments{
		DataType:      MapRoleData,
		RoleARN:       mapRole.RoleARN,
		Username:      mapRole.Username,
		Groups:        mapRole.Groups,
		WithRetries:   svc.cfg.WithRetries,
		MaxRetryCount: svc.cfg.MaxRetryCount,
		MaxRetryTime:  svc.cfg.MaxRetryTime,
		MinRetryTime:  svc.cfg.MinRetryTime,
	})
	if err != nil {
		svc.cfg.Log.Error(err, "failure to upsert mapRole",
			"username", mapRole.Username,
			"rolearn", mapRole.RoleARN,
			"groups", mapRole.Groups,
		)
	}
	return err
}

// RemoveMapRole removes a MapRole from the configmap keyed by username.
func (svc impl) RemoveMapRole(mapRole MapRole) error {
	svc.cfg.Log.V(1).Info("RemoveMapRole", "username", mapRole.Username, "rolearn", mapRole.RoleARN)
	mapper := NewMapper(svc.cfg.KubeClient)
	err := mapper.Remove(&Arguments{
		DataType:      MapRoleData,
		Username:      mapRole.Username,
		RoleARN:       mapRole.RoleARN,
		Groups:        mapRole.Groups,
		WithRetries:   svc.cfg.WithRetries,
		MaxRetryCount: svc.cfg.MaxRetryCount,
		MaxRetryTime:  svc.cfg.MaxRetryTime,
		MinRetryTime:  svc.cfg.MinRetryTime,
	})
	if err != nil {
		svc.cfg.Log.Info("mapRole not found",
			"username", mapRole.Username,
			"rolearn", mapRole.RoleARN,
			"groups", mapRole.Groups)
	}
	return err
}

func (svc impl) CheckMapUserExists(mapUser MapUser) error {
	svc.cfg.Log.V(1).Info("CheckMapUserExists", "username", mapUser.Username, "userarn", mapUser.UserARN)
	mapper := NewMapper(svc.cfg.KubeClient)
	err := mapper.CheckExists(&Arguments{
		DataType:      MapUserData,
		Username:      mapUser.Username,
		UserARN:       mapUser.UserARN,
		WithRetries:   svc.cfg.WithRetries,
		MaxRetryCount: svc.cfg.MaxRetryCount,
		MaxRetryTime:  svc.cfg.MaxRetryTime,
		MinRetryTime:  svc.cfg.MinRetryTime,
	})
	if err != nil {
		svc.cfg.Log.Info("username or userarn already exists", "username", mapUser.Username, "userarn", mapUser.UserARN)
	}
	return err
}

// UpsertMapUser upserts a MapUser into the configmap keyed by username.
func (svc impl) UpsertMapUser(mapUser MapUser) error {
	svc.cfg.Log.V(1).Info("UpsertMapUser", "username", mapUser.Username, "userarn", mapUser.UserARN)
	mapper := NewMapper(svc.cfg.KubeClient)
	err := mapper.Upsert(&Arguments{
		DataType:      MapUserData,
		UserARN:       mapUser.UserARN,
		Username:      mapUser.Username,
		Groups:        mapUser.Groups,
		WithRetries:   svc.cfg.WithRetries,
		MaxRetryCount: svc.cfg.MaxRetryCount,
		MaxRetryTime:  svc.cfg.MaxRetryTime,
		MinRetryTime:  svc.cfg.MinRetryTime,
	})
	if err != nil {
		svc.cfg.Log.Error(err, "failure to upsert mapUser", "username", mapUser.Username, "userarn", mapUser.UserARN)
	}
	return err
}

// RemoveMapUser removes a MapUser from the configmap keyed by username.
func (svc impl) RemoveMapUser(mapUser MapUser) error {
	svc.cfg.Log.V(1).Info("RemoveMapUser", "username", mapUser.Username, "userarn", mapUser.UserARN)
	mapper := NewMapper(svc.cfg.KubeClient)
	err := mapper.Remove(&Arguments{
		DataType:      MapUserData,
		Username:      mapUser.Username,
		UserARN:       mapUser.UserARN,
		WithRetries:   svc.cfg.WithRetries,
		MaxRetryCount: svc.cfg.MaxRetryCount,
		MaxRetryTime:  svc.cfg.MaxRetryTime,
		MinRetryTime:  svc.cfg.MinRetryTime,
	})
	if err != nil {
		svc.cfg.Log.Info("mapUser not found", "username", mapUser.Username, "userarn", mapUser.UserARN)
	}
	return err
}

func (svc impl) CheckMapAccountExists(mapAccount MapAccount) error {
	svc.cfg.Log.V(1).Info("CheckAccountExists", "accountid", mapAccount.AccountID)
	mapper := NewMapper(svc.cfg.KubeClient)
	err := mapper.CheckExists(&Arguments{
		DataType:      MapAccountData,
		AccountID:     mapAccount.AccountID,
		WithRetries:   svc.cfg.WithRetries,
		MaxRetryCount: svc.cfg.MaxRetryCount,
		MaxRetryTime:  svc.cfg.MaxRetryTime,
		MinRetryTime:  svc.cfg.MinRetryTime,
	})
	if err != nil {
		svc.cfg.Log.Info("mapAccount already exists", "accountid", mapAccount.AccountID)
	}
	return err
}

// UpsertMapAccount upserts a MapAccount into the configmap keyed by username.
func (svc impl) UpsertMapAccount(mapAccount MapAccount) error {
	svc.cfg.Log.V(1).Info("UpsertMapAccount", "accountid", mapAccount.AccountID)
	mapper := NewMapper(svc.cfg.KubeClient)
	err := mapper.Upsert(&Arguments{
		DataType:      MapAccountData,
		AccountID:     mapAccount.AccountID,
		WithRetries:   svc.cfg.WithRetries,
		MaxRetryCount: svc.cfg.MaxRetryCount,
		MaxRetryTime:  svc.cfg.MaxRetryTime,
		MinRetryTime:  svc.cfg.MinRetryTime,
	})
	if err != nil {
		svc.cfg.Log.Error(err, "failure to upsert mapAccount", "accountid", mapAccount.AccountID)
	}
	return err
}

// RemoveMapAccount removes a MapAccount from the configmap keyed by account id.
func (svc impl) RemoveMapAccount(mapAccount MapAccount) error {
	svc.cfg.Log.V(1).Info("RemoveMapAccount", "accountid", mapAccount.AccountID)
	mapper := NewMapper(svc.cfg.KubeClient)
	err := mapper.Remove(&Arguments{
		DataType:      MapAccountData,
		AccountID:     mapAccount.AccountID,
		WithRetries:   svc.cfg.WithRetries,
		MaxRetryCount: svc.cfg.MaxRetryCount,
		MaxRetryTime:  svc.cfg.MaxRetryTime,
		MinRetryTime:  svc.cfg.MinRetryTime,
	})
	if err != nil {
		svc.cfg.Log.Info("mapAccount not found", "accountid", mapAccount.AccountID)
	}
	return err
}
