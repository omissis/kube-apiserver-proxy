package kube

import (
	"fmt"
	"net/http"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"k8s.io/kubectl/pkg/scheme"
)

type K8sRESTClientFactory interface {
	Client(group string, version string) (*rest.RESTClient, error)
	Request(r http.Request) (*rest.Request, error)
}

func NewDefaultK8sRESTClientFactory(
	restConfigFactory RESTConfigFactory,
	httpClient *http.Client,
	kubeconfigPath string,
) *DefaultK8sRESTClientFactory {
	return &DefaultK8sRESTClientFactory{
		restConfigFactory: restConfigFactory,
		httpClient:        httpClient,
		kubeconfigPath:    kubeconfigPath,
	}
}

type DefaultK8sRESTClientFactory struct {
	clients           map[string]map[string]*rest.RESTClient
	restConfigFactory RESTConfigFactory
	httpClient        *http.Client
	kubeconfigPath    string
}

func (k *DefaultK8sRESTClientFactory) Client(group string, version string) (*rest.RESTClient, error) {
	if k.clients == nil {
		k.clients = make(map[string]map[string]*rest.RESTClient)
	}

	_, ok := k.clients[group][version]
	if !ok {
		cfg, err := k.newRESTConfig(group, version)
		if err != nil {
			return nil, fmt.Errorf("cannot create rest config: %w", err)
		}

		clt, err := k.newRESTClient(cfg)
		if err != nil {
			return nil, fmt.Errorf("cannot create rest client: %w", err)
		}

		if k.clients[group] == nil {
			k.clients[group] = make(map[string]*rest.RESTClient)
		}

		k.clients[group][version] = clt
	}

	return k.clients[group][version], nil
}

func (k *DefaultK8sRESTClientFactory) Request(r http.Request) (*rest.Request, error) {
	group, version, err := GetGroupVersionFromURI(r.URL.Path)
	if err != nil {
		return nil, fmt.Errorf("cannot get group and version from request uri: %w", err)
	}

	rc, err := k.Client(group, version)
	if err != nil {
		return nil, fmt.Errorf("cannot create rest client: %w", err)
	}

	return rest.NewRequest(rc).
			Verb(r.Method).
			RequestURI(r.URL.Path).
			Body(r.Body),
		nil
}

func (k *DefaultK8sRESTClientFactory) newRESTClient(config *rest.Config) (*rest.RESTClient, error) {
	if k.httpClient == nil {
		return rest.RESTClientFor(config)
	}

	return rest.RESTClientForConfigAndClient(config, k.httpClient)
}

func (k *DefaultK8sRESTClientFactory) newRESTConfig(group string, version string) (*rest.Config, error) {
	config, err := k.restConfigFactory.New(k.kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("cannot create rest config: %w", err)
	}

	config.GroupVersion = &schema.GroupVersion{
		Group:   group,
		Version: version,
	}

	// config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	return config, nil
}
