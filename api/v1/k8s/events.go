package v1

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lmxia/nightwatcher/app"
	"github.com/lmxia/nightwatcher/controllers/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type EventsQuery struct {
	Namespace string `json:"namespace" form:"namespace" binding:"required"`
	Name      string `json:"name" form:"name" binding:"required"`
	Kind      string `json:"kind" form:"kind" binding:"required"`
	Uid       string `json:"uid" form:"uid" binding:"required"`
}

func GetEvents(c *gin.Context) {
	appG := app.Gin{C: c}
	var (
		q        EventsQuery
		listOpts metav1.ListOptions
	)
	if err := appG.C.ShouldBindQuery(&q); err != nil {
		appG.Fail(http.StatusBadRequest, err, nil)
		return
	}

	k8sClient, err := k8s.GetClientWithPanic()
	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}
	listOpts.FieldSelector = fmt.Sprintf(
		"involvedObject.name=%s,involvedObject.namespace=%s,involvedObject.kind=%s,involvedObject.uid=%s",
		q.Name, q.Namespace, q.Kind, q.Uid,
	)
	listOpts.TypeMeta = metav1.TypeMeta{Kind: q.Kind}
	events, err := k8sClient.K8sClient.CoreV1().Events(q.Namespace).List(context.TODO(), listOpts)
	for i := 0; i < len(events.Items); i++ {
		events.Items[i].CreationTimestamp = metav1.NewTime(events.Items[i].CreationTimestamp.Add(8 * time.Hour))
	}

	if err != nil {
		appG.Fail(http.StatusInternalServerError, err, nil)
		return
	}

	appG.Success(http.StatusOK, "ok", events)
}
