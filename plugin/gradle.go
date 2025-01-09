package plugin

import (
	"fmt"
	"log"
)

var GradleConfigJsonTagToExeFlagMapStringItemList = []JsonTagToExeFlagMapStringItem{
	{"--deploy-ivy-desc=", "PLUGIN_DEPLOY_IVY_DESC", false, false},
	{"--deploy-maven-desc=", "PLUGIN_DEPLOY_MAVEN_DESC", false, false},
	{"--global=", "PLUGIN_GLOBAL", false, false},
	{"--ivy-artifacts-pattern=", "PLUGIN_IVY_ARTIFACTS_PATTERN", false, false},
	{"--ivy-desc-pattern=", "PLUGIN_IVY_DESC_PATTERN", false, false},
	{"--repo-deploy=", "PLUGIN_REPO_DEPLOY", false, false},
	{"--repo-resolve=", "PLUGIN_REPO_RESOLVE", false, false},
	{"--server-id-deploy=", "PLUGIN_SERVER_ID_DEPLOY", false, false},
	{"--server-id-resolve=", "PLUGIN_SERVER_ID_RESOLVE", false, false},
	{"--use-wrapper=", "PLUGIN_USE_WRAPPER", false, false},
	{"--uses-plugin=", "PLUGIN_USES_PLUGIN", false, false},
}

var GradleRunJsonTagToExeFlagMapStringItemList = []JsonTagToExeFlagMapStringItem{
	{"--build-name=", "PLUGIN_BUILD_NAME", false, false},
	{"--build-number=", "PLUGIN_BUILD_NUMBER", false, false},
	{"--detailed-summary=", "PLUGIN_DETAILED_SUMMARY", false, false},
	{"--format=", "PLUGIN_FORMAT", false, false},
	{"--project=", "PLUGIN_PROJECT", false, false},
	{"--scan=", "PLUGIN_SCAN", false, false},
	{"--threads=", "PLUGIN_THREADS", false, false},
}

func GetGradleCommandArgs(args Args) ([][]string, error) {

	var cmdList [][]string

	jfrogConfigAddConfigCommandArgs, err := GetConfigAddConfigCommandArgs(args.ResolverId,
		args.Username, args.Password, args.URL, args.AccessToken, args.APIKey)
	if err != nil {
		return cmdList, err
	}

	gradleConfigCommandArgs := []string{GradleConfig}
	err = PopulateArgs(&gradleConfigCommandArgs, &args, GradleConfigJsonTagToExeFlagMapStringItemList)
	if err != nil {
		return cmdList, err
	}

	gradleTaskCommandArgs := []string{GradleCmd, args.GradleTasks}
	err = PopulateArgs(&gradleTaskCommandArgs, &args, GradleRunJsonTagToExeFlagMapStringItemList)
	if err != nil {
		return cmdList, err
	}

	if len(args.BuildFile) > 0 {
		gradleTaskCommandArgs = append(gradleTaskCommandArgs, "-b "+args.BuildFile)
	}

	cmdList = append(cmdList, jfrogConfigAddConfigCommandArgs)
	cmdList = append(cmdList, gradleConfigCommandArgs)
	cmdList = append(cmdList, gradleTaskCommandArgs)

	return cmdList, nil
}

func GetGradlePublishCommand(args Args) ([][]string, error) {

	fmt.Println("GetGradlePublishCommand")

	var cmdList [][]string
	var jfrogConfigAddConfigCommandArgs []string

	tmpServerId := args.DeployerId // "tmpSrvConfig"
	jfrogConfigAddConfigCommandArgs, err := GetConfigAddConfigCommandArgs(tmpServerId,
		args.Username, args.Password, args.URL, args.AccessToken, args.APIKey)
	if err != nil {
		log.Println("GetConfigAddConfigCommandArgs error: ", err)
		return cmdList, err
	}

	gradleConfigCommandArgs := []string{GradleConfig}
	err = PopulateArgs(&gradleConfigCommandArgs, &args, MavenConfigCmdJsonTagToExeFlagMapStringItemList)
	if err != nil {
		log.Println("PopulateArgs error: ", err)
		return cmdList, err
	}
	gradleConfigCommandArgs = append(gradleConfigCommandArgs, "--server-id-deploy="+tmpServerId)
	gradleConfigCommandArgs = append(gradleConfigCommandArgs, "--server-id-resolve="+tmpServerId)

	rtPublishCommandArgs := []string{"gradle", Publish}
	switch {
	case args.Username != "":
		rtPublishCommandArgs = append(rtPublishCommandArgs, "-Pusername="+args.Username)
		rtPublishCommandArgs = append(rtPublishCommandArgs, "-Ppassword="+args.Password)
	case args.AccessToken != "":
		rtPublishCommandArgs = append(rtPublishCommandArgs, "-PaccessToken="+args.AccessToken)
	}
	rtPublishCommandArgs = append(rtPublishCommandArgs, "--build-name="+args.BuildName)
	rtPublishCommandArgs = append(rtPublishCommandArgs, "--build-number="+args.BuildNumber)

	rtPublishBuildInfoCommandArgs := []string{"rt", BuildPublish, args.BuildName, args.BuildNumber,
		"--server-id=" + tmpServerId}
	err = PopulateArgs(&rtPublishBuildInfoCommandArgs, &args, RtBuildInfoPublishCmdJsonTagToExeFlagMap)
	if err != nil {
		log.Println("PopulateArgs error: ", err)
		return cmdList, err
	}

	cmdList = append(cmdList, jfrogConfigAddConfigCommandArgs)
	cmdList = append(cmdList, gradleConfigCommandArgs)
	cmdList = append(cmdList, rtPublishCommandArgs)
	cmdList = append(cmdList, rtPublishBuildInfoCommandArgs)

	for _, cmd := range cmdList {
		fmt.Println("#################################")
		fmt.Println(cmd)
	}
	fmt.Println("")
	return cmdList, nil
}
