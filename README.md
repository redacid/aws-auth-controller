# aws-auth-controller
[![codecov](https://codecov.io/github/redacid/aws-auth-controller/graph/badge.svg?token=3B6KI6EJCR)](https://codecov.io/github/redacid/aws-auth-controller)
[![unit-test](https://github.com/redacid/aws-auth-controller/actions/workflows/test.yml/badge.svg?branch=init)](https://github.com/redacid/aws-auth-controller/actions/workflows/test.yml)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/redacid/aws-auth-controller/blob/init/LICENSE)

[Helm chart README.md](charts/aws-auth-controller/README.md)

This repository contains the Golang implementation of a [Kubernetes Operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)
managing the `aws-auth` ConfigMap(for testing you can set another name, see configMapName in values.yaml), built with [Kubebuilder](https://kubebuilder.io/)
For self-signed certs used cert-manager.io, must be present in a cluster. 
Allow creation user-names, role-names with special characters. 
Manipulations with all of crd resources controlled by Validation Webhook and CRDs spec.properties.

```yaml
apiVersion: aws-auth.prozorro.sale/v1beta1
kind: MapUser
metadata:
  name: user1
spec:
  description: A test mapuser1
  username: user1@domain.com # < this name will be stored in cm
  groups:
    - system:masters
    - admins
  userarn: arn:aws:iam::123456789012:user/user1@domain.com
```
```yaml
apiVersion: aws-auth.prozorro.sale/v1beta1
kind: MapRole
metadata:
  name: maprole-managed-node-group
spec:
  rolearn: arn:aws:iam::123456789012:role/managed-node-group
  groups:
    - sample:group1
    - sample:group2
  description: A sample maprole
  username: system:node:{{EC2PrivateDNSName}} # < this name will be stored in cm
```
```yaml
apiVersion: aws-auth.prozorro.sale/v1beta1
kind: MapAccount
metadata:
  name: account-123456789012
spec:
  accountid: "123456789012" # < if this account id is present in mustPresentAccountID, you can't delete this crd resource 
```

## Custom Resource Definitions

This operator provides the following CRD kinds in the `aws-auth.prozorro.sale` API group.

- [MapAccount](charts/examples/mapaccount.yaml)
- [MapRole](charts/examples/maprole.yaml)
- [MapUser](charts/examples/mapuser.yaml)

## External Resources

- [Kubebuilder documentation](https://book.kubebuilder.io/)
- [kubernetes Custom Resources documentation](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)
- [Kubernetes Controllers documentation](https://kubernetes.io/docs/concepts/architecture/controller/)
- [Kubernetes Operator documentation](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)
- [AWS EKS aws-auth ConfigMap documentation](https://docs.aws.amazon.com/eks/latest/userguide/add-user-role.html)

## Related work

- [sambatv aws-auth-operator](https://github.com/sambatv/aws-auth-operator) This repository contains the Golang 
  implementation of a Kubernetes Operator managing the aws-auth ConfigMap, built with Kubebuilder

- [ops42 aws-auth-operator](https://ops42.org/aws-auth-operator/) provides a
  way to map AWS IAM users to the `data.mapUsers` section of the `kube-system:aws-auth`
  ConfigMap.
- [rustrial aws-eks-iam-auth-controller](https://github.com/rustrial/aws-eks-iam-auth-controller)
  similarly provides a way to map AWS IAM users to the `data.mapUsers` section of the
  `kube-system:aws-auth` ConfigMap.
- [aws-auth](https://github.com/keikoproj/aws-auth) provides a CLI and Golang
  package enabling management of the `data.mapRoles` and `data.mapUsers`
  sections of the `kube-system:aws-auth` ConfigMap.