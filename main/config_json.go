package main

import (
	"io"
	"os"
	"os/exec"

	"v2ray.com/core"
	"v2ray.com/core/app/log"
	"v2ray.com/core/common/platform"
	jsonconf "v2ray.com/ext/tools/conf/serial"
)

func jsonToProto(input io.Reader) (*core.Config, error) {
	v2ctl := platform.GetToolLocation("v2ctl")
	_, err := os.Stat(v2ctl)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(v2ctl, "config")
	cmd.Stdin = input
	cmd.Stderr = os.Stderr

	stdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	defer stdoutReader.Close()
	return core.LoadConfig(core.ConfigFormat_Protobuf, stdoutReader)
}

func init() {
	core.RegisterConfigLoader(core.ConfigFormat_JSON, func(input io.Reader) (*core.Config, error) {
		config, err := jsonToProto(input)
		if err != nil {
			log.Trace(newError("failed to execute v2ctl to convert config file.").Base(err).AtWarning())
			return jsonconf.LoadJSONConfig(input)
		}
		return config, nil
	})
}
