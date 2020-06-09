package main

import (
    "flag"
    "fmt"
    "github.com/spacetimi/timi_shared_server/scripts/scripting_utilities"
    "github.com/spacetimi/timi_shared_server/utils/go_vars_helper"
    "os"
    "os/exec"
    "sync"
)

func usage() {
    fmt.Println("!! Usage: timi_build -app=APP_NAME -env=ENVIRONMENT -appdir=<path to your app> [-v] [-run]")
    flag.PrintDefaults()
}

func main() {

    appPtr          := flag.String("app", "", "Name of a valid spacetimi app")
    appDirPtr       := flag.String("appdir", "", "Path to your app. This is the path to the app's directory in GOPATH/src/.../<your_app_name>")
    envPtr          := flag.String("env", "", "Local, Test, Staging, Production")
    verbosePtr      := flag.Bool("v", false, "Verbose output from this build tool")
    runPtr          := flag.Bool("run", false, "Run after building. If absent, build only")

    flag.Usage = usage
    flag.Parse()

    appName   := *appPtr
    appDir    := *appDirPtr
    appEnv    := *envPtr
    verbose   := *verbosePtr
    shouldRun := *runPtr

    /** Validate parameters or die **/
    if len(*appPtr) == 0 ||
       len(*envPtr) == 0 ||
       (appName != "bonda" && appName != "passman") ||
       len(appDir) == 0 ||
       (appEnv != "Local" && appEnv != "Test" && appEnv != "Staging" && appEnv != "Production"){

        flag.Usage()
        os.Exit(1)
    }

    var waitGroup sync.WaitGroup
    waitGroup.Add(1)
    go func() {
        defer waitGroup.Done()
        err := build_and_start_local_server(appDir, appName, appEnv, verbose, shouldRun)
        if err != nil {
            fmt.Println("Build failed|error=" + err.Error())
            os.Exit(1)
        }
    }()
    waitGroup.Wait()
}


func build_and_start_local_server(appDirPath string, appName string, appEnv string, verbose bool, shouldRunAfterBuilding bool) error {

	appDir, err := os.Stat(appDirPath)
	if err != nil {
	    return scripting_utilities.ScriptError{err.Error()}
    }
	if !appDir.IsDir() {
	    return scripting_utilities.ScriptError{appDirPath + " is not a directory"}
    }

	outputFolderPath := go_vars_helper.GOPATH + "/bin/" + appName
    outputFilePath := outputFolderPath + "/" + appName + "-server"

    if verbose {
        fmt.Println("Building executable from package: " + appDirPath + "/main")
        fmt.Println("Output path: " + outputFilePath)
    }

    buildCommand := exec.Command(go_vars_helper.GOROOT + "/bin/go", "build", "-o", outputFilePath, "main/main.go")
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
        runCommand.Stdout = os.Stdout
        runCommand.Stderr = os.Stderr
        err = runCommand.Run()
        if err != nil {
        	return scripting_utilities.ScriptError{ "Run command failed with: " + err.Error()}
        }
    }

    return nil
}
