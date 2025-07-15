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
	"context"
	"fmt"
	"time"

	// "log"
	"sync"

	"gopkg.in/yaml.v2"
	kcorev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"strings"
)

const (
	ConfigMapNamespace = "kube-system"
)

var (
	ConfigMapName       = "aws-auth"
	CrdFinalizerName    = "aws-auth.prozorro.sale/finalizer"
	UsernameMustBeEmail = false
	// TODO automatic add to configmap
	MustPresentAccountID    string = ""
	CrdItemAllowedNamespace string = ""

	ReconcileTime  time.Duration = time.Minute * 30
	configMapMutex sync.RWMutex
)

type ControllerArgs struct {
	UsernameMustBeEmail     bool
	MustPresentAccountID    string
	CrdItemAllowedNamespace string
	ConfigMapName           string
	ReconcileTime           time.Duration
}

func DeclareVariables(controllerArgs ControllerArgs) {
	UsernameMustBeEmail = controllerArgs.UsernameMustBeEmail
	MustPresentAccountID = controllerArgs.MustPresentAccountID
	CrdItemAllowedNamespace = controllerArgs.CrdItemAllowedNamespace
	ConfigMapName = controllerArgs.ConfigMapName
	ReconcileTime = controllerArgs.ReconcileTime
}

// ReadAuthMap reads the auth ConfigMap and returns AwsAuthData and the read ConfigMap.
func ReadAuthMap(k kubernetes.Interface) (AwsAuthData, *kcorev1.ConfigMap, error) {
	configMapMutex.RLock()
	defer configMapMutex.RUnlock()

	var authData AwsAuthData

	cm, err := k.CoreV1().ConfigMaps(ConfigMapNamespace).Get(context.Background(), ConfigMapName, apismetav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			cm, err = CreateAuthMap(k)
			if err != nil {
				return authData, cm, err
			}
		} else {
			return authData, cm, err
		}
	}

	err = yaml.Unmarshal([]byte(cm.Data["mapRoles"]), &authData.MapRoles)
	if err != nil {
		return authData, cm, err
	}

	/* err = yaml.Unmarshal([]byte(cm.Data["mapAccounts"]), &authData.MapAccounts)
	 if err != nil {
	   	return authData, cm, err
	} */

	var accountIDs []string
	err = yaml.Unmarshal([]byte(cm.Data["mapAccounts"]), &accountIDs)
	if err != nil {
		return authData, cm, err
	}
	authData.MapAccounts = ConvertStringToMapAccount(accountIDs)

	err = yaml.Unmarshal([]byte(cm.Data["mapUsers"]), &authData.MapUsers)
	return authData, cm, err

}

func CreateAuthMap(k kubernetes.Interface) (*kcorev1.ConfigMap, error) {
	configMapObject := &kcorev1.ConfigMap{
		ObjectMeta: apismetav1.ObjectMeta{
			Name:      ConfigMapName,
			Namespace: ConfigMapNamespace,
		},
	}
	return k.CoreV1().ConfigMaps(ConfigMapNamespace).Create(
		context.Background(),
		configMapObject,
		apismetav1.CreateOptions{},
	)
}

// UpdateAuthMap updates a given ConfigMap
func UpdateAuthMap(k kubernetes.Interface, authData AwsAuthData, cm *kcorev1.ConfigMap) error {
	configMapMutex.Lock()
	defer configMapMutex.Unlock()

	mapRoles, err := yaml.Marshal(authData.MapRoles)
	if err != nil {
		return err
	}

	mapUsers, err := yaml.Marshal(authData.MapUsers)
	if err != nil {
		return err
	}

	mapAccounts, err := yaml.Marshal(ConvertMapAccountToString(authData.MapAccounts))
	if err != nil {
		return err
	}

	cm.Data = map[string]string{
		"mapRoles":    string(mapRoles),
		"mapUsers":    string(mapUsers),
		"mapAccounts": string(mapAccounts),
	}

	_, err = k.CoreV1().ConfigMaps(ConfigMapNamespace).Update(context.Background(), cm, apismetav1.UpdateOptions{})
	return err
}

// AwsAuthData represents the data of the aws-auth configmap
type AwsAuthData struct {
	MapRoles    []*MapRole    `yaml:"mapRoles"`
	MapUsers    []*MapUser    `yaml:"mapUsers"`
	MapAccounts []*MapAccount `yaml:"mapAccounts"`
}

// SetMapRoles sets the MapRoles element
func (m *AwsAuthData) SetMapRoles(authMap []*MapRole) {
	m.MapRoles = authMap
}

// SetMapUsers sets the MapUsers element
func (m *AwsAuthData) SetMapUsers(authMap []*MapUser) {
	m.MapUsers = authMap
}

func (m *AwsAuthData) SetMapAccounts(authMap []*MapAccount) {
	m.MapAccounts = authMap
}

// MapRole is the basic structure of a mapRoles authentication object
type MapRole struct {
	RoleARN  string   `yaml:"rolearn"`
	Username string   `yaml:"username"`
	Groups   []string `yaml:"groups,omitempty"`
}

func (r *MapRole) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("- rolearn: %v\n  ", r.RoleARN))
	s.WriteString(fmt.Sprintf("username: %v\n  ", r.Username))
	s.WriteString("groups:\n")
	for _, group := range r.Groups {
		s.WriteString(fmt.Sprintf("  - %v\n", group))
	}
	return s.String()

}

// SetGroups sets the Groups value
func (r *MapRole) SetGroups(g []string) *MapRole {
	r.Groups = g
	return r
}

// SetRoleARN sets the Username value
func (r *MapRole) SetRoleARN(v string) *MapRole {
	r.Username = v
	return r
}

// NewMapRole returns a new NewMapRole
func NewMapRole(rolearn, username string, groups []string) *MapRole {
	return &MapRole{
		RoleARN:  rolearn,
		Username: username,
		Groups:   groups,
	}
}

// MapUser is the basic structure of a mapUsers authentication object

type MapUser struct {
	UserARN  string   `yaml:"userarn"`
	Username string   `yaml:"username"`
	Groups   []string `yaml:"groups,omitempty"`
}

func (r *MapUser) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("- userarn: %v\n  ", r.UserARN))
	s.WriteString(fmt.Sprintf("username: %v\n  ", r.Username))
	s.WriteString("groups:\n")
	for _, group := range r.Groups {
		s.WriteString(fmt.Sprintf("  - %v\n", group))
	}
	return s.String()
}

// SetGroups sets the Groups value
func (r *MapUser) SetGroups(g []string) *MapUser {
	r.Groups = g
	return r
}

// SetUserARN sets the UserARN value
func (r *MapUser) SetUserARN(v string) *MapUser {
	r.UserARN = v
	return r
}

func (r *MapUser) SetUsername(v string) *MapUser {
	r.Username = v
	return r
}

// NewMapUser returns a new NewMapUser
func NewMapUser(userarn, username string, groups []string) *MapUser {
	return &MapUser{
		UserARN:  userarn,
		Username: username,
		Groups:   groups,
	}
}

type MapAccount struct {
	AccountID string `yaml:"accountid"`
}

func (r *MapAccount) SetAccount(v string) *MapAccount {
	r.AccountID = v
	return r
}

func (r *MapAccount) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("- accountid: %v\n  ", r.AccountID))
	return s.String()
}

func (r *MapAccount) String2() string {
	return fmt.Sprintf("- %s", r.AccountID)
}

func NewMapAccount(accountid string) *MapAccount {
	return &MapAccount{
		AccountID: accountid,
	}
}
