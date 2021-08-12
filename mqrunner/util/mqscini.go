package util

import (
	"log"
	"os"
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

	// mq yaml file is mounted on /etc/mqm/mqyaml
	mqyamlDir := "/etc/mqm/mqyaml"
	mqyamlOutFile := GetMqscic()

	err := MqYamlMerge(mqyamlDir, mqyamlOutFile)
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

	err =QminiMerge(qminiDir, qminiOutFile)
	if err != nil {
		return err
	}

	return nil
}

func MergeGitConfigFiles2(fc FetchConfig) error {

	if GetDebugFlag() {
		log.Printf("run-main: giturl '%s', gitref '%s', gitdir '%s'\n", fc.Url, fc.ReferenceName, fc.Dir)
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

	fc := FetchConfig{
		Url:           os.Getenv("GIT_CONFIG_URL"),
		ReferenceName: os.Getenv("GIT_CONFIG_REF"),
		Dir:           os.Getenv("GIT_CONFIG_DIR"),
	}

	return MergeGitConfigFiles2(fc)
}
