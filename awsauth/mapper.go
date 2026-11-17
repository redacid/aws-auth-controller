/*

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
	"fmt"
	"reflect"
	"time"

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"k8s.io/client-go/kubernetes"
)

// Arguments are the arguments for management of the auth map.
type Arguments struct {
	OperationType OperationType
	DataType      DataType
	AccountID     string
	RoleARN       string
	UserARN       string
	Username      string
	Groups        []string
	WithRetries   bool
	MinRetryTime  time.Duration
	MaxRetryTime  time.Duration
	MaxRetryCount int
}

var logger = ctrl.Log.WithName("mapper")

// Validate validates if all Arguments fields are valid.
func (args *Arguments) Validate() {

	if args.WithRetries && args.MaxRetryCount < 1 {
		logger.Info("error: retry max count is invalid, must be greater than zero", "MaxRetryCount", args.MaxRetryCount)
	}
	if args.Username == "" {
		logger.Info("error: username not provided")
	}
	if args.OperationType == "" {
		logger.Info("error: operation type not provided")
	}
	if args.OperationType != UpsertOperation && args.OperationType != RemoveOperation {
		logger.Info("error: operation type '%s' not valid\n", args.OperationType)
	}
	if args.DataType == "" {
		logger.Info("error: data type not provided")
	}
	if args.DataType != MapRoleData && args.DataType != MapUserData && args.DataType != MapAccountData {
		logger.Info("error: data type '%s' not valid\n", args.DataType)
	}
	if args.OperationType == UpsertOperation && args.DataType == MapRoleData && args.RoleARN == "" {
		logger.Info("error: role arn not provided")
	}
	if args.OperationType == UpsertOperation && args.DataType == MapUserData && args.UserARN == "" {
		logger.Info("error: user arn not provided")
	}
}

// OperationType indicates the auth map management operation.
type OperationType string

const (
	UpsertOperation OperationType = "upsert"
	RemoveOperation OperationType = "remove"
	// CheckExistsOperation OperationType = "checkExists"
)

// DataType indicates the auth map management scope.
type DataType string

const (
	MapRoleData    DataType = "mapRole"
	MapUserData    DataType = "mapUser"
	MapAccountData DataType = "mapAccount"
)

// NewMapper returns a new Mapper object.
func NewMapper(client kubernetes.Interface) *Mapper {
	var mapper = &Mapper{
		KubernetesClient: client,
		Log:              logf.Log.WithName("mapper"),
	}

	return mapper

}

// Mapper is responsible for managing the auth map.
type Mapper struct {
	KubernetesClient kubernetes.Interface
	Log              logr.Logger
}

// Remove removes a mapRole or mapUser from the auth map.
func (m *Mapper) Remove(args *Arguments) error {
	// TODO Remove this validate
	// args.Validate()
	if args.WithRetries {
		return WithRetry(m.removeAuth, args)
	}
	return m.removeAuth(args)
}

func (m *Mapper) removeAuth(args *Arguments) error {
	authData, configMap, err := ReadAuthMap(m.KubernetesClient)
	if err != nil {
		return err
	}

	var removed bool

	if args.DataType == MapRoleData {
		var newRolesAuthMap []*MapRole
		for _, mapRole := range authData.MapRoles {
			if args.Username != mapRole.Username {
				newRolesAuthMap = append(newRolesAuthMap, mapRole)
			} else {
				removed = true
			}
		}
		authData.SetMapRoles(newRolesAuthMap)
	}

	if args.DataType == MapUserData {
		var newUsersAuthMap []*MapUser
		for _, mapUser := range authData.MapUsers {
			if args.Username != mapUser.Username && args.UserARN != mapUser.UserARN {
				newUsersAuthMap = append(newUsersAuthMap, mapUser)
			} else {
				removed = true
			}
		}
		authData.SetMapUsers(newUsersAuthMap)
	}

	if args.DataType == MapAccountData {
		var newAccountsAuthMap []*MapAccount
		for _, mapAccount := range authData.MapAccounts {
			if args.AccountID != mapAccount.AccountID {
				newAccountsAuthMap = append(newAccountsAuthMap, mapAccount)
			} else {
				removed = true
			}
		}
		authData.SetMapAccounts(newAccountsAuthMap)
	}

	if !removed {
		return fmt.Errorf("%s with fields '%v' not found in auth map", args.DataType, args)
	}
	return UpdateAuthMap(m.KubernetesClient, authData, configMap)
}

func (m *Mapper) CheckExists(args *Arguments) error {

	if args.WithRetries {
		return WithRetry(m.existsAuth, args)
	}
	return m.existsAuth(args)
}

func (m *Mapper) existsAuth(args *Arguments) error {
	authData, _, err := ReadAuthMap(m.KubernetesClient)
	if err != nil {
		return err
	}

	if args.DataType == MapAccountData {
		m.Log.Info("existsAuth", "DataType", args.DataType, "AccountID", args.AccountID)
		mapAccount := NewMapAccount(args.AccountID)
		err, exists := existsAccount(authData.MapAccounts, mapAccount)
		if exists {
			m.Log.Info("%v", err)
			return err
		} else {
			m.Log.Info("MapAccount not exists", "DataType", args.DataType, "AccountID", args.AccountID)
			return nil
		}
	}
	if args.DataType == MapUserData {
		m.Log.Info("existsAuth", "DataType", args.DataType, "Username", args.Username, "UserARN", args.UserARN)
		mapUser := NewMapUser(args.UserARN, args.Username, args.Groups)
		err, exists := existsUser(authData.MapUsers, mapUser)
		if exists {
			m.Log.Info("%v", err)
			return err
		} else {
			m.Log.Info("MapUser not exists", "DataType", args.DataType, "Username", args.Username, "UserARN", args.UserARN)
			return nil
		}
	}

	if args.DataType == MapRoleData {
		m.Log.Info("existsAuth", "DataType", args.DataType, "Username", args.Username, "RoleARN", args.RoleARN)
		mapRole := NewMapRole(args.RoleARN, args.Username, args.Groups)
		err, exists := existsRole(authData.MapRoles, mapRole)
		if exists {
			m.Log.Info("%v", err)
			return err
		} else {
			m.Log.Info("MapRole not exists", "DataType", args.DataType, "Username", args.Username, "RoleARN", args.RoleARN)
			return nil
		}
	}

	return nil
}

func existsAccount(authMaps []*MapAccount, resource *MapAccount) (error, bool) {
	var found = false
	logger.Info("existsAccount", "AccountID", resource.AccountID)
	for _, existing := range authMaps {
		logger.Info("existsAccount.Compare", "cm account id", existing.AccountID, "new account id", resource.AccountID)
		if existing.AccountID == resource.AccountID {
			found = true
			return fmt.Errorf("account with id '%s' already exists", resource.AccountID), found
		} else {
			found = false
		}
	}
	return nil, found
}

func existsUser(authMaps []*MapUser, resource *MapUser) (error, bool) {
	var found = false
	logger.Info("existsUser", "Username", resource.Username, "UserARN", resource.UserARN)
	for _, existing := range authMaps {
		logger.Info("existsUser.Compare", "cm username", existing.Username, "new username", resource.Username)
		logger.Info("existsUser.Compare", "cm userarn", existing.UserARN, "new userarn", resource.UserARN)
		if existing.Username == resource.Username {
			found = true
			return fmt.Errorf("existsUser: username  '%s' already exists", resource.Username), found
		} else if existing.UserARN == resource.UserARN {
			found = true
			return fmt.Errorf("existsUser: userarn  '%s' already exists", resource.UserARN), found
		} else {
			found = false
		}
	}
	return nil, found
}

func existsRole(authMaps []*MapRole, resource *MapRole) (error, bool) {
	var found = false
	logger.Info("existsRole",
		"Username", resource.Username,
		"RoleARN", resource.RoleARN)
	for _, existing := range authMaps {
		logger.Info("existsRole.Compare",
			"cm username", existing.Username,
			"new username", resource.Username)
		logger.Info("existsRole.Compare",
			"cm rolearn", existing.RoleARN,
			"new rolearn", resource.RoleARN)
		if existing.Username == resource.Username {
			found = true
			return fmt.Errorf("existsRole: username  '%s' already exists", resource.Username), found
		} else if existing.RoleARN == resource.RoleARN {
			found = true
			return fmt.Errorf("existsRole: rolearn  '%s' already exists", resource.RoleARN), found
		} else {
			found = false
		}
	}
	return nil, found
}

// Upsert updates or inserts a mapRole or mapUser item into the auth map.
func (m *Mapper) Upsert(args *Arguments) error {
	// TODO Remove this validate
	// args.Validate()
	if args.WithRetries {
		return WithRetry(m.upsertAuth, args)
	}
	return m.upsertAuth(args)
}

func (m *Mapper) upsertAuth(args *Arguments) error {
	authData, configMap, err := ReadAuthMap(m.KubernetesClient)
	if err != nil {
		return err
	}

	if args.DataType == MapRoleData {
		mapRole := NewMapRole(args.RoleARN, args.Username, args.Groups)
		newMap, ok := upsertRole(authData.MapRoles, mapRole)
		if ok {
			logger.Info("upsertAuth.Updated",
				"DataType", args.DataType,
				"Username", args.Username,
				"RoleARN", args.RoleARN)
		} else {
			logger.V(1).Info("upsertAuth.NoNeedUpdate",
				"DataType", args.DataType,
				"Username", args.Username,
				"RoleARN", args.RoleARN)
		}
		authData.SetMapRoles(newMap)
	}

	if args.DataType == MapUserData {
		mapUser := NewMapUser(args.UserARN, args.Username, args.Groups)
		newMap, ok := upsertUser(authData.MapUsers, mapUser)
		if ok {
			logger.Info("upsertAuth.Updated",
				"DataType", args.DataType,
				"Username", args.Username,
				"UserARN", args.UserARN)
		} else {
			logger.V(1).Info("upsertAuth.NoNeedUpdate",
				"DataType", args.DataType,
				"Username", args.Username,
				"UserARN", args.UserARN)
		}
		authData.SetMapUsers(newMap)
	}

	if args.DataType == MapAccountData {
		mapAccount := NewMapAccount(args.AccountID)
		newMap, ok := upsertAccount(authData.MapAccounts, mapAccount)
		if ok {
			logger.Info("upsertAuth.Updated",
				"DataType", args.DataType,
				"AccountID", args.AccountID)
		} else {
			logger.V(1).Info("upsertAuth.NoNeedUpdate",
				"DataType", args.DataType,
				"AccountID", args.AccountID)
		}
		authData.SetMapAccounts(newMap)
	}

	return UpdateAuthMap(m.KubernetesClient, authData, configMap)
}

func upsertRole(authMaps []*MapRole, resource *MapRole) ([]*MapRole, bool) {
	var found, updated bool
	for _, existing := range authMaps {
		// Update existing role in auth map.
		if existing.Username == resource.Username {
			found = true
			if !reflect.DeepEqual(existing.Groups, resource.Groups) {
				existing.SetGroups(resource.Groups)
				updated = true
			}
			if existing.RoleARN != resource.RoleARN {
				existing.SetRoleARN(resource.RoleARN)
				updated = true
			}
		}
	}

	// Insert new role in auth map.
	if !found {
		updated = true
		authMaps = append(authMaps, resource)
	}
	return authMaps, updated
}

func upsertUser(authMaps []*MapUser, resource *MapUser) ([]*MapUser, bool) {
	var found, updated bool
	for _, existing := range authMaps {
		if existing.Username+existing.UserARN == resource.Username+resource.UserARN {
			found = true
			if !reflect.DeepEqual(existing.Groups, resource.Groups) {
				existing.SetGroups(resource.Groups)
				updated = true
			}
		}
	}

	// Insert new user in auth map.
	if !found {
		updated = true
		authMaps = append(authMaps, resource)
	}
	return authMaps, updated
}

func upsertAccount(authMaps []*MapAccount, resource *MapAccount) ([]*MapAccount, bool) {
	var found, updated bool
	for _, existing := range authMaps {

		if existing.AccountID == resource.AccountID {
			found = true
			if existing.AccountID != resource.AccountID {
				existing.SetAccount(resource.AccountID)
				updated = true
			}
		}
	}

	// Insert new account in auth map.
	if !found {
		updated = true
		authMaps = append(authMaps, resource)
	}
	return authMaps, updated
}
