package conf

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"strings"
	"sync"
)

var (
	opConf   *BaseOperatorConf
	initConf sync.Once
)

const prefixVar = "VM"

//genvars:true
type BaseOperatorConf struct {
	VMAlertDefault struct {
		Image    string `default:"victoriametrics/vmalert"`
		Version  string `default:"v1.37.0"`
		Port     string `default:"8080"`
		Resource struct {
			Limit struct {
				Mem string `default:"500Mi"`
				Cpu string `default:"200m"`
			}
			Request struct {
				Mem string `default:"200Mi"`
				Cpu string `default:"50m"`
			}
		}
		ConfigReloaderCPU    string `default:"100m"`
		ConfigReloaderMemory string `default:"25Mi"`
		ConfigReloadImage    string `default:"jimmidyson/configmap-reload:v0.3.0"`
	}
	VMAgentDefault struct {
		Image             string `default:"victoriametrics/vmagent"`
		Version           string `default:"v1.37.0"`
		ConfigReloadImage string `default:"quay.io/coreos/prometheus-config-reloader:v0.30.1"`
		Port              string `default:"8429"`
		Resource          struct {
			Limit struct {
				Mem string `default:"500Mi"`
				Cpu string `default:"200m"`
			}
			Request struct {
				Mem string `default:"200Mi"`
				Cpu string `default:"50m"`
			}
		}
		ConfigReloaderCPU    string `default:"100m"`
		ConfigReloaderMemory string `default:"25Mi"`
	}

	VMSingleDefault struct {
		Image    string `default:"victoriametrics/victoria-metrics"`
		Version  string `default:"v1.37.0"`
		Port     string `default:"8429"`
		Resource struct {
			Limit struct {
				Mem string `default:"1500Mi"`
				Cpu string `default:"1200m"`
			}
			Request struct {
				Mem string `default:"500Mi"`
				Cpu string `default:"150m"`
			}
		}
		ConfigReloaderCPU    string `default:"100m"`
		ConfigReloaderMemory string `default:"25Mi"`
	}
	VMAlertManager struct {
		ConfigReloaderImage          string `default:"jimmidyson/configmap-reload:v0.3.0"`
		ConfigReloaderCPU            string `default:"100m"`
		ConfigReloaderMemory         string `default:"25Mi"`
		AlertmanagerDefaultBaseImage string `default:"quay.io/prometheus/alertmanager"`
		AlertManagerVersion          string `default:"v0.20.0"`
		LocalHost                    string `default:"127.0.0.1"`
		LogLevel                     string `default:"INFO"`
		LogFormat                    string
		PromSelector                 string
		Namespaces                   Namespaces `ignored:"true"`
		AlertManagerSelector         string
		ClusterDomain                string `default:""`
		KubeletObject                string
	}
	DisabledServiceMonitorCreation bool   `default:"false"`
	Host                           string `default:"0.0.0.0"`
	ListenAddress                  string `default:"0.0.0.0"`
	DefaultLabels                  string `default:"managed-by=vm-operator"`
	Labels                         Labels `ignored:"true"`
	LogLevel                       string
	LogFormat                      string
}

func MustGetBaseConfig() *BaseOperatorConf {
	initConf.Do(func() {
		c := &BaseOperatorConf{}
		err := envconfig.Process(prefixVar, c)
		if err != nil {
			panic(err)
		}
		if c.DefaultLabels != "" {
			defL := Labels{}
			err := defL.Set(c.DefaultLabels)
			if err != nil {
				panic(err)
			}
			c.Labels = defL
		}
		//if c.DefaultNamespaces != "" {
		//}
		opConf = c
		fmt.Printf("conf inited \n")
	})
	return opConf
}

type Labels struct {
	LabelsString string
	LabelsMap    map[string]string
}

// Implement the flag.Value interface
func (labels *Labels) String() string {
	return labels.LabelsString
}

// Merge labels create a new map with labels merged.
func (labels *Labels) Merge(otherLabels map[string]string) map[string]string {
	mergedLabels := map[string]string{}

	for key, value := range otherLabels {
		mergedLabels[key] = value
	}

	for key, value := range labels.LabelsMap {
		mergedLabels[key] = value
	}
	return mergedLabels
}

// Set implements the flag.Set interface.
func (labels *Labels) Set(value string) error {
	m := map[string]string{}
	if value != "" {
		splited := strings.Split(value, ",")
		for _, pair := range splited {
			sp := strings.Split(pair, "=")
			m[sp[0]] = sp[1]
		}
	}
	(*labels).LabelsMap = m
	(*labels).LabelsString = value
	return nil
}

type Namespaces struct {
	// allow list/deny list for common custom resources
	AllowList, DenyList map[string]struct{}
	// allow list for prometheus/alertmanager custom resources

}
