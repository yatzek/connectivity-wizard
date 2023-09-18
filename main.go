package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	app := fiber.New()

	app.Get("/", helloWorld)
	app.Get("/pods", podHandler)

	err := app.Listen(":8000")
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

type Response struct {
	Message   string `json:"message"`
}

func helloWorld(c *fiber.Ctx) error {
	return c.JSON(Response{"Hello, WorldXXX!"})
}

func podHandler(c *fiber.Ctx) error {
	p, err := pods()
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(p)
}

// Runs this to grant permissions to default service account:
// 	kubectl create clusterrolebinding default-admin \
// 		--clusterrole=admin  \
// 		--serviceaccount=default:default

func pods() ([]v1.Pod, error) {
		// robbed from:
		// https://github.com/kubernetes/client-go/blob/master/examples/in-cluster-client-configuration/main.go

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

		// get pods in the default namespaces
		pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Printf("Error getting pods", err)
			return nil, err
		}

		return pods.Items, nil
}