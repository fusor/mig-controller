package compat

import (
	"context"
	"strconv"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	batchv1beta "k8s.io/api/batch/v1beta1"
	batchv2alpha "k8s.io/api/batch/v2alpha1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	dapi "k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
)

//
// A smart client.
// Provides seamless API version compatibility.
type Client struct {
	k8sclient.Client
	dapi.DiscoveryInterface
	// major k8s version.
	Major int
	// minor k8s version.
	Minor int
}

//
// Create a new client.
func NewClient(restCfg *rest.Config) (k8sclient.Client, error) {
	rClient, err := k8sclient.New(
		restCfg,
		k8sclient.Options{
			Scheme: scheme.Scheme,
		})
	if err != nil {
		return nil, err
	}
	dClient, err := dapi.NewDiscoveryClientForConfig(restCfg)
	if err != nil {
		return nil, err
	}
	version, err := dClient.ServerVersion()
	if err != nil {
		return nil, err
	}
	major, err := strconv.Atoi(version.Major)
	if err != nil {
		return nil, err
	}
	minor, err := strconv.Atoi(strings.Trim(version.Minor, "+"))
	if err != nil {
		return nil, err
	}
	nClient := &Client{
		Client:             rClient,
		DiscoveryInterface: dClient,
		Major:              major,
		Minor:              minor,
	}

	return nClient, nil
}

//
// supportedVersion will determine correct version of the object provided, based on cluster version
func (c Client) supportedVersion(obj runtime.Object) runtime.Object {
	if c.Minor < 16 {
		switch obj.(type) {

		// Deployment
		case *appsv1.Deployment:
			return &appsv1beta1.Deployment{}
		case *appsv1.DeploymentList:
			return &appsv1beta1.DeploymentList{}

		// DaemonSet
		case *appsv1.DaemonSet:
			return &extv1beta1.DaemonSet{}
		case *appsv1.DaemonSetList:
			return &extv1beta1.DaemonSetList{}

		// ReplicaSet
		case *appsv1.ReplicaSet:
			return &extv1beta1.ReplicaSet{}
		case *appsv1.ReplicaSetList:
			return &extv1beta1.ReplicaSetList{}

		// StatefulSet
		case *appsv1.StatefulSet:
			return &appsv1beta1.StatefulSet{}
		case *appsv1.StatefulSetList:
			return &appsv1beta1.StatefulSetList{}
		}
	}

	if c.Minor < 8 {
		switch obj.(type) {

		// CronJob
		case *batchv1beta.CronJobList:
			return &batchv2alpha.CronJobList{}
		case *batchv1beta.CronJob:
			return &batchv2alpha.CronJob{}
		}
	}

	return obj
}

//
// Down convert a resource as needed based on cluster version.
func (c Client) downConvert(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	new := c.supportedVersion(obj)
	if new == obj {
		return obj, nil
	}

	err := scheme.Scheme.Convert(obj, new, ctx)
	if err != nil {
		return nil, err
	}

	return new, nil
}

//
// upConvert will convert src resource to dst as needed based on cluster version
func (c Client) upConvert(ctx context.Context, src runtime.Object, dst runtime.Object) error {
	if c.supportedVersion(dst) == dst {
		dst = src
		return nil
	}

	return scheme.Scheme.Convert(src, dst, ctx)
}

//
// Get the specified resource.
// The resource will be converted to a compatible version as needed.
func (c Client) Get(ctx context.Context, key k8sclient.ObjectKey, in runtime.Object) error {
	obj := c.supportedVersion(in)
	err := c.Client.Get(ctx, key, obj)
	if err != nil {
		return err
	}

	return c.upConvert(ctx, obj, in)
}

//
// List the specified resource.
// The resource will be converted to a compatible version as needed.
func (c Client) List(ctx context.Context, opt *k8sclient.ListOptions, in runtime.Object) error {
	obj, err := c.downConvert(ctx, in)
	if err != nil {
		return err
	}
	err = c.Client.List(ctx, opt, obj)
	if err != nil {
		return err
	}

	return c.upConvert(ctx, obj, in)
}

// Create the specified resource.
// The resource will be converted to a compatible version as needed.
func (c Client) Create(ctx context.Context, in runtime.Object) error {
	obj, err := c.downConvert(ctx, in)
	if err != nil {
		return err
	}
	return c.Client.Create(ctx, obj)
}

// Delete the specified resource.
// The resource will be converted to a compatible version as needed.
func (c Client) Delete(ctx context.Context, in runtime.Object, opt ...k8sclient.DeleteOptionFunc) error {
	obj, err := c.downConvert(ctx, in)
	if err != nil {
		return err
	}
	return c.Client.Delete(ctx, obj, opt...)
}

// Update the specified resource.
// The resource will be converted to a compatible version as needed.
func (c Client) Update(ctx context.Context, in runtime.Object) error {
	obj, err := c.downConvert(ctx, in)
	if err != nil {
		return err
	}
	return c.Client.Update(ctx, obj)
}
