
Install kind, skaffold and ko:

```
brew install kind skaffold ko
```

Install kubectl-netshoot plugin:
https://github.com/nilic/kubectl-netshoot


Create a kind cluster:

```
kind create cluster --name kind --image kindest/node:v1.23.17
```

Grant admin role to default:default service acount:

```
kubectl create clusterrolebinding default-admin \
 	--clusterrole=admin  \
 	--serviceaccount=default:default
```

Create 'customer-x' namespace:

```
kubectl create namespace customer-x
```

Deploy the application to your local cluster in the default namespace:

```
skaffold run
```

Attach a debug container to the application:

```
 kubectl netshoot debug web-8674d4b855-4gwmh
```

Test the application from inside a debug container:

```
curl web:8000
curl web:8000/pods
curl web:8000/deployment
```

