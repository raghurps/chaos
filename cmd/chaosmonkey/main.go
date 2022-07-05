package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"chaosmonkey.monke/chaos/pkg/k8s"
	"chaosmonkey.monke/chaos/pkg/logger"
	"chaosmonkey.monke/chaos/pkg/monitoring"
	"chaosmonkey.monke/chaos/pkg/random"
	"github.com/urfave/cli/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

var log = logger.Logger.Sugar()

func main() {
	defer log.Sync()

	log.Debug("Setting up signal handler")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGABRT, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1)

	log.Debug("Initializing chaosmonkey")
	app := &cli.App{
		Name:    "chaosmonkey",
		Usage:   "create chaos in an otherwise orderly world",
		Version: "v0.0.1",
		Flags: []cli.Flag{
			&cli.DurationFlag{
				Name:        "interval",
				DefaultText: "60s",
				Value:       60 * time.Second,
				EnvVars:     []string{"INTERVAL"},
				Usage:       "Provide interval for killing pods randomly",
			},
			&cli.StringSliceFlag{
				Name:    "namespaces",
				EnvVars: []string{"NAMESPACES"},
				Value:   cli.NewStringSlice(),
				Usage:   "Provide multiple namespaces separated by comma to select a random pod from",
			},
			&cli.StringSliceFlag{
				Name:        "excluded_namespaces",
				DefaultText: "kube-system",
				EnvVars:     []string{"EXCLUDED_NAMESPACES"},
				Value:       cli.NewStringSlice("kube-system"),
				Usage:       "Provide multiple namespaces separated by comma to avoid selecting random pod from",
			},
			&cli.StringSliceFlag{
				Name:        "deployment_annotations",
				DefaultText: "",
				Usage:       "Provide multiple annotations as key=value",
				EnvVars:     []string{"DEPLOYMENT_ANNOTATIONS"},
			},
		},
		Action: func(ctx *cli.Context) error {
			go func() {
				monitoring.Start()
			}()

			return beginChaos(ctx, sig, done)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err.Error())
	}
}

func beginChaos(ctx *cli.Context, sig chan os.Signal, done chan bool) error {
	interval := ctx.Duration("interval")
	ticker := time.NewTicker(interval)
	namespaces := ctx.StringSlice("namespaces")
	excludedNamespaces := ctx.StringSlice("excluded_namespaces")
	annotations := ctx.StringSlice("deployment_annotations")
	deployment_annotations := map[string]string{}

	log.Infof(`Received arguments: {"interval": %s, "namespaces": [%s], "excluded_namespace": [%s], "deployment_annotations": [%s]}`,
		interval,
		namespaces,
		excludedNamespaces,
		annotations,
	)

	for _, anntn := range annotations {
		anntnSlice := strings.Split(anntn, "=")
		deployment_annotations[anntnSlice[0]] = anntnSlice[1]
	}

	log.Debugf("Deployment annotations map: %v", deployment_annotations)

	go func() {
		log.Info("Starting chaos")
		if err := killRandomPod(namespaces, excludedNamespaces, deployment_annotations); err != nil {
			log.Errorf(err.Error())
			done <- true
			return
		}

		for {
			select {
			case s := <-sig:
				log.Infof("Received signal %s", s)
				done <- true
				return
			case <-ticker.C:
				if err := killRandomPod(namespaces, excludedNamespaces, deployment_annotations); err != nil {
					log.Errorf(err.Error())
				}
			}
		}
	}()

	<-done
	return nil
}

func killRandomPod(namespaces, excludedNamespaces []string, annotations map[string]string) error {
	log.Info("Initiating random pod kill")
	client := k8s.K8sClientInstance
	podCollection := []string{}

	// Generate list of namespaces to select the random pod from
	excludedNsMap := map[string]struct{}{}
	ns := []string{}
	var err error

	for _, val := range excludedNamespaces {
		if val == "" {
			err = fmt.Errorf("empty namespace not allowed")
			return err
		}
		excludedNsMap[val] = struct{}{}
	}

	if len(namespaces) > 0 {
		for _, namespace := range namespaces {
			if namespace == "" {
				err = fmt.Errorf("empty namespace not allowed")
				return err
			}
			if _, ok := excludedNsMap[namespace]; !ok {
				ns = append(ns, namespace)
			}
		}
	} else {
		ns, err = client.ListFilteredNamespaces(excludedNsMap)
		if err != nil {
			return err
		}
	}

	log.Debugf("Processing namespaces: [%s]", ns)

	log.Info("Creating pods list")
	if len(annotations) > 0 {
		log.Debugf("Fetching deployments with annotations [%s]", annotations)

		deployments, err := client.ListDeploymentsWithFilters(ns, annotations)
		if err != nil {
			return err
		}

		for _, deployment := range deployments {
			log.Debugf("Fetching pods associated with deployment [%s/%s]", deployment.GetNamespace(), deployment.GetName())

			selectorLabels := labels.Set(deployment.Spec.Selector.MatchLabels).String()
			pods, err := client.ListPodsWithOpts(deployment.GetNamespace(), metav1.ListOptions{
				LabelSelector: selectorLabels,
			})
			if err != nil {
				return err
			}

			podCollection = append(podCollection, pods...)
		}
	} else {
		for _, namespace := range ns {
			log.Debugf("Querying pods in Namespace: %s\n", namespace)
			pods, err := client.ListPodsWithOpts(namespace, metav1.ListOptions{})
			if err != nil {
				return err
			}
			podCollection = append(podCollection, pods...)
		}
	}

	log.Debugf("Collected pods list: [%s]", podCollection)

	if len(podCollection) == 0 {
		log.Info("Nothing to delete as no pods match the provided input flags")
		return nil
	}
	log.Info("Picking a random pod from list")
	rndIndx := random.GetRandomNumber(len(podCollection))

	log.Infof("Killing pod %s\n", podCollection[rndIndx])
	namespaceAndPodName := strings.Split(podCollection[rndIndx], "/")
	// TODO:  handle self termination
	err = client.KillPod(namespaceAndPodName[1], namespaceAndPodName[0])
	if err != nil {
		return err
	}

	log.Infof("Pod [%s] successfully killed", podCollection[rndIndx])

	return nil
}
