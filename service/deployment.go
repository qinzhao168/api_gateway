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
	"k8s.io/client-go/1.4/pkg/util/intstr"
)

var (
	NameSpace        string = "default"
	app_status       string
	app_status_slice = make([]v1.PodPhase, 5)
	status_Pending   = "Pending"
	status_Running   = "Running"
	status_Succeeded = "Succeeded"
	status_Failed    = "Failed"
	status_Unknown   = "Unknown"
)

func RegisterDeploymentHandler(router *mux.Router) {
	RegisterHttpHandler(router, "/deployment", HTTP_POST, DeploymentApp)           //部署应用
	RegisterHttpHandler(router, "/deployment", HTTP_DELETE, Delete)                //删除应用
	RegisterHttpHandler(router, "/deployment", HTTP_PUT, Update)                   //更新应用配置 包括弹性伸缩  容器个数  启动  停止
	RegisterHttpHandler(router, "/deployment/redployment", HTTP_POST, ReDeloyment) //重新部署
	RegisterHttpHandler(router, "/deployment/status", HTTP_GET, GetAppStatus)
	RegisterHttpHandler(router, "/deployment/containers", HTTP_GET, GetContainers)
}

func DeploymentApp(req *http.Request) (code string, ret interface{}) {
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

	//create replicationController
	rc := new(v1.ReplicationController)

	rcTypeMeta := unversioned.TypeMeta{Kind: "ReplicationController", APIVersion: "v1"}
	rc.TypeMeta = rcTypeMeta

	rcObjectMeta := v1.ObjectMeta{Name: app.Name, Labels: map[string]string{"name": app.Name}}
	rc.ObjectMeta = rcObjectMeta

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
				Containers: []v1.Container{
					v1.Container{
						Name:  app.Name,
						Image: app.Image,
						Ports: []v1.ContainerPort{
							v1.ContainerPort{
								ContainerPort: 9080,
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
	rc.Spec = rcSpec

	result, err := dao.Clientset.Core().ReplicationControllers("default").Create(rc)
	if err != nil {
		log.Errorf("deploy application failed ,the reason is %s", err.Error())
		app.Status = dao.AppFailed
		if err = app.Update(); err != nil {
			log.Errorf("update application status failed,the reason is %s", err.Error())
		}
		code = StatusInternalServerError
		ret = JSON_EMPTY_OBJ
		return
	} else {
		//create service
		service := new(v1.Service)

		svTypemeta := unversioned.TypeMeta{Kind: "Service", APIVersion: "v1"}
		service.TypeMeta = svTypemeta

		svObjectMeta := v1.ObjectMeta{Name: app.Name, Labels: map[string]string{"name": app.Name}}
		service.ObjectMeta = svObjectMeta

		svServiceSpec := v1.ServiceSpec{
			Ports: []v1.ServicePort{
				v1.ServicePort{
					Name:       app.Name,
					Port:       9080,
					TargetPort: intstr.FromInt(9080),
					Protocol:   "TCP",
					// NodePort:   32107,
				},
			},
			Selector: map[string]string{"name": app.Name},
			Type:     v1.ServiceTypeNodePort,
			// LoadBalancerIP: "172.17.11.2",
			// Status: v1.ServiceStatus{
			// 	LoadBalancer: v1.LoadBalancerStatus{
			// 		Ingress: []v1.LoadBalancerIngress{
			// 			v1.LoadBalancerIngress{IP: "172.17.11.2"},
			// 		},
			// 	},
			// },
		}
		service.Spec = svServiceSpec

		result, err := dao.Clientset.Core().Services("default").Create(service)

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

		log.Infof("the result is %v", result)
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

func GetAppStatus(req *http.Request) (code string, ret interface{}) {
	generateName := req.FormValue("appName") + "-"
	podList, err := dao.Clientset.Core().Pods("default").List(api.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	//get every container's status of pods
	for i := 0; i < len(podList.Items); i++ {
		if podList.Items[i].ObjectMeta.GenerateName == generateName && podList.Items[i].ObjectMeta.Namespace == NameSpace {
			app_status_slice = append(app_status_slice, podList.Items[i].Status.Phase)
		}
	}

	//determine the app status by all the container's status of pod
	for _, status := range app_status_slice {
		if status == v1.PodPending {
			app_status = status_Pending
			break
		}

		if status == v1.PodRunning {
			app_status = status_Running
		}

		if status == v1.PodSucceeded {
			app_status = status_Succeeded
		}

		if status == v1.PodFailed || status == v1.PodUnknown {
			app_status = status_Failed
			break
		}
	}

	code = StatusOK
	ret = app_status
	return
}

func GetContainers(req *http.Request) (code string, ret interface{}) {
	generateName := req.FormValue("appName") + "-"

	podList, err := dao.Clientset.Core().Pods("default").List(api.ListOptions{})
	if err != nil {
		code = StatusInternalServerError
		ret = JSON_EMPTY_ARRAY
		log.Error(err.Error())
	}

	//get every container's status of pods
	var containers []dao.Container
	var pod v1.Pod
	for _, pod = range podList.Items {
		//Get the pod of app's
		if pod.ObjectMeta.GenerateName == generateName && pod.ObjectMeta.Namespace == NameSpace {
			container := dao.Container{}
			container.Name = pod.ObjectMeta.Name
			container.Image = pod.Spec.Containers[0].Image
			cpu := pod.Spec.Containers[0].Resources.Requests[v1.ResourceCPU]
			container.Cpu = cpu.String()
			memory := pod.Spec.Containers[0].Resources.Requests[v1.ResourceCPU]
			container.Memory = memory.String()
			container.Created = pod.Status.ContainerStatuses[0].State.Running.StartedAt.String()
			container.Ports = pod.Spec.Containers[0].Ports
			container.Envs = pod.Spec.Containers[0].Env
			container.IntranetIp = pod.Status.PodIP
			container.ExtranetIp = pod.Status.HostIP
			container.Mounts = pod.Spec.Containers[0].VolumeMounts

			containers = append(containers, container)
		}
	}

	code = StatusOK
	ret = containers
	return
}
