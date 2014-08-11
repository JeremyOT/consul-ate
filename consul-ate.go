package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/JeremyOT/consul-ate/cmd"
	"github.com/JeremyOT/consul-ate/consul"
	"os"
	"strings"
	"time"
)

var consulAddr = flag.String("consul", "localhost:8500", "The consul agent to connect to.")
var serviceName = flag.String("name", "", "The service name to register with.")
var serviceId = flag.String("id", "", "The ID to register with.")
var ttl = flag.Duration("ttl", 30, "The health check ttl.")
var interval = flag.Duration("interval", 10, "The health check update interval.")
var check = flag.String("check", "", "JSON containing a Check. If specified, will replace the default TTL check.")
var port = flag.Int("port", 0, "The port the service operates on.")
var tags = flag.String("tags", "", "Comma delimited tags to add when registering the service.")

func main() {
	flag.Usage = func() {
		fmt.Println("Usage:")
		fmt.Println("  Registers a service with the local consul agent and tracks its health. If the command")
		fmt.Println("  exits without an error, the corresponding service will be deregistered. Otherwise, it")
		fmt.Println("  will fail its status checks.")
		fmt.Println("  To monitor an external command add the command argurments after all consulate options.")
		fmt.Println("  E.g. consulate -name my-service -- my_script.sh arg1 arg2")
		fmt.Println("  stdin, stdout and stderr will be piped through from the monitored command.")
		fmt.Println()
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}
	flag.Parse()
	if *serviceName == "" {
		flag.Usage()
	}
	id := *serviceId
	if id == "" && *port > 0 {
		id = fmt.Sprintf("%s:%d", *serviceName, *port)
		fmt.Println("ID: ", id)
	}
	tagList := strings.Split(*tags, ",")
	client := consul.NewClient(*consulAddr)
	checkDefinition := map[string]string{}
	if *check != "" {
		err := json.Unmarshal([]byte(*check), checkDefinition)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	} else {
		checkDefinition["TTL"] = fmt.Sprintf("%ds", *ttl)
	}
	serviceId, checkId, err := client.RegisterService(*serviceName, id, tagList, *port, checkDefinition)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	args := flag.Args()
	quit := make(chan int)
	var command *cmd.Command
	if len(args) == 0 {
		fmt.Println("No command specified. Running indefinitely.")
	} else {
		fmt.Println("Monitoring command:", strings.Join(args, " "))
		command = cmd.NewCommand(args)
		go command.RunCommand(quit)
	}
	client.RegisterCheckHeartbeat(checkId, "Service OK", *interval*time.Second, quit)
	if command != nil {
		if command.Error() != nil {
			client.UpdateCheck(checkId, fmt.Sprintf("Command %s exited with error: %s", command.String(), command.Error()), consul.CheckFail)
		} else {
			client.DeregisterService(serviceId)
		}
	}
}
