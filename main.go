package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	kubernetesClient "chatops/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/gin-gonic/gin"
)

func main() {

	// Initialize Kubernetes client
	kubeConfig, err := kubernetesClient.GetKubeConfig()
	if err != nil {
		log.Fatalf("Failed to get Kubernetes config: %v", err)
	}

	// Initialize Gin router
	router := gin.Default()

	// Define API routes
	api := router.Group("/api")
	{
		api.POST("/pod/restart/:namespace/:deploymentName", RestartHandler(kubeConfig))
	}

	// Start the server
	port := 8080
	serverAddr := fmt.Sprintf(":%d", port)
	log.Printf("Starting ChatOps service on port %d...", port)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start CubeFlow service: %v", err)
	}
}

func RestartHandler(clientset *kubernetes.Clientset) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Param("namespace")
		deploymentName := c.Param("deploymentName")

		err := TriggerPodRestart(clientset, namespace, deploymentName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Pod restarted successfully"})
	}
}

func TriggerPodRestart(clientset *kubernetes.Clientset, namespace string, deploymentName string) error {
	// Retrieve the current Deployment
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// restart the deployment
	deployment.Spec.Template.ObjectMeta.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
	_, err = clientset.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		log.Printf("Failed to restart deployment %s in namespace %s: %v", deploymentName, namespace, err)
		return err
	}

	return nil
}
