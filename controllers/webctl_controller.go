/*
Copyright 2021.

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

package controllers

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"reflect"

	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	cachev1alpha1 "github.com/misteruly/k8s-webctl/api/v1alpha1"
)

// WebctlReconciler reconciles a Webctl object
type WebctlReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=cache.harbor.misteruly.cn,resources=webctls,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cache.harbor.misteruly.cn,resources=webctls/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cache.harbor.misteruly.cn,resources=webctls/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Webctl object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *WebctlReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// your logic here
	webctl := &cachev1alpha1.Webctl{}
	err := r.Get(ctx, req.NamespacedName, webctl)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("webctl resource not found.")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get webctl")
		return ctrl.Result{}, nil
	}
	//检查是否安装了应用，没有就再部署一个
	found := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: webctl.Name, Namespace: webctl.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		dep := r.deploymentforwebctl(webctl)
		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to create new Deployment")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "fail to get deployment")
		return ctrl.Result{}, err
	}
	size := webctl.Spec.Size
	if *found.Spec.Replicas != size {
		found.Spec.Replicas = &size
		err = r.Update(ctx, found)
		if err != nil {
			log.Error(err, "faild to update deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)

			return ctrl.Result{}, err
		}
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	}
	podlist := &corev1.PodList{}
	lisOpts := []client.ListOption{
		client.InNamespace(webctl.Namespace),
		client.MatchingLabels(labelsForwebctl(webctl.Name)),
	}
	if err = r.List(ctx, podlist, lisOpts...); err != nil {
		log.Error(err, "failed to list pods", "webctl.Namespaces", webctl.Namespace, "webctl.Name", webctl.Name)
		return ctrl.Result{}, err
	}
	podnames := getpodnames(podlist.Items)
	if !reflect.DeepEqual(podnames, webctl.Status.Nodes) {
		webctl.Status.Nodes = podnames
		err := r.Status().Update(ctx, webctl)
		if err != nil {
			log.Error(err, "failed to update webctl status")
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

func (r *WebctlReconciler) deploymentforwebctl(m *cachev1alpha1.Webctl) *appsv1.Deployment {
	ls := labelsForwebctl(m.Name)
	replicas := m.Spec.Size
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: "misteruly/hello-node:host1",
						Name:  "k8s-webctl",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 34000,
							Name:          "k8s-webctl",
						}},
					}},
				},
			},
		},
	}
	_ = ctrl.SetControllerReference(m, dep, r.Scheme)
	return dep
}

func labelsForwebctl(name string) map[string]string {
	return map[string]string{"map": "webctl", "webctl_cr": name}
}
func getpodnames(pods []corev1.Pod) []string {
	var podnames []string
	for _, pod := range pods {
		podnames = append(podnames, pod.Name)
	}
	return podnames
}

// SetupWithManager sets up the controller with the Manager.
func (r *WebctlReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1alpha1.Webctl{}).
		Complete(r)
}
