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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	awsauthv1beta1 "github.com/redacid/aws-auth-controller/api/v1beta1"
)

var _ = Describe("MapAccount Webhook", func() {
	var (
		obj       *awsauthv1beta1.MapAccount
		oldObj    *awsauthv1beta1.MapAccount
		validator MapAccountCustomValidator
		defaulter MapAccountCustomDefaulter
	)
	const (
		ResourceName        = "test-account-resource"
		PresentResourceName = "test-account-resource-present"
		PresentAccountID    = "333333333333"
		AccountID           = "111111111111"
		OldAccountID        = "222222222222"
		Namespace           = "default"
	)

	BeforeEach(func() {
		obj = &awsauthv1beta1.MapAccount{}
		oldObj = &awsauthv1beta1.MapAccount{}
		validator = MapAccountCustomValidator{}
		Expect(validator).NotTo(BeNil(), "Expected validator to be initialized")
		defaulter = MapAccountCustomDefaulter{}
		Expect(defaulter).NotTo(BeNil(), "Expected defaulter to be initialized")
		Expect(oldObj).NotTo(BeNil(), "Expected oldObj to be initialized")
		Expect(obj).NotTo(BeNil(), "Expected obj to be initialized")
	})

	AfterEach(func() {
		// TODO (user): Add any teardown logic common to all tests
	})

	Context("When creating MapAccount under Defaulting Webhook", func() {
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

	Context("When creating or updating MapAccount under Validating Webhook", func() {

		It("Should MapAccount created", func() {
			By("By creating a new MapAccount")
			ctx := context.Background()
			mapAccount := &awsauthv1beta1.MapAccount{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "aws-auth.prozorro.sale/v1beta1",
					Kind:       "MapAccount",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      PresentResourceName,
					Namespace: Namespace,
				},
				Spec: awsauthv1beta1.MapAccountSpec{
					AccountID: PresentAccountID,
				},
			}
			Expect(k8sClient.Create(ctx, mapAccount)).To(Succeed())
		})

		It("Should MapAccount create new account with same data", func() {
			By("By invalid creating a new MapAccount same account id")
			ctx := context.Background()
			mapAccount := &awsauthv1beta1.MapAccount{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "aws-auth.prozorro.sale/v1beta1",
					Kind:       "MapAccount",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      ResourceName,
					Namespace: Namespace,
				},
				Spec: awsauthv1beta1.MapAccountSpec{
					AccountID: PresentAccountID,
				},
			}
			Expect(k8sClient.Create(ctx, mapAccount)).Error().To(HaveOccurred())
		})

		It("Should validate incorrect updates", func() {
			By("simulating a invalid update scenario")
			oldObj.Spec.AccountID = OldAccountID
			oldObj.Namespace = Namespace
			obj.Name = ResourceName
			obj.Spec.AccountID = AccountID
			obj.Namespace = Namespace
			Expect(validator.ValidateUpdate(ctx, oldObj, obj)).Error().To(HaveOccurred())
		})
	})

})
