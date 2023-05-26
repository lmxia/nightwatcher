package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lmxia/nightwatcher/app"
	"github.com/lmxia/nightwatcher/controllers/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ServicesQuery struct {
	Namespace string `form:"namespace"`
}

type ServiceQuery struct {
}

type ServiceUri struct {
	Namespace   string `uri:"namespace" binding:"required"`
	ServiceName string `uri:"serviceName" binding:"required"`
}

func GetServices(c *gin.Context) {
	appG := app.Gin{C: c}
	var q ServicesQuery
	if err := appG.C.ShouldBindQuery(&q); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}

	k8sClient, err := k8s.GetClientWithPanic()
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
	}

	services, err := k8sClient.K8sClient.CoreV1().Services(q.Namespace).List(context.TODO(), metav1.ListOptions{})
	for i := 0; i < len(services.Items); i++ {
		services.Items[i].CreationTimestamp = metav1.NewTime(services.Items[i].CreationTimestamp.Add(8 * time.Hour))
	}

	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", services)
}

func GetService(c *gin.Context) {
	appG := app.Gin{C: c}
	var u ServiceUri
	if err := appG.C.ShouldBindUri(&u); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}

	k8sClient, err := k8s.GetClientWithPanic()
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
	}

	service, err := k8sClient.K8sClient.CoreV1().Services(u.Namespace).Get(context.TODO(), u.ServiceName, metav1.GetOptions{})
	service.CreationTimestamp = metav1.NewTime(service.CreationTimestamp.Add(8 * time.Hour))

	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	appG.Success(http.StatusOK, "ok", service)
}
