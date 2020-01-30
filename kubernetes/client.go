package kubernetes

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/satori/go.uuid"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	rbac "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	client kubernetes.Interface
}

func NewKubeClient() (*Client, error) {
	config, err := readKubeConfig()
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{client: client}, nil
}

func (c *Client) CreateNamespace(prefix string) (*corev1.Namespace, error) {
	name := prefix + "-" + uuid.NewV4().String()
	labels := make(map[string]string)
	labels["test"] = prefix
	namespaceObject := corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: name, Labels: labels}}
	return c.client.CoreV1().Namespaces().Create(&namespaceObject)
}

func (c *Client) DeleteNamespace(namespace string) error {
	var gracePeriod int64
	gracePeriod = 0
	return c.client.CoreV1().Namespaces().Delete(namespace, &metav1.DeleteOptions{GracePeriodSeconds: &gracePeriod})
}

func (c *Client) CreateDeployment(namespace string, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	return c.client.AppsV1().Deployments(namespace).Create(deployment)
}

func (c *Client) GetDeployment(namespace, deploymentName string) (*appsv1.Deployment, error) {
	return c.client.AppsV1().Deployments(namespace).Get(deploymentName, metav1.GetOptions{})
}

func (c *Client) GetDeployments(namespace, selector string) (*appsv1.DeploymentList, error) {
	return c.client.AppsV1().Deployments(namespace).List(metav1.ListOptions{LabelSelector: selector})
}

func (c *Client) DeleteDeployment(namespace, deploymentName string) error {
	return c.client.AppsV1().Deployments(namespace).Delete(deploymentName, &metav1.DeleteOptions{})
}

func (c *Client) WaitForDeployment(namespace, deploymentName string, timeout time.Duration, writer io.Writer) error {
	w, err := c.client.AppsV1().Deployments(namespace).Watch(metav1.ListOptions{})
	if err != nil {
		return err
	}

	_, err = watch.Until(timeout, w, func(event watch.Event) (bool, error) {
		deployment, ok := event.Object.(*appsv1.Deployment)
		if !ok {
			return false, fmt.Errorf("expected `%#v` to be of type appsv1.Deployment", event.Object)
		}

		if deployment.Name == deploymentName {
			if deployment.Status.AvailableReplicas == deployment.Status.UpdatedReplicas {
				return true, nil
			}
			fmt.Fprintf(writer, "Waiting for %d to be equal to %d\n", deployment.Status.AvailableReplicas, deployment.Status.UpdatedReplicas)
		}

		return false, nil
	})

	if err != nil {
		return fmt.Errorf("deployment `%s` did not finish rolling out with error: %s", deploymentName, err)
	}

	return nil
}

func (c *Client) CreateServiceAccount(namespace string, serviceAccount *corev1.ServiceAccount) (*corev1.ServiceAccount, error) {
	return c.client.CoreV1().ServiceAccounts(namespace).Create(serviceAccount)
}

func (c *Client) DeleteServiceAccount(namespace string, serviceAccountName string) error {
	return c.client.CoreV1().ServiceAccounts(namespace).Delete(serviceAccountName, &metav1.DeleteOptions{})
}

func (c *Client) CreatePodSecurityPolicy(podSecurityPolicy *policyv1.PodSecurityPolicy) (*policyv1.PodSecurityPolicy, error) {
	return c.client.PolicyV1beta1().PodSecurityPolicies().Create(podSecurityPolicy)
}

func (c *Client) DeletePodSecurityPolicy(podSecurityPolicyName string) error {
	return c.client.PolicyV1beta1().PodSecurityPolicies().Delete(podSecurityPolicyName, &metav1.DeleteOptions{})
}

func (c *Client) CreateRole(namespace string, role *rbac.Role) (*rbac.Role, error) {
	return c.client.RbacV1().Roles(namespace).Create(role)
}

func (c *Client) DeleteRole(namespace, roleName string) error {
	return c.client.RbacV1().Roles(namespace).Delete(roleName, &metav1.DeleteOptions{})
}

func (c *Client) CreateRoleBinding(namespace string, roleBinding *rbac.RoleBinding) (*rbac.RoleBinding, error) {
	return c.client.RbacV1().RoleBindings(namespace).Create(roleBinding)
}

func (c *Client) DeleteRoleBinding(namespace, roleBindingName string) error {
	return c.client.RbacV1().RoleBindings(namespace).Delete(roleBindingName, &metav1.DeleteOptions{})
}

func (c *Client) IsHealthy() bool {
	var status int
	c.client.CoreV1().RESTClient().Get().RequestURI("/healthz").Do().StatusCode(&status)
	if status == http.StatusOK {
		return true
	}
	return false
}

func readKubeConfig() (*restclient.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	).ClientConfig()
}
