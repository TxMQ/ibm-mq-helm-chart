package util

import (
	"fmt"
	"os"
	"szesto.com/mqrunner/logger"
	"time"
)

const _mqscic = "/etc/mqm/mqscic.mqsc"
const _qmini = "/etc/mqm/qmini.ini"

func GetMqscic() string {
	return _mqscic
}

func GetQmini() string {
	return _qmini
}

func MergeMqscFiles() error {
	if GetDebugFlag() {
		logger.Logmsg("merging mqsc, ini, and yaml files")
	}

	t := time.Now()
	defer logger.Logmsg(fmt.Sprintf("time to merge files: %v", time.Since(t)))

	// mq yaml file is mounted on /etc/mqm/mqyaml
	mqyamlDir := "/etc/mqm/mqyaml"
	mqyamlOutFile := GetMqscic()

	// delete existing mqscic file
	err := deleteFile(mqyamlOutFile)
	if err != nil {
		logger.Logmsg(err)
	}

	err = MqYamlMerge(mqyamlDir, mqyamlOutFile)
	if err != nil {
		return err
	}

	// mqscic file is mounted on /etc/mqm/mqsc
	mqscicDir := "/etc/mqm/mqsc"
	mqscicOutFile := GetMqscic()

	err = MqscMerge(mqscicDir, mqscicOutFile)
	if err != nil {
		return err
	}

	// qmini file is mounted on /etc/mqm/qmini
	qminiDir := "/etc/mqm/qmini"
	qminiOutFile := GetQmini()

	err = QminiMerge(qminiDir, qminiOutFile)
	if err != nil {
		return err
	}

	return nil
}

func MergeGitConfigFiles2(fc FetchConfig) error {

	if GetDebugFlag() {
		logger.Logmsg(fmt.Sprintf("giturl '%s', gitref '%s', gitdir '%s'", fc.Url, fc.ReferenceName, fc.Dir))
	}

	// protect againts "''" case
	if len(fc.Url) <= 2 {
		return nil
	}

	mqscicOutFile := GetMqscic()
	qminiOutFile := GetQmini()

	err := FetchMergeConfigFiles(fc, mqscicOutFile, qminiOutFile)
	if err != nil {
		return err
	}

	return nil
}

func MergeGitConfigFiles() error {
	if GetDebugFlag() {
		logger.Logmsg("merging git config files")
	}

	t := time.Now()
	defer logger.Logmsg(fmt.Sprintf("time to merge git config files: %v", time.Since(t)))

	fc := FetchConfig{
		Url:           os.Getenv("GIT_CONFIG_URL"),
		ReferenceName: os.Getenv("GIT_CONFIG_REF"),
		Dir:           os.Getenv("GIT_CONFIG_DIR"),
	}

	return MergeGitConfigFiles2(fc)
}
