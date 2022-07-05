package k8s

import (
	"context"
	"fmt"

	"chaosmonkey.monke/chaos/pkg/logger"
	"chaosmonkey.monke/chaos/pkg/monitoring"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// TODO: Implement pagination for all get requests
// to scale.

var K8sClientInstance = GetClient()

var log = logger.Logger.Sugar()

type Client interface {
	ListFilteredNamespaces(map[string]struct{}) ([]string, error)
	ListDeploymentsWithFilters([]string, map[string]string) ([]v1.Deployment, error)
	ListPodsWithOpts(string, metav1.ListOptions) ([]string, error)
	KillPod(string, string) error
}

type k8sClient struct {
	*kubernetes.Clientset
}

func GetClient() Client {
	log.Info("Generating kubernetes client")
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Errorf("Failed to fetch in-cluster kube config with error: %s", err.Error())
		/*
			log.Info("Fetching kube config from file-system")
			home := homedir.HomeDir()
			kubeConfigFile := os.Getenv("KUBECONFIG")
			if kubeConfigFile == "" {
				log.Warn("No kube config file specified as env variable KUBECONFIG")
				kubeConfigFile = filepath.Join(home, ".kube", "configure")
			}

			log.Infof("Using kube config file: [%s]", kubeConfigFile)
			config, err = clientcmd.BuildConfigFromFlags("", kubeConfigFile)
			if err != nil {
				log.Warnf("Failed to fetch kube config from file-system with error: [%s]", err.Error())
			}
		*/
		return nil
	}

	log.Debug("Get kubernetes client from config")
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error(err.Error())
	}

	return &k8sClient{clientSet}
}

// TODO: pagination when number of pods is too large.
func (c *k8sClient) ListPodsWithOpts(namespace string, opts metav1.ListOptions) ([]string, error) {
	log.Debugf("Listing pods in namespace [%s] with options [%s]", namespace, opts)

	podsSlice := []string{}
	pods, err := c.CoreV1().Pods(namespace).List(context.TODO(), opts)
	if err != nil {
		return nil, err
	}

	log.Debugf("Creating list of namespace/podName")
	for _, pod := range pods.Items {
		podsSlice = append(podsSlice, fmt.Sprintf(`%s/%s`, pod.GetNamespace(), pod.GetName()))
	}

	return podsSlice, nil
}

func (c *k8sClient) ListDeploymentsWithFilters(namespaces []string, annotations map[string]string) ([]v1.Deployment, error) {
	log.Debugf("Listing deployments in namespaces [%s] with annotations [%s]", namespaces, annotations)
	deployments := []v1.Deployment{}

	for _, val := range namespaces {
		log.Debugf("Listing deployments in namespace [%s]", val)
		deploys, err := c.ListDeployments(val, metav1.ListOptions{})
		if err != nil {
			return nil, err
		}

		log.Debugf("Filtering deployments with annotations [%s]", annotations)
		if len(annotations) > 0 {
			for _, deploy := range deploys.Items {
				match := true
				depAnntns := deploy.GetAnnotations()
				for key := range annotations {
					if annotations[key] != depAnntns[key] {
						match = false
						log.Debugf("Deployment [%s] does not have matching annotation [%s: %s]", deploy.GetName(), key, depAnntns[key])
						break
					}
				}
				if match {
					deployments = append(deployments, deploy)
				}
			}
		} else {
			deployments = append(deployments, deploys.Items...)
		}
	}

	return deployments, nil
}

func (c *k8sClient) ListDeployments(namespace string, opts metav1.ListOptions) (*v1.DeploymentList, error) {
	log.Debugf("Listing deployment in namespace [%s] with options [%s]", namespace, opts)

	deploys, err := c.AppsV1().Deployments(namespace).List(context.TODO(), opts)
	if err != nil {
		return nil, err
	}

	return deploys, nil
}

// TODO: More graceful solution would be to get excluded/included namespaces
// based on labels
func (c *k8sClient) ListFilteredNamespaces(excludedNamespaces map[string]struct{}) ([]string, error) {
	log.Debugf("Lisiting namespaces except the excluded namespaces")
	namespaces, err := c.ListNamspaces()
	if err != nil {
		return nil, err
	}

	retval := []string{}

	for _, item := range namespaces.Items {
		if _, ok := excludedNamespaces[item.GetName()]; !ok {
			log.Debugf("Namespace [%s] is not in excluded namespace list", item.GetName())
			retval = append(retval, item.GetName())
		}
	}
	return retval, nil
}

func (c *k8sClient) ListNamspaces() (*corev1.NamespaceList, error) {
	log.Debugf("Listing namespaces")
	namespaces, err := c.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return namespaces, nil
}

func (c *k8sClient) KillPod(name, namespace string) error {
	log.Debugf("Received request to kill pod [%s] in namespace [%s]", name, namespace)

	log.Debugf("Fetching pod [%s] in namespace [%s]", name, namespace)
	pod, err := c.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	log.Debugf("Deleting pid [%s] in namespace [%s]", name, namespace)
	err = c.CoreV1().Pods(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	// Increment counters
	log.Debugf("Incrementing metric counters")
	monitoring.MetricsServerInstance.IncWithLabels("chaosmonkey_pods_by_labels_killed_count", pod.Labels)
	monitoring.MetricsServerInstance.Inc("chaosmonkey_pods_killed_count_total")

	return nil
}
