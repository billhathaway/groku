package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/codegangsta/cli"
)

const VERSION = "0.1"

func main() {
	app := cli.NewApp()
	app.Name = "groku"
	app.Version = VERSION
	app.Usage = "roku CLI remote"
	app.Commands = commands()
	app.Commands = append(app.Commands, discover())
	app.Run(os.Args)
}

func findRoku() string {
	ssdp, _ := net.ResolveUDPAddr("udp", "239.255.255.250:1900")
	addr, _ := net.ResolveUDPAddr("udp", ":0")
	socket, _ := net.ListenUDP("udp", addr)

	socket.WriteToUDP([]byte("M-SEARCH * HTTP/1.1\r\n"+
		"HOST: 239.255.255.250:1900\r\n"+
		"MAN: \"ssdp:discover\"\r\n"+
		"ST: roku:ecp\r\n"+
		"MX: 3 \r\n\r\n"), ssdp)

	answerBytes := make([]byte, 1024)
	socket.ReadFromUDP(answerBytes[:])

	ret := strings.Split(string(answerBytes), "\r\n")
	return strings.TrimPrefix(ret[len(ret)-3], "LOCATION: ")
}

func commands() []cli.Command {
	cmds := []cli.Command{}
	for _, cmd := range []string{
		"Home",
		"Rev",
		"Fwd",
		"Select",
		"Left",
		"Right",
		"Down",
		"Up",
		"Back",
		"Info",
		"Backspace",
		"Search",
		"Enter",
	} {
		cmds = append(cmds, cli.Command{
			Name:  strings.ToLower(cmd),
			Usage: strings.ToLower(cmd),
			Action: func(c *cli.Context) {
				http.PostForm(fmt.Sprintf("%vkeypress/%v", findRoku(), cmd), nil)
			},
		})
	}
	cmds = append(cmds, cli.Command{
		Name:  "replay",
		Usage: "replay",
		Action: func(c *cli.Context) {
			http.PostForm(fmt.Sprintf("%vkeypress/%v", findRoku(), "InstantReplay"), nil)
		},
	})
	cmds = append(cmds, cli.Command{
		Name:  "play",
		Usage: "play/pause",
		Action: func(c *cli.Context) {
			http.PostForm(fmt.Sprintf("%vkeypress/%v", findRoku(), "Play"), nil)
		},
	})
	return cmds
}

func discover() cli.Command {
	return cli.Command{
		Name:  "discover",
		Usage: "discover roku on your local network",
		Action: func(c *cli.Context) {
			fmt.Println("Found roku at", findRoku())
		},
	}
}