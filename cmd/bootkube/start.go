package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/kubernetes-incubator/bootkube/pkg/bootkube"
	"github.com/kubernetes-incubator/bootkube/pkg/util"
	"github.com/kubernetes-incubator/bootkube/pkg/util/etcdutil"
)

var (
	cmdStart = &cobra.Command{
		Use:          "start",
		Short:        "Start the bootkube service",
		Long:         "",
		PreRunE:      validateStartOpts,
		RunE:         runCmdStart,
		SilenceUsage: true,
	}

	startOpts struct {
		assetDir              string
		etcdServer            string
		etcdAuthEnabled       bool
		selfHostedEtcd        bool
		serviceClusterIPRange net.IPNet
		clusterCIDR           net.IPNet
	}
)

func init() {
	cmdRoot.AddCommand(cmdStart)
	cmdStart.Flags().StringVar(&startOpts.etcdServer, "etcd-server", "http://127.0.0.1:2379", "Single etcd node to use during bootkube bootstrap process.")
	cmdStart.Flags().BoolVar(&startOpts.etcdAuthEnabled, "etcd-auth-enabled", false, "Etcd requires authentication through client certificates.")
	cmdStart.Flags().StringVar(&startOpts.assetDir, "asset-dir", "", "Path to the cluster asset directory. Expected layout genereted by the `bootkube render` command.")
	cmdStart.Flags().BoolVar(&startOpts.selfHostedEtcd, "experimental-self-hosted-etcd", false, "Self hosted etcd mode. Includes starting the initial etcd member by bootkube.")
	cmdStart.Flags().IPNetVar(&startOpts.serviceClusterIPRange, "service-cluster-ip-range", net.IPNet{IP: net.IP("10.3.0.0"), Mask: net.IPMask("24")}, "A CIDR notation IP range from which to assign service cluster IPs.")
	cmdStart.Flags().IPNetVar(&startOpts.clusterCIDR, "cluster-cidr", net.IPNet{IP: net.IP("10.2.0.0"), Mask: net.IPMask("16")}, "CIDR Range for Pods in cluster.")
}

func runCmdStart(cmd *cobra.Command, args []string) error {
	etcdServer, err := url.Parse(startOpts.etcdServer)
	if err != nil {
		return fmt.Errorf("Invalid etcd etcdServer %q: %v", startOpts.etcdServer, err)
	}

	// TODO: this should likely move into bootkube.Run() eventually.
	if startOpts.selfHostedEtcd {
		if err := etcdutil.StartEtcd(startOpts.etcdServer); err != nil {
			return fmt.Errorf("fail to start etcd: %v", err)
		}
	}

	bk, err := bootkube.NewBootkube(bootkube.Config{
		AssetDir:              startOpts.assetDir,
		EtcdServer:            etcdServer,
		EtcdAuthEnabled:       startOpts.etcdAuthEnabled,
		SelfHostedEtcd:        startOpts.selfHostedEtcd,
		ServiceClusterIPRange: startOpts.serviceClusterIPRange,
		ClusterCIDR:           startOpts.clusterCIDR,
	})

	if err != nil {
		return err
	}

	// set in util init() func, but lets not depend on that
	flag.Set("logtostderr", "true")
	util.InitLogs()
	defer util.FlushLogs()

	return bk.Run()
}

func validateStartOpts(cmd *cobra.Command, args []string) error {
	if startOpts.etcdServer == "" {
		return errors.New("missing required flag: --etcd-server")
	}
	if startOpts.assetDir == "" {
		return errors.New("missing required flag: --asset-dir")
	}
	return nil
}
