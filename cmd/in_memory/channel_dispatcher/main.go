/*
Copyright 2019 The Knative Authors

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
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"os"

	filteredFactory "knative.dev/pkg/client/injection/kube/informers/factory/filtered"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/injection/sharedmain"
	"knative.dev/pkg/signals"

	"knative.dev/eventing/pkg/eventingtls"
	inmemorychannel "knative.dev/eventing/pkg/reconciler/inmemorychannel/dispatcher"
)

func main() {
	ctx := signals.NewContext()
	ns := os.Getenv("NAMESPACE")
	if ns != "" {
		ctx = injection.WithNamespaceScope(ctx, ns)
	}

	ctx = filteredFactory.WithSelectors(ctx,
		eventingtls.TrustBundleLabelSelector,
	)

	sharedmain.MainWithContext(ctx, "inmemorychannel-dispatcher",
		inmemorychannel.NewController,
	)
}
