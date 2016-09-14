package service

import (
	"api_gateway/dao"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"k8s.io/client-go/1.4/pkg/api"
	"k8s.io/client-go/1.4/pkg/api/resource"
	"k8s.io/client-go/1.4/pkg/api/unversioned"
	"k8s.io/client-go/1.4/pkg/api/v1"
)

func RegisterDeploymentHandler(router *mux.Router) {
	RegisterHttpHandler(router, "/deployment", HTTP_POST, Create)                  //创建应用
	RegisterHttpHandler(router, "/deployment", HTTP_DELETE, Delete)                //删除应用
	RegisterHttpHandler(router, "/deployment", HTTP_PUT, Update)                   //更新应用配置 包括弹性伸缩  容器个数  启动  停止
	RegisterHttpHandler(router, "/deployment/redployment", HTTP_POST, ReDeloyment) //重新部署
}

func Create(req *http.Request) (code string, ret interface{}) {
	decoder := json.NewDecoder(req.Body)
	app := &dao.App{}
	err := decoder.Decode(app)
	if err != nil {
		log.Errorf("请求参数有误：%s", err.Error())
		code = StatusBadRequest
		ret = map[string]interface{}{"success": false, "reason": "请求参数有误，请检查！"}
		return
	}

	app.Status = dao.AppBuilding
	if err = app.Insert(); err != nil {
		log.Errorf("插入数据库失败：%s", err.Error())
		code = StatusInternalServerError
		ret = map[string]interface{}{"success": false, "reason": "插入数据库失败！"}
		return
	}

	//创建replicationController
	rcTypeMeta := unversioned.TypeMeta{Kind: "ReplicationController", APIVersion: "v1"}

	rcObjectMeta := v1.ObjectMeta{
		Name: app.Name,
		Labels: map[string]string{
			"name": app.Name,
		},
	}

	rcSpec := v1.ReplicationControllerSpec{
		Replicas: &app.InstanceCount,
		Selector: map[string]string{
			"name": app.Name,
		},
		Template: &v1.PodTemplateSpec{
			v1.ObjectMeta{
				Name: app.Name,
				Labels: map[string]string{
					"name": app.Name,
				},
			},
			v1.PodSpec{
				RestartPolicy: v1.RestartPolicyAlways,
				NodeSelector: map[string]string{
					"name": app.Name,
				},
				Containers: []v1.Container{
					v1.Container{
						Name:  app.Name,
						Image: app.Image,
						Ports: []v1.ContainerPort{
							v1.ContainerPort{
								ContainerPort: 6379,
								Protocol:      v1.ProtocolTCP,
							},
						},
						Resources: v1.ResourceRequirements{
							Requests: v1.ResourceList{
								v1.ResourceCPU:    resource.MustParse(app.Cpu),
								v1.ResourceMemory: resource.MustParse(app.Memory),
							},
						},
					},
				},
			},
		},
	}

	rc := new(v1.ReplicationController)
	rc.TypeMeta = rcTypeMeta
	rc.ObjectMeta = rcObjectMeta
	rc.Spec = rcSpec

	result, err := dao.Clientset.Core().ReplicationControllers("").Create(rc)
	if err != nil {
		log.Errorf("deploy application failed ,the reason is %s", err.Error())
		app.Status = dao.AppFailed
		if err = app.Update(); err != nil {
			log.Errorf("update application status failed,the reason is %s", err.Error())
		}
		code = StatusInternalServerError
		ret = JSON_EMPTY_OBJ
		return
	}

	app.Status = dao.AppSuccessed
	if err = app.Update(); err != nil {
		log.Errorf("update application status failed,the reason is %s", err.Error())
		code = StatusInternalServerError
		ret = JSON_EMPTY_OBJ
		return
	}

	code = StatusCreated
	ret = result
	return
}

func Delete(req *http.Request) (code string, ret interface{}) {
	decoder := json.NewDecoder(req.Body)
	app := &dao.App{}
	err := decoder.Decode(app)
	if err != nil {
		log.Errorf("请求参数有误：%s", err.Error())
		code = StatusBadRequest
		ret = map[string]interface{}{"success": false, "reason": "请求参数有误，请检查！"}
		return
	}

	err = dao.Clientset.Core().ReplicationControllers("").Delete("name", &api.DeleteOptions{})
	if err != nil {
		log.Errorf("delete application failed ：%s", err.Error())
		code = StatusInternalServerError
		ret = map[string]interface{}{"success": false, "reason": err.Error()}
		return
	}

	code = StatusNoContent
	ret = OK
	return
}

