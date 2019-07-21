package driver

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"k8s.io/klog"
	"path/filepath"
)

//go:generate mockgen -destination=../../mocks/mock_node_utils.go -package=mocks github.com/ibm/ibm-block-csi-driver/node/pkg/driver NodeUtilsInterface

type NodeUtilsInterface interface {
	ParseIscsiInitiators(path string) (string, error)
	GetInfoFromPublishContext(publishContext map[string]string, configYaml ConfigFile) (string, int, string, error)
	GetIscsiSessionHostsForArrayIQN(array_iqn string) ([]int, error)
}

type NodeUtils struct {
}

func NewNodeUtils() *NodeUtils {
	return &NodeUtils{}

}

func (n NodeUtils) ParseIscsiInitiators(path string) (string, error) {

	file, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer file.Close()

	file_out, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	fileSplit := strings.Split(string(file_out), "InitiatorName=")
	if len(fileSplit) != 2 {
		return "", fmt.Errorf(ErrorWhileTryingToReadIQN, string(file_out))
	}

	iscsiIqn := strings.TrimSpace(fileSplit[1])

	return iscsiIqn, nil
}

func (n NodeUtils) GetInfoFromPublishContext(publishContext map[string]string, configYaml ConfigFile) (string, int, string, error) {
	// this will return :  connectivityType, lun, array_iqn, error
	str_lun := publishContext[configYaml.Controller.Publish_context_lun_parameter]
	
	lun, err := strconv.Atoi(str_lun)
	if err != nil {
		return "", -1, "", err
	}
	
	connectivityType := publishContext[configYaml.Controller.Publish_context_connectivity_parameter]
	array_iqn := publishContext[configYaml.Controller.Publish_context_array_iqn]

	return connectivityType, lun, array_iqn, nil
}


func (n NodeUtils) GetIscsiSessionHostsForArrayIQN(array_iqn string) ([]int, error){
	sysPath := "/sys/class/iscsi_host/"
	var sessionHosts []int
	if hostDirs, err := ioutil.ReadDir(sysPath); err != nil {
		klog.Errorf("cannot read sys dir : {%v}. error : {%v}", sysPath, err)
		return sessionHosts, err
	} else {
		klog.V(5).Infof("host dirs : {%v}", hostDirs)
		for _, hostDir := range hostDirs {
			// get the host session number : "host34"
			hostName := hostDir.Name()
			hostNumber := -1
			if !strings.HasPrefix(hostName, "host") {
				continue
			} else{
				hostNumber, err = strconv.Atoi(strings.TrimPrefix(hostName, "host"))
				if err != nil {
					klog.V(4).Infof("cannot get host id from host : {%V}", hostName)
					continue
				 }
			}

			targetPath := sysPath + hostName + "/device/session*/iscsi_session/session*/targetname" 
			
			//devicePath + sessionName + "/iscsi_session/" + sessionName + "/targetname"
			matches, err := filepath.Glob(targetPath)
			if err != nil{
				klog.Errorf("error while finding targetPath : {%V}. err : {%v}", targetPath, err)
				return sessionHosts, err
			}
			
			klog.V(5).Infof("matches were found : {%V}", matches)
			
			//TODO: can there be more then 1 session?? 
			//sessionNumber, err :=  strconv.Atoi(strings.TrimPrefix(matches[0], "session"))
			
			if len(matches)  == 0{
				klog.V(4).Infof("could not find targe name for host : {%v}, path : {%v}", hostName, targetPath)
				continue
			}
			
			targetNamePath := matches[0]
			targetName, err := ioutil.ReadFile(targetNamePath)
			if err != nil {
				klog.V(4).Infof("could not read target name from file : {%v}, error : {%v}", targetNamePath, err)
				continue
			}
			
			klog.V(5).Infof("target name found : {%V}", targetName)

			 if strings.TrimSpace(string(targetName)) == array_iqn{
			 	sessionHosts = append(sessionHosts, hostNumber)
			 	klog.V(5).Infof("host nunber appended : {%V}. sessionhosts is : {%v}", hostNumber, sessionHosts)
			 }
		}
		
		return sessionHosts, nil
	}
}
