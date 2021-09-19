helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx

helm show values ingress-nginx/ingress-nginx

https://artifacthub.io/packages/helm/ingress-nginx/ingress-nginx
source: https://github.com/kubernetes/ingress-nginx

helm install -n ingress-nginx --version 4.0.0-beta.3 -f ./custom-values.yaml  ingress-nginx ingress-nginx/ingress-nginx

ingress nginx helm chart released
helm install -n ingress-nginx -f ./custom-values.yaml  ingress-nginx ingress-nginx/ingress-nginx
