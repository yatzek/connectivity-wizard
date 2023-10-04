package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	app := fiber.New()
	app.Static("/", "./frontend/build")

	app.Get("/hello", helloWorld)
	app.Get("/pods", getPodsHandler)
	app.Get("/deployment", createDeploymentHandler)

	err := app.Listen(":8000")
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

type Response struct {
	Message string `json:"message"`
}

// func home(c *fiber.Ctx) error {
// 	return c.Render("index", fiber.Map{})
// }

func helloWorld(c *fiber.Ctx) error {
	return c.JSON(Response{"Hello, WorldXXX!"})
}

func getPodsHandler(c *fiber.Ctx) error {
	p, err := getPods()
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(p)
}

func createDeploymentHandler(c *fiber.Ctx) error {
	d, err := createDeployment()
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(d)
}

// Runs this to grant permissions to default service account:
// 	kubectl create clusterrolebinding default-admin \
// 		--clusterrole=admin  \
// 		--serviceaccount=default:default

// robbed from:
// https://github.com/kubernetes/client-go/blob/master/examples/in-cluster-client-configuration/main.go

func getPods() ([]apiv1.Pod, error) {
	clientset, err := clientset()
	if err != nil {
		return nil, err
	}

	// get pods in the default namespaces
	pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Error getting pods", err)
		return nil, err
	}

	return pods.Items, nil
}

func createDeployment() (*appsv1.Deployment, error) {
	clientset, err := clientset()
	if err != nil {
		return nil, err
	}

	deploymentsClient := clientset.AppsV1().Deployments("customer-x")
	deployment := deploymentDefinition()

	// Create Deployment
	log.Println("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		log.Printf("Error creating deployment", err)
		return nil, err
	}
	log.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
	return result, nil
}

func clientset() (*kubernetes.Clientset, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Printf("Error creating in-cluster config", err)
		return nil, err
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("Error creating clientset", err)
		return nil, err
	}

	return clientset, nil
}

func deploymentDefinition() *appsv1.Deployment {
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: "nginx:1.12",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	return dep
}

func int32Ptr(i int32) *int32 { return &i }
