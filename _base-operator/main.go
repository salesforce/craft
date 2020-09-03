// Copyright (c) 2020, salesforce.com, inc.
// All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// For full license text, see LICENSE.txt file in the repo root or https://opensource.org/licenses/BSD-3-Clause

/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"os"
	"strconv"

	"{{ .Repo }}/controllers"

	mygroup{{ .Version }} "{{ .Repo }}/api/{{ .Version }}"
	"{{ .Repo }}/reconciler"
	"k8s.io/apimachinery/pkg/runtime"
	uberzap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = mygroup{{ .Version }}.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	logger := uberzap.NewAtomicLevelAt(zapcore.InfoLevel)
	//ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
	ctrl.SetLogger(zap.New(func(o *zap.Options) {
		o.Development = true
		o.Level =    &logger
	}))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:	"{{.LeaderElectionID}}",
		Port:               9443,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	var requeueAfter int
	freq_min, err := strconv.Atoi("{{.ReconcileFreq}}")
	
	if err != nil {
		requeueAfter = 30000
	} else {
		requeueAfter = freq_min * 60 * 1000		//mins to milliseconds
	}
	controllerParams := reconciler.ReconcileParameters{
		RequeueAfter	  	: requeueAfter,
		RequeueAfterSuccess	: 15000,
		RequeueAfterFailure	: 30000,
	}
	if err = (&controllers.ControllerFactory{
		ResourceManagerCreator: controllers.CreateResourceManager,
		Scheme:                 scheme,
	}).SetupWithManager(mgr, controllerParams, nil); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "B")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
