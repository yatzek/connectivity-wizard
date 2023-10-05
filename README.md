# Conectivity Wizard

## Try it locally

```
make run
```
Visit http://localhost:8080


============================ LOCAL K8S ============================


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

============================ FRONTEND ============================

CREATE REACT APP:

```
npx create-react-app frontend --template typescript
cd frontend
# run in development mode
npm run start

# build static assets
npm run build
```

JSON SERVER:
```
cd frontend
npx json-server --watch data/db.json --port 8000
```

Docs:

React tutorial:
	https://www.youtube.com/watch?v=j942wKiXFu8