package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/spacetimi/timi_shared_server/scripts/scripting_utilities"
	"github.com/spacetimi/timi_shared_server/utils/go_vars_helper"
)

func usage() {
	fmt.Println("!! Usage: timi_build -app=APP_NAME -env=ENVIRONMENT -appdir=<path to your app's code> -shareddir=<path to shared code> [-awsprofile=<aws profile to use for aws-sessions>] [-v] [-run]")
	flag.PrintDefaults()
}

func main() {

	appPtr := flag.String("app", "", "Name of a valid spacetimi app")
	appDirPtr := flag.String("appdir", "", "Path to your app's code.")
	sharedDirPtr := flag.String("shareddir", "", "Path to shared code.")
	envPtr := flag.String("env", "", "Local, Test, Staging, Production")
	awsProfilePtr := flag.String("awsprofile", "", "Optional AWS profile to use for creating AWS-sessions in the AWS sdk")
	verbosePtr := flag.Bool("v", false, "Verbose output from this build tool")
	runPtr := flag.Bool("run", false, "Run after building. If absent, build only")

	flag.Usage = usage
	flag.Parse()

	appName := *appPtr
	appDir := *appDirPtr
	sharedDir := *sharedDirPtr
	appEnv := *envPtr
	awsProfile := *awsProfilePtr
	verbose := *verbosePtr
	shouldRun := *runPtr

	/** Validate parameters or die **/
	if len(*appPtr) == 0 ||
		len(*envPtr) == 0 ||
		(appName != "bonda" && appName != "passman_server" && appName != "pfh_reader_server") ||
		len(appDir) == 0 ||
		len(sharedDir) == 0 ||
		(appEnv != "Local" && appEnv != "Test" && appEnv != "Staging" && appEnv != "Production") {

		flag.Usage()
		os.Exit(1)
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		err := build_and_start_local_server(appDir, sharedDir, appName, appEnv, awsProfile, verbose, shouldRun)
		if err != nil {
			fmt.Println("Build failed|error=" + err.Error())
			os.Exit(1)
		}
	}()
	waitGroup.Wait()
}

func build_and_start_local_server(appDirPath string, sharedDirPath string, appName string, appEnv string, awsProfile string, verbose bool, shouldRunAfterBuilding bool) error {

	appDir, err := os.Stat(appDirPath)
	if err != nil {
		return scripting_utilities.ScriptError{err.Error()}
	}
	if !appDir.IsDir() {
		return scripting_utilities.ScriptError{appDirPath + " is not a directory"}
	}

	sharedDir, err := os.Stat(sharedDirPath)
	if err != nil {
		return scripting_utilities.ScriptError{err.Error()}
	}
	if !sharedDir.IsDir() {
		return scripting_utilities.ScriptError{sharedDirPath + " is not a directory"}
	}

	outputFolderPath := go_vars_helper.GOPATH + "/bin/" + appName
	outputFilePath := outputFolderPath + "/" + appName + "-server"

	if verbose {
		fmt.Println("Building executable from package: " + appDirPath + "/main")
		fmt.Println("Output path: " + outputFilePath)
	}

	buildCommand := exec.Command(go_vars_helper.GOROOT+"/bin/go", "build", "-o", outputFilePath, "main/main.go")
	buildCommand.Dir = appDirPath
	buildCommand.Stdout = os.Stdout
	buildCommand.Stderr = os.Stderr
	err = buildCommand.Run()
	if err != nil {
		return scripting_utilities.ScriptError{"Build command failed with: " + err.Error()}
	}

	if shouldRunAfterBuilding {
		runCommand := exec.Command(outputFilePath)
		runCommand.Env = os.Environ()
		runCommand.Env = append(runCommand.Env, "app_environment="+appEnv)
		runCommand.Env = append(runCommand.Env, "app_name="+appName)
		runCommand.Env = append(runCommand.Env, "app_dir_path="+appDirPath)
		runCommand.Env = append(runCommand.Env, "shared_dir_path="+sharedDirPath)
		runCommand.Env = append(runCommand.Env, "aws_profile="+awsProfile)
		runCommand.Stdout = os.Stdout
		runCommand.Stderr = os.Stderr
		err = runCommand.Run()
		if err != nil {
			return scripting_utilities.ScriptError{"Run command failed with: " + err.Error()}
		}
	}

	return nil
}
