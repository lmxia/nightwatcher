package gaia

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lmxia/nightwatcher/app"
	"github.com/lmxia/nightwatcher/controllers/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

type ClusterUri struct {
	Namespace string `uri:"namespace" binding:"required"`
	Cluster   string `uri:"cluster" binding:"required"`
}

type LabelUri struct {
	Label string `uri:"label" binding:"required"`
}

// @Summary 查看全部clusters
// @accept application/json
// @Produce  application/json
// @Param cluster path string true "Cluster"
// @Param namespace path string true "Namespace"
// @Param deploymentName path string true "DeploymentName"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /clusters [get]
func GetClusters(c *gin.Context) {
	appG := app.Gin{C: c}
	k8sClient, err := k8s.GetClientWithPanic()
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	clusters, err := k8sClient.Gaiaclient.PlatformV1alpha1().ManagedClusters(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", clusters)
}

// @Summary 查看全部具体某个cluster 的label的value
// @accept application/json
// @Produce  application/json
// @Param cluster path string true "Cluster"
// @Param namespace path string true "Namespace"
// @Param deploymentName path string true "DeploymentName"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /clusters/{namespace}/{cluster}/{label} [get]
func GetClusterLabel(c *gin.Context) {
	appG := app.Gin{C: c}

	var (
		l LabelUri
		u ClusterUri
	)

	if err := appG.C.ShouldBindUri(&u); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}
	if err := appG.C.ShouldBindUri(&l); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}

	k8sClient, err := k8s.GetClientWithPanic()
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	clusters, err := k8sClient.Gaiaclient.PlatformV1alpha1().ManagedClusters(u.Namespace).Get(context.TODO(), u.Cluster, metav1.GetOptions{})
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	netEnvironmentMap, nodeRoleMap, resFormMap, runtimeStateMap, snMap, geolocationMap, providers := clusters.GetHypernodeLabelsMapFromManagedCluster()

	if l.Label == "supplier-name" {
		keys := make([]string, 0, len(providers))
		for k := range providers {
			keys = append(keys, k)
		}
		appG.Success(http.StatusOK, "ok", keys)
	} else if l.Label == "geo-location" {
		keys := make([]string, 0, len(geolocationMap))
		for k := range geolocationMap {
			keys = append(keys, k)
		}
		appG.Success(http.StatusOK, "ok", keys)
	} else if l.Label == "net-environment" {
		keys := make([]string, 0, len(netEnvironmentMap))
		for k := range netEnvironmentMap {
			keys = append(keys, k)
		}
		appG.Success(http.StatusOK, "ok", keys)
	} else if l.Label == "node-role" {
		keys := make([]string, 0, len(nodeRoleMap))
		for k := range nodeRoleMap {
			keys = append(keys, k)
		}
		appG.Success(http.StatusOK, "ok", keys)
	} else if l.Label == "runtime-state" {
		keys := make([]string, 0, len(runtimeStateMap))
		for k := range runtimeStateMap {
			keys = append(keys, k)
		}
		appG.Success(http.StatusOK, "ok", keys)
	} else if l.Label == "sn" {
		keys := make([]string, 0, len(snMap))
		for k := range snMap {
			keys = append(keys, k)
		}
		appG.Success(http.StatusOK, "ok", keys)
	} else if l.Label == "res-form" {
		keys := make([]string, 0, len(resFormMap))
		for k := range resFormMap {
			keys = append(keys, k)
		}
		appG.Success(http.StatusOK, "ok", keys)
	} else {
		appG.Fail(http.StatusBadRequest, errors.New("unsupported labels"), "")
	}
}
