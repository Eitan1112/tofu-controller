package tfctl

import (
	"context"
	"fmt"
	"io"
	"time"

	infrav1 "github.com/flux-iac/tofu-controller/api/v1alpha2"
	"github.com/fluxcd/pkg/apis/meta"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Reconcile annotates the given object
func (c *CLI) Reconcile(out io.Writer, resource string) error {
	key := types.NamespacedName{
		Name:      resource,
		Namespace: c.namespace,
	}

	err := requestReconciliation(context.TODO(), c.client, key)
	if err != nil {
		return err
	}

	fmt.Fprintf(out, " Reconcile requested for %s/%s\n", c.namespace, resource)

	return nil
}

func requestReconciliation(ctx context.Context, kubeClient client.Client, namespacedName types.NamespacedName) error {
	return retry.RetryOnConflict(retry.DefaultBackoff, func() (err error) {
		terraform := &infrav1.Terraform{}
		if err := kubeClient.Get(ctx, namespacedName, terraform); err != nil {
			return err
		}
		patch := client.MergeFrom(terraform.DeepCopy())
		if ann := terraform.GetAnnotations(); ann == nil {
			terraform.SetAnnotations(map[string]string{
				meta.ReconcileRequestAnnotation: time.Now().Format(time.RFC3339Nano),
			})
		} else {
			ann[meta.ReconcileRequestAnnotation] = time.Now().Format(time.RFC3339Nano)
			terraform.SetAnnotations(ann)
		}
		return kubeClient.Patch(ctx, terraform, patch)
	})
}
