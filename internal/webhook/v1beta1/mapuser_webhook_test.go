/*
Copyright 2025.

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

package v1beta1

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	awsauthv1beta1 "github.com/redacid/aws-auth-controller/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("MapUser Webhook", func() {
	var (
		obj       *awsauthv1beta1.MapUser
		oldObj    *awsauthv1beta1.MapUser
		validator MapUserCustomValidator
		defaulter MapUserCustomDefaulter
	)

	const ResourceName = "test-user-resource"
	const PresentResourceName = "test-user-resource-present"
	const PresentUsername = "test-user-present"
	const PresentUserARN = "arn:aws:iam::123456789012:user/test-user-present"
	var PresentGroups = []string{"system:users", "system:viewers"}
	const UserName = "test-user"
	const OldUserName = "test-user-old"
	const OldUserARN = "arn:aws:iam::123456789012:user/test-user-old"
	var OldGroups = []string{"system:users", "system:viewers"}
	var Groups = []string{"system:masters", "system:nodes"}
	const UserARN = "arn:aws:iam::123456789012:user/test-user"
	const Namespace = "default"

	BeforeEach(func() {
		obj = &awsauthv1beta1.MapUser{}
		oldObj = &awsauthv1beta1.MapUser{}
		validator = MapUserCustomValidator{}
		Expect(validator).NotTo(BeNil(), "Expected validator to be initialized")
		defaulter = MapUserCustomDefaulter{}
		Expect(defaulter).NotTo(BeNil(), "Expected defaulter to be initialized")
		Expect(oldObj).NotTo(BeNil(), "Expected oldObj to be initialized")
		Expect(obj).NotTo(BeNil(), "Expected obj to be initialized")
	})

	AfterEach(func() {
		// TODO (user): Add any teardown logic common to all tests
	})

	Context("When creating MapUser under Defaulting Webhook", func() {
		// TODO (user): Add logic for defaulting webhooks
		// Example:
		// It("Should apply defaults when a required field is empty", func() {
		//     By("simulating a scenario where defaults should be applied")
		//     obj.SomeFieldWithDefault = ""
		//     By("calling the Default method to apply defaults")
		//     defaulter.Default(ctx, obj)
		//     By("checking that the default values are set")
		//     Expect(obj.SomeFieldWithDefault).To(Equal("default_value"))
		// })
	})

	Context("When creating or updating MapUser under Validating Webhook", func() {

		It("Should MapUser created", func() {
			By("By creating a new MapUser")
			ctx := context.Background()
			mapUser := &awsauthv1beta1.MapUser{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "aws-auth.prozorro.sale/v1beta1",
					Kind:       "MapUser",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      PresentResourceName,
					Namespace: Namespace,
				},
				Spec: awsauthv1beta1.MapUserSpec{
					Username: PresentUsername,
					Groups:   PresentGroups,
					UserARN:  PresentUserARN,
				},
			}
			Expect(k8sClient.Create(ctx, mapUser)).To(Succeed())
		})
		It("Should MapUser create new user with same Username", func() {
			By("By invalid creating a new MapUser same Username")
			ctx := context.Background()
			mapUser := &awsauthv1beta1.MapUser{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "aws-auth.prozorro.sale/v1beta1",
					Kind:       "MapUser",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      ResourceName,
					Namespace: Namespace,
				},
				Spec: awsauthv1beta1.MapUserSpec{
					Username: PresentUsername,
					Groups:   PresentGroups,
					UserARN:  UserARN,
				},
			}
			Expect(k8sClient.Create(ctx, mapUser)).Error().To(HaveOccurred())
		})
		It("Should MapUser create new user with same UserARN", func() {
			By("By invalid creating a new MapUser same UserARN")
			ctx := context.Background()
			mapUser := &awsauthv1beta1.MapUser{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "aws-auth.prozorro.sale/v1beta1",
					Kind:       "MapUser",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      ResourceName,
					Namespace: Namespace,
				},
				Spec: awsauthv1beta1.MapUserSpec{
					Username: UserName,
					Groups:   PresentGroups,
					UserARN:  PresentUserARN,
				},
			}
			Expect(k8sClient.Create(ctx, mapUser)).Error().To(HaveOccurred())
		})
		It("Should validate incorrect updates", func() {
			By("simulating a invalid update scenario")
			oldObj.Spec.Username = OldUserName
			oldObj.Spec.Groups = OldGroups
			oldObj.Spec.UserARN = OldUserARN
			oldObj.Namespace = Namespace
			obj.Name = ResourceName
			obj.Spec.Groups = Groups
			obj.Spec.UserARN = UserARN
			obj.Spec.Username = UserName
			obj.Namespace = Namespace
			Expect(validator.ValidateUpdate(ctx, oldObj, obj)).Error().To(HaveOccurred())
		})
		It("Should validate updates correctly", func() {
			By("simulating a valid update scenario")
			oldObj.Spec.Username = UserName
			oldObj.Spec.Groups = OldGroups
			oldObj.Spec.UserARN = UserARN
			oldObj.Namespace = Namespace
			obj.Name = ResourceName
			obj.Spec.Groups = Groups
			obj.Spec.UserARN = UserARN
			obj.Spec.Username = UserName
			obj.Namespace = Namespace
			Expect(validator.ValidateUpdate(ctx, oldObj, obj)).To(BeNil())
		})
	})
})
