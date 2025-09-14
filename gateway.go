package main

type GatewayService struct {
	jobSchedulerEndpoints []string // ["http://scheduler-1:3000", "http://scheduler-2:3000"]
	dockerClient          *client.Client
}
