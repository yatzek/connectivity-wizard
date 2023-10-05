# Conectivity Wizard

## Try it locally

```
make run
```
Visit http://localhost:8080


## Try it from local k8s

Install kind and skaffold:

```
brew install kind skaffold
```

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

Port-forward the application to localhost:8080:

```
kubectl port-forward service/web 8080:8080
```

Test the application:

```
curl localhost:8080
curl localhost:8080/pods
curl localhost:8080/deployment
```

## React frontend notes

Create react app:

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

React tutorial: https://www.youtube.com/watch?v=j942wKiXFu8