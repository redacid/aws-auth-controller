# aws-auth-controller

kubebuilder init --domain prozorro.sale --plugins go/v4 --repo github.com/redacid/aws-auth-controller --project-name=aws-auth-controller-prozorro
kubebuilder create api --group aws-auth --version v1beta1 --kind MapRole --namespaced true --controller --resource
kubebuilder create api --group aws-auth --version v1beta1 --kind MapUser --namespaced true --controller --resource
kubebuilder create api --group aws-auth --version v1beta1 --kind MapAccount --namespaced true --controller --resource

kubebuilder create webhook --group aws-auth --version v1beta1 --kind MapUser --defaulting --programmatic-validation
kubebuilder create webhook --group aws-auth --version v1beta1 --kind MapRole --defaulting --programmatic-validation
kubebuilder create webhook --group aws-auth --version v1beta1 --kind MapAccount --defaulting --programmatic-validation

git tag 0.1.0 -m "Test Release 0.1.0"
git push origin 0.1.0

cd ./charts helm upgrade -i --namespace kube-system aws-auth aws-auth-operator
cd ./config/samples kubectl apply -f aws-auth_v1beta1_mapuser.yaml -n kube-system
cd ./config/samples kubectl apply -f aws-auth_v1beta1_maprole.yaml -n kube-system
cd ./config/samples kubectl apply -f aws-auth_v1beta1_mapaccount.yaml -n kube-system

https://book.kubebuilder.io/reference/using-finalizers
https://book.kubebuilder.io/reference/markers/crd-validation

https://github.com/kubernetes-sigs/controller-runtime/blob/main/pkg/reconcile/reconcile.go#L46
https://book.kubebuilder.io/reference/watching-resources.html?highlight=RequeueAfter#when-requeueafter-x-is-useful

https://github.com/gp42/aws-auth-operator

Usage of ./manager:
-config-map-name string
ConfigMap-name where store items, --config-map-name=aws-auth, if not set, use default name aws-auth (default "aws-auth")
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
Time to reconcile CRDs and recreate  configmap item if need, --reconcile-time=30m, s - seconds, m - minutes, h - hours, d - days (default 30m0s)
-username-must-be-email
Enables checking of username, should it be an email address, default false.Use --username-must-be-email=true to enable this feature.
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