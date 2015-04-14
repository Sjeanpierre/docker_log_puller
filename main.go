package main

import (
	"archive/tar"
	"bytes"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

var client *docker.Client

type containerLog struct {
	ContainerID string
	ImageName   string
	LogName     string
	LogBody     string
}

func printContainers(containers *[]docker.APIContainers) {
	fmt.Println("Currently Running containers are:")
	for _, container := range *containers {
		fmt.Println("ID: ", container.ID)
		fmt.Println("Name: ", container.Names[0])
		fmt.Printf("Image: %v\n\n", container.Image)
	}
}

func outputStdLogs(container *docker.APIContainers) string {
	var log bytes.Buffer
	client.Logs(docker.LogsOptions{Container: container.ID, RawTerminal: true, Stdout: true, OutputStream: &log})
	return log.String()

}

func processLogFiles(containers *[]docker.APIContainers) []containerLog {
	var RetrievedLogs []containerLog
	for _, container := range *containers {
		for _, containerConfig := range config.Containers {
			if strings.HasPrefix(container.Image, containerConfig.Image) {
				logBody := outputStdLogs(&container)
				stdoutLog := containerLog{ContainerID: container.ID, ImageName: containerConfig.Image, LogName: "stdout.log", LogBody: logBody}
				RetrievedLogs = append(RetrievedLogs, stdoutLog)
				for _, logPath := range containerConfig.LogPaths {
					logContent := getFileConteants(logPath, &container)
					logPathParts := strings.Split(logPath, "/")
					logName := logPathParts[len(logPathParts)-1]
					log := containerLog{ContainerID: container.ID, ImageName: containerConfig.Image, LogName: logName, LogBody: logContent}
					RetrievedLogs = append(RetrievedLogs, log)
				}
			}
		}
	}
	return RetrievedLogs
}

func getFileConteants(filename string, container *docker.APIContainers) string {
	var buf bytes.Buffer
	err := client.CopyFromContainer(docker.CopyFromContainerOptions{
		Container:    container.ID,
		Resource:     filename,
		OutputStream: &buf,
	})
	if err != nil {
		fmt.Printf("Error while copying from %s: %s\n", container.ID, err)
	}
	content := readFile(&buf)
	return content
}

func readFile(fileBuffer *bytes.Buffer) string {
	content := new(bytes.Buffer)
	r := bytes.NewReader(fileBuffer.Bytes())
	tr := tar.NewReader(r)
	_, err := tr.Next()
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	if _, err := io.Copy(content, tr); err != nil {
		log.Fatal(err)
	}
	return content.String()
}

func writeLogsToDestination(containerLogs *[]containerLog) {
	for _, containerLog := range *containerLogs {
		fileDir := strings.Join([]string{config.LogDestination, containerLog.ImageName, containerLog.ContainerID[:10]}, "/")
		filePath := strings.Join([]string{fileDir, containerLog.LogName}, "/")
		setupLogDir(fileDir)
		fmt.Printf("%v\n", filePath)
		err := ioutil.WriteFile(filePath, []byte(containerLog.LogBody), 0644)
		if err != nil {
			log.Printf("Could not write log file to %v", filePath)
			log.Println(err)
		}

	}
}

func setup() *docker.Client {
	loadConfigFile()
	endpoint := "unix:///var/run/docker.sock"
	client, _ = docker.NewClient(endpoint)
	return client
}

func main() {
	setup()
	containers, _ := client.ListContainers(docker.ListContainersOptions{})
	printContainers(&containers)
	logs := processLogFiles(&containers)
	fmt.Printf("Writing logs to disk\n")
	writeLogsToDestination(&logs)
}
