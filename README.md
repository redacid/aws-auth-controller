# aws-auth-controller

kubebuilder init --domain prozorro.sale --plugins go/v4 --repo github.com/redacid/aws-auth-controller --project-name=aws-auth-controller-prozorro
kubebuilder create api --group aws-auth --version v1beta1 --kind MapRole --namespaced true --controller --resource
kubebuilder create api --group aws-auth --version v1beta1 --kind MapUser --namespaced true --controller --resource
kubebuilder create api --group aws-auth --version v1beta1 --kind MapAccount --namespaced true --controller --resource

kubebuilder create webhook --group aws-auth --version v1beta1 --kind MapUser --defaulting --programmatic-validation
kubebuilder create webhook --group aws-auth --version v1beta1 --kind MapRole --defaulting --programmatic-validation
kubebuilder create webhook --group aws-auth --version v1beta1 --kind MapAccount --defaulting --programmatic-validation

git tag 0.0.1 -m "Test 0.0.1"
git push origin 0.0.1

cd ./charts helm upgrade -i --namespace kube-system aws-auth aws-auth-operator
cd ./config/samples kubectl apply -f aws-auth_v1beta1_mapuser.yaml -n kube-system
cd ./config/samples kubectl apply -f aws-auth_v1beta1_maprole.yaml -n kube-system
cd ./config/samples kubectl apply -f aws-auth_v1beta1_mapaccount.yaml -n kube-system

https://book.kubebuilder.io/reference/using-finalizers
https://book.kubebuilder.io/reference/markers/crd-validation

https://github.com/gp42/aws-auth-operator