package env

//App application
type App struct {
	Namespace string     `yaml:"namespace"`
	Tags      []string   `yaml:"tags"`
	Helm      HelmConfig `yaml:"helm"`
	Priority  int        `yaml:"priority"`
	Ingress   struct {
		Enabled bool   `yaml:"enabled"`
		Host    bool   `yaml:"host"`
		Path    string `yaml:"path"`
	}
}

//HelmConfig helm configuration
type HelmConfig struct {
	Chart       string      `yaml:"chart"`
	Repo        string      `yaml:"repo"`
	Version     string      `yaml:"version"`
	Values      interface{} `yaml:"values"`
	ValuesFiles []string    `yaml:"files"`
}

//LocalEnvironment local environment
type LocalEnvironment struct {
	Cluster struct {
		Namespace string `yaml:"namespace"`
	} `yaml:"cluster"`
	Apps map[string]*App `yaml:"apps"`
}

//Namespace namespace of the application
func (e *LocalEnvironment) Namespace(appName string) string {
	app, exists := e.Apps[appName]
	if exists {
		if len(app.Namespace) > 0 {
			return app.Namespace
		}
	}
	return e.Cluster.Namespace
}

func id(namespace, name string) string {
	return namespace + "-" + name
}

//AppAction application action
type AppAction int

const (
	//AppActionNothing application action nothing
	AppActionNothing AppAction = iota
	//AppActionNotFound application action not found
	AppActionNotFound
	//AppActionInstall application action install
	AppActionInstall
	//AppActionUpgrade application action upgrade
	AppActionUpgrade
	//AppActionDowngrade application action downgrade
	AppActionDowngrade
	//AppActionUninstall application action uninstall
	AppActionUninstall
)

//String string representation of the application action
func (a AppAction) String() string {
	switch a {
	case AppActionNothing:
		return ""
	case AppActionNotFound:
		return ""
	case AppActionInstall:
		return "install"
	case AppActionUpgrade:
		return "upgrade"
	case AppActionDowngrade:
		return "downgrade"
	case AppActionUninstall:
		return "uninstall"
	}
	return ""
}