func Update(req *http.Request) (code string, ret interface{}) {
	vebrType := req.FormValue("vebrType")

	decoder := json.NewDecoder(req.Body)
	app := &dao.App{}
	err := decoder.Decode(app)
	if err != nil {
		log.Errorf("请求参数有误：%s", err.Error())
		code = StatusBadRequest
		ret = map[string]interface{}{"success": false, "reason": "请求参数有误，请检查！"}
		return
	}

	//创建replicationController
	rcTypeMeta := unversioned.TypeMeta{Kind: "ReplicationController", APIVersion: "v1"}

	rcObjectMeta := v1.ObjectMeta{
		Name: app.Name,
		Labels: map[string]string{
			"name": app.Name,
		},
	}

	rcSpec := v1.ReplicationControllerSpec{
		Replicas: &app.InstanceCount,
		Selector: map[string]string{
			"name": app.Name,
		},
		Template: &v1.PodTemplateSpec{
			v1.ObjectMeta{
				Name: app.Name,
				Labels: map[string]string{
					"name": app.Name,
				},
			},
			v1.PodSpec{
				RestartPolicy: v1.RestartPolicyAlways,
				NodeSelector: map[string]string{
					"name": app.Name,
				},
				Containers: []v1.Container{
					v1.Container{
						Name:  app.Name,
						Image: app.Image,
						Ports: []v1.ContainerPort{
							v1.ContainerPort{
								ContainerPort: 6379,
								Protocol:      v1.ProtocolTCP,
							},
						},
						Resources: v1.ResourceRequirements{
							Requests: v1.ResourceList{
								v1.ResourceCPU:    resource.MustParse(app.Cpu),
								v1.ResourceMemory: resource.MustParse(app.Memory),
							},
						},
					},
				},
			},
		},
	}

	rc := new(v1.ReplicationController)
	rc.TypeMeta = rcTypeMeta
	rc.ObjectMeta = rcObjectMeta
	rc.Spec = rcSpec

	if vebrType == "updateStatus" {
		_, err = dao.Clientset.Core().ReplicationControllers("").UpdateStatus(rc)

		if err != nil {
			code = StatusInternalServerError
			ret = JSON_EMPTY_OBJ
			return
		}
	}

	if vebrType == "update" {
		_, err = dao.Clientset.Core().ReplicationControllers("").Update(rc)

		if err != nil {
			code = StatusInternalServerError
			ret = JSON_EMPTY_OBJ
			return
		}
	}

	code = StatusNoContent
	ret = OK
	return
}

func ReDeloyment(req *http.Request) (code string, ret interface{}) {
	//删除
	err := dao.Clientset.Core().ReplicationControllers("").Delete("name", &api.DeleteOptions{})
	if err != nil {
		log.Errorf("delete application failed ：%s", err.Error())
		code = StatusInternalServerError
		ret = map[string]interface{}{"success": false, "reason": err.Error()}
		return
	}

	decoder := json.NewDecoder(req.Body)
	app := &dao.App{}
	err = decoder.Decode(app)
	if err != nil {
		log.Errorf("请求参数有误：%s", err.Error())
		code = StatusBadRequest
		ret = map[string]interface{}{"success": false, "reason": "请求参数有误，请检查！"}
		return
	}

	//创建replicationController
	rcTypeMeta := unversioned.TypeMeta{Kind: "ReplicationController", APIVersion: "v1"}

	rcObjectMeta := v1.ObjectMeta{
		Name: app.Name,
		Labels: map[string]string{
			"name": app.Name,
		},
	}

	rcSpec := v1.ReplicationControllerSpec{
		Replicas: &app.InstanceCount,
		Selector: map[string]string{
			"name": app.Name,
		},
		Template: &v1.PodTemplateSpec{
			v1.ObjectMeta{
				Name: app.Name,
				Labels: map[string]string{
					"name": app.Name,
				},
			},
			v1.PodSpec{
				RestartPolicy: v1.RestartPolicyAlways,
				NodeSelector: map[string]string{
					"name": app.Name,
				},
				Containers: []v1.Container{
					v1.Container{
						Name:  app.Name,
						Image: app.Image,
						Ports: []v1.ContainerPort{
							v1.ContainerPort{
								ContainerPort: 6379,
								Protocol:      v1.ProtocolTCP,
							},
						},
						Resources: v1.ResourceRequirements{
							Requests: v1.ResourceList{
								v1.ResourceCPU:    resource.MustParse(app.Cpu),
								v1.ResourceMemory: resource.MustParse(app.Memory),
							},
						},
					},
				},
			},
		},
	}

	rc := new(v1.ReplicationController)
	rc.TypeMeta = rcTypeMeta
	rc.ObjectMeta = rcObjectMeta
	rc.Spec = rcSpec

	result, err := dao.Clientset.Core().ReplicationControllers("").Create(rc)
	if err != nil {
		log.Errorf("deploy application failed ,the reason is %s", err.Error())
		app.Status = dao.AppFailed
		if err = app.Update(); err != nil {
			log.Errorf("update application status failed,the reason is %s", err.Error())
		}
		code = StatusInternalServerError
		ret = JSON_EMPTY_OBJ
		return
	}

	app.Status = dao.AppSuccessed
	if err = app.Update(); err != nil {
		log.Errorf("update application status failed,the reason is %s", err.Error())
		code = StatusInternalServerError
		ret = JSON_EMPTY_OBJ
		return
	}

	code = StatusCreated
	ret = result
	return
}
