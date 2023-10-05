package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

	app.Get("/blogs", func(c *fiber.Ctx) error {
		return c.JSON(blogs)
	})

	app.Get("/blogs/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(http.StatusNotFound).JSON(err)
		}
		blog, err := findBlogById(id)
		if err != nil {
			return c.Status(http.StatusNotFound).JSON(err)
		}
		return c.JSON(blog)
	})

	app.Delete("/blogs/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(http.StatusNotFound).JSON(err)
		}
		err = deleteBlogById(id)
		if err != nil {
			return c.Status(http.StatusNotFound).JSON(err)
		}
		return c.JSON(nil)
	})

	app.Post("/blogs", func(c *fiber.Ctx) error {
		blog := new(Blog)
		if err := c.BodyParser(blog); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(err)
		}
		addBlog(blog)
		return c.JSON(nil)
	})

	err := app.Listen(":8080")
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

// Blog structs

type Blog struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	Author string `json:"author"`
	Id     int    `json:"id"`
}

type Blogs []Blog

var blogs Blogs = []Blog{
	{
		Title:  "My First Blog",
		Body:   "Why do we use it?\nIt is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout. The point of using Lorem Ipsum is that it has a more-or-less normal distribution of letters, as opposed to using 'Content here, content here', making it look like readable English. Many desktop publishing packages and web page editors now use Lorem Ipsum as their default model text, and a search for 'lorem ipsum' will uncover many web sites still in their infancy. Various versions have evolved over the years, sometimes by accident, sometimes on purpose (injected humour and the like).\n\n\nWhere does it come from?\nContrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. Richard McClintock, a Latin professor at Hampden-Sydney College in Virginia, looked up one of the more obscure Latin words, consectetur, from a Lorem Ipsum passage, and going through the cites of the word in classical literature, discovered the undoubtable source. Lorem Ipsum comes from sections 1.10.32 and 1.10.33 of \"de Finibus Bonorum et Malorum\" (The Extremes of Good and Evil) by Cicero, written in 45 BC. This book is a treatise on the theory of ethics, very popular during the Renaissance. The first line of Lorem Ipsum, \"Lorem ipsum dolor sit amet..\", comes from a line in section 1.10.32.\n\nThe standard chunk of Lorem Ipsum used since the 1500s is reproduced below for those interested. Sections 1.10.32 and 1.10.33 from \"de Finibus Bonorum et Malorum\" by Cicero are also reproduced in their exact original form, accompanied by English versions from the 1914 translation by H. Rackham.\n\nWhere can I get some?\nThere are many variations of passages of Lorem Ipsum available, but the majority have suffered alteration in some form, by injected humour, or randomised words which don't look even slightly believable. If you are going to use a passage of Lorem Ipsum, you need to be sure there isn't anything embarrassing hidden in the middle of text. All the Lorem Ipsum generators on the Internet tend to repeat predefined chunks as necessary, making this the first true generator on the Internet. It uses a dictionary of over 200 Latin words, combined with a handful of model sentence structures, to generate Lorem Ipsum which looks reasonable. The generated Lorem Ipsum is therefore always free from repetition, injected humour, or non-characteristic words etc.",
		Author: "mario",
		Id:     1,
	},
	{
		Title:  "My Second Blog",
		Body:   "Why do we use it?\nIt is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout. The point of using Lorem Ipsum is that it has a more-or-less normal distribution of letters, as opposed to using 'Content here, content here', making it look like readable English. Many desktop publishing packages and web page editors now use Lorem Ipsum as their default model text, and a search for 'lorem ipsum' will uncover many web sites still in their infancy. Various versions have evolved over the years, sometimes by accident, sometimes on purpose (injected humour and the like).\n\n\nWhere does it come from?\nContrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. Richard McClintock, a Latin professor at Hampden-Sydney College in Virginia, looked up one of the more obscure Latin words, consectetur, from a Lorem Ipsum passage, and going through the cites of the word in classical literature, discovered the undoubtable source. Lorem Ipsum comes from sections 1.10.32 and 1.10.33 of \"de Finibus Bonorum et Malorum\" (The Extremes of Good and Evil) by Cicero, written in 45 BC. This book is a treatise on the theory of ethics, very popular during the Renaissance. The first line of Lorem Ipsum, \"Lorem ipsum dolor sit amet..\", comes from a line in section 1.10.32.\n\nThe standard chunk of Lorem Ipsum used since the 1500s is reproduced below for those interested. Sections 1.10.32 and 1.10.33 from \"de Finibus Bonorum et Malorum\" by Cicero are also reproduced in their exact original form, accompanied by English versions from the 1914 translation by H. Rackham.\n\nWhere can I get some?\nThere are many variations of passages of Lorem Ipsum available, but the majority have suffered alteration in some form, by injected humour, or randomised words which don't look even slightly believable. If you are going to use a passage of Lorem Ipsum, you need to be sure there isn't anything embarrassing hidden in the middle of text. All the Lorem Ipsum generators on the Internet tend to repeat predefined chunks as necessary, making this the first true generator on the Internet. It uses a dictionary of over 200 Latin words, combined with a handful of model sentence structures, to generate Lorem Ipsum which looks reasonable. The generated Lorem Ipsum is therefore always free from repetition, injected humour, or non-characteristic words etc.",
		Author: "yoshi",
		Id:     2,
	},
}

func findBlogById(id int) (*Blog, error) {
	for _, blog := range blogs {
		if blog.Id == id {
			return &blog, nil
		}
	}
	return nil, fmt.Errorf("Blog not found for id: %d", id)
}

func deleteBlogById(id int) error {
	for i, blog := range blogs {
		if blog.Id == id {
			blogs = append(blogs[:i], blogs[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Blog not found for id: %d", id)
}

func addBlog(blog *Blog) error {
	blog.Id = len(blogs) + 1
	blogs = append(blogs, *blog)
	return nil
}
