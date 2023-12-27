kind create cluster --config kind-config.yaml

bash ./scripts/apply-gateways-api.sh
helmfile sync -f helmfiles/helmfile.yaml

kubectl apply -f manifests/metallb.yaml

# cilium connectivity test

kubectl apply -f manifests/argocd-ingress.yaml
