package cmd

import (
	"io/ioutil"

	"github.com/okteto/okteto/pkg/config"
	"github.com/okteto/okteto/pkg/log"
)

type upState string

const (
	activating    upState = "activating"
	starting      upState = "starting"
	attaching     upState = "attaching"
	pulling       upState = "pulling"
	startingSync  upState = "startingSync"
	synchronizing upState = "synchronizing"
	ready         upState = "ready"
	failed        upState = "failed"
)

func (up *UpContext) updateStateFile(state upState) {
	if len(up.Dev.Namespace) == 0 {
		log.Info("can't update state file, namespace is empty")
	}

	if len(up.Dev.Name) == 0 {
		log.Info("can't update state file, name is empty")
	}

	s := config.GetStateFile(up.Dev.Namespace, up.Dev.Name)
	log.Debugf("updating statefile %s: '%s'", s, state)
	if err := ioutil.WriteFile(s, []byte(state), 0644); err != nil {
		log.Infof("can't update state file, %s", err)
	}
}
