package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	k3dCluster "github.com/rancher/k3d/v4/pkg/client"
	"github.com/rancher/k3d/v4/pkg/config"
	conf "github.com/rancher/k3d/v4/pkg/config/v1alpha1"
	"github.com/rancher/k3d/v4/pkg/runtimes"
	k3d "github.com/rancher/k3d/v4/pkg/types"
	"github.com/rancher/k3d/v4/version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func CmdEnvCreate() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create environment",
		Long:  `Create environment base on the envorinment configuration.`,
		Run: func(cmd *cobra.Command, args []string) {
			flags := readEnvFlags()
			createEnv(cmd, flags)
		},
	}

	return cmd
}

func createEnv(cmd *cobra.Command, flags envFlags) {

	envConfig := envConfig(flags)

	cliConfig := &conf.SimpleConfig{
		Servers: 0,
		Agents:  0,
		Image:   fmt.Sprintf("%s:%s", k3d.DefaultK3sImageRepo, version.GetK3sVersion(false)),
		Options: conf.SimpleConfigOptions{
			K3dOptions: conf.SimpleConfigOptionsK3d{
				Wait:                       true,
				Timeout:                    0 * time.Second,
				DisableLoadbalancer:        false,
				NoRollback:                 false,
				PrepDisableHostIPInjection: false,
				DisableImageVolume:         false,
			},
			KubeconfigOptions: conf.SimpleConfigOptionsKubeconfig{
				UpdateDefaultKubeconfig: true,
				SwitchCurrentContext:    true,
			},
		},
	}

	if envConfig.Cluster.K3d.Config != "" {
		configFromFile, err := config.ReadConfig(envConfig.Cluster.K3d.Config)
		if err != nil {
			log.Fatalln(err)
		}
		cliConfig, err = config.MergeSimple(*cliConfig, configFromFile.(conf.SimpleConfig))
		if err != nil {
			log.Fatalln(err)
		}
		cliConfig = absVolumes(*cliConfig)
	} else {
		log.Fatal("Missing k3d cluster configuration")
	}

	clusterConfig, err := config.TransformSimpleToClusterConfig(cmd.Context(), runtimes.SelectedRuntime, *cliConfig)
	if err != nil {
		log.Fatalln(err)
	}
	if err := config.ValidateClusterConfig(cmd.Context(), runtimes.SelectedRuntime, *clusterConfig); err != nil {
		log.Fatalln("Failed Cluster Configuration Validation: ", err)
	}

	// check if a cluster with that name exists already
	if _, err := k3dCluster.ClusterGet(cmd.Context(), runtimes.SelectedRuntime, &clusterConfig.Cluster); err == nil {
		log.Fatalf("Failed to create cluster '%s' because a cluster with that name already exists", clusterConfig.Cluster.Name)
	}

	// create cluster
	if clusterConfig.KubeconfigOpts.UpdateDefaultKubeconfig {
		log.Debugln("'--kubeconfig-update-default set: enabling wait-for-server")
		clusterConfig.ClusterCreateOpts.WaitForServer = true
	}
	if err := k3dCluster.ClusterRun(cmd.Context(), runtimes.SelectedRuntime, clusterConfig); err != nil {
		// rollback if creation failed
		log.Errorln(err)
		if cliConfig.Options.K3dOptions.NoRollback { // TODO: move rollback mechanics to pkg/
			log.Fatalln("Cluster creation FAILED, rollback deactivated.")
		}
		// rollback if creation failed
		log.Errorln("Failed to create cluster >>> Rolling Back")
		if err := k3dCluster.ClusterDelete(cmd.Context(), runtimes.SelectedRuntime, &clusterConfig.Cluster); err != nil {
			log.Errorln(err)
			log.Fatalln("Cluster creation FAILED, also FAILED to rollback changes!")
		}
		log.Fatalln("Cluster creation FAILED, all changes have been rolled back!")
	}
	log.Infof("Cluster '%s' created successfully!", clusterConfig.Cluster.Name)

	if clusterConfig.KubeconfigOpts.UpdateDefaultKubeconfig && clusterConfig.KubeconfigOpts.SwitchCurrentContext {
		log.Infoln("--kubeconfig-update-default=false --> sets --kubeconfig-switch-context=false")
		clusterConfig.KubeconfigOpts.SwitchCurrentContext = false
	}

	if clusterConfig.KubeconfigOpts.UpdateDefaultKubeconfig {
		log.Debugf("Updating default kubeconfig with a new context for cluster %s", clusterConfig.Cluster.Name)
		if _, err := k3dCluster.KubeconfigGetWrite(cmd.Context(), runtimes.SelectedRuntime, &clusterConfig.Cluster, "", &k3dCluster.WriteKubeConfigOptions{UpdateExisting: true, OverwriteExisting: false, UpdateCurrentContext: cliConfig.Options.KubeconfigOptions.SwitchCurrentContext}); err != nil {
			log.Warningln(err)
		}
	}

	log.Infoln("You can now use it like this:")
	if clusterConfig.KubeconfigOpts.UpdateDefaultKubeconfig && !clusterConfig.KubeconfigOpts.SwitchCurrentContext {
		fmt.Printf("kubectl config use-context %s\n", fmt.Sprintf("%s-%s", k3d.DefaultObjectNamePrefix, clusterConfig.Cluster.Name))
	} else if !clusterConfig.KubeconfigOpts.SwitchCurrentContext {
		if runtime.GOOS == "windows" {
			fmt.Printf("$env:KUBECONFIG=(%s kubeconfig write %s)\n", os.Args[0], clusterConfig.Cluster.Name)
		} else {
			fmt.Printf("export KUBECONFIG=$(%s kubeconfig write %s)\n", os.Args[0], clusterConfig.Cluster.Name)
		}
	}
	fmt.Println("kubectl cluster-info")
}

func absVolumes(p conf.SimpleConfig) *conf.SimpleConfig {
	for i, vol := range p.Volumes {
		s := strings.Split(vol.Volume, ":")
		a, e := filepath.Abs(s[0])
		if e != nil {
			log.WithField("volume", vol.Volume).Fatal("Error convert volume to absolute path")
		}
		n := a + ":" + s[1]
		log.WithFields(log.Fields{"volume": vol.Volume, "new": n}).Info("Update volume to absolute path")
		p.Volumes[i].Volume = n
	}
	return &p
}
