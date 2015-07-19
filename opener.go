package slacktv

import "os/exec"

func open(url string) (err error) {
	cmd := cmdForUrl(url)
	dbg("%s", cmd.Args)
	err = cmd.Start()
	if err != nil {
		return
	}
	err = cmd.Wait()
	return
}

func cmdForUrl(url string) *exec.Cmd {
	return exec.Command(
		"open",
		"-n",
		"-a",
		"Google Chrome",
		"--args",
		"--start-fullscreen",
		"--kiosk",
		url,
	)
}
