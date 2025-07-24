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
- [AWS EKS aws-auth ConfigMap documentation](https://docs.aws.amazon.com/eks/latest/best-practices/identity-and-access-management.html)
- [CRD Validation](https://book.kubebuilder.io/reference/markers/crd-validation)
- [Using Finalizers](https://book.kubebuilder.io/reference/using-finalizers)
- [About Webhooks](https://book.kubebuilder.io/multiversion-tutorial/webhooks)
- [Writing tests](https://book.kubebuilder.io/cronjob-tutorial/writing-tests)
- [RequeueAfter](https://book.kubebuilder.io/reference/watching-resources.html?highlight=RequeueAfter#when-requeueafter-x-is-useful)

## Controller command line parameters
```text
 ./bin/manager --help
Usage of ./bin/manager:
  -config-map-name string
        ConfigMap-name where stored items, --config-map-name=aws-auth-test (default "aws-auth")
  -crd-item-allowed-namespace string
        Namespace where allowed creation of Resources --crd-item-allowed-namespace=kube-system, If is set, only in this namespace allowed creation of Resources, if not set, all namespaces allowed creation of Resources.
  -enable-http2
        If set, HTTP/2 will be enabled for the metrics and webhook servers
  -health-probe-bind-address string
        The address the probe endpoint binds to. (default ":8081")
  -kubeconfig string
        Paths to a kubeconfig. Only required if out-of-cluster.
  -leader-elect
        Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.
  -metrics-bind-address string
        The address the metrics endpoint binds to. Use :8443 for HTTPS or :8080 for HTTP, or leave as 0 to disable the metrics service. (default "0")
  -metrics-cert-key string
        The name of the metrics server key file. (default "tls.key")
  -metrics-cert-name string
        The name of the metrics server certificate file. (default "tls.crt")
  -metrics-cert-path string
        The directory that contains the metrics server certificate.
  -metrics-secure
        If set, the metrics endpoint is served securely via HTTPS. Use --metrics-secure=false to use HTTP instead. (default true)
  -must-present-account-id string
        This account must be allways present in ConfigMap, --must-present-account-id=123456789012
  -reconcile-time duration
        Time to reconcile CRDs and recreate configmap item if need, --reconcile-time=30m, s - seconds, m - minutes, h - hours, d - days (default 30m0s)
  -username-must-be-email
        Enables checking of username in MapUser, should it be an email address, default false.Use --username-must-be-email=true to enable this feature.
  -webhook-cert-key string
        The name of the webhook key file. (default "tls.key")
  -webhook-cert-name string
        The name of the webhook certificate file. (default "tls.crt")
  -webhook-cert-path string
        The directory that contains the webhook certificate.
  -zap-devel
        Development Mode defaults(encoder=consoleEncoder,logLevel=Debug,stackTraceLevel=Warn). Production Mode defaults(encoder=jsonEncoder,logLevel=Info,stackTraceLevel=Error) (default true)
  -zap-encoder value
        Zap log encoding (one of 'json' or 'console')
  -zap-log-level value
        Zap Level to configure the verbosity of logging. Can be one of 'debug', 'info', 'error', 'panic'or any integer value > 0 which corresponds to custom debug levels of increasing verbosity
  -zap-stacktrace-level value
        Zap Level at and above which stacktraces are captured (one of 'info', 'error', 'panic').
  -zap-time-encoding value
        Zap time encoding (one of 'epoch', 'millis', 'nano', 'iso8601', 'rfc3339' or 'rfc3339nano'). Defaults to 'epoch'.
```

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