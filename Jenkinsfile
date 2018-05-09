#!/usr/bin/groovy
@Library('github.com/stakater/fabric8-pipeline-library@master')

def dummy

toolsNode(toolsImage: 'stakater/pipeline-tools:1.5.2') {
    container(name: 'tools') {
        withCurrentRepo(type: 'go') { def repoUrl, def repoName, def repoOwner, def repoBranch ->
            String chartPackageName = ""
            String srcDir = WORKSPACE
            def kubernetesDir = WORKSPACE + "/deployments/kubernetes"

            def chartTemplatesDir = kubernetesDir + "/templates/chart"
            def chartDir = kubernetesDir + "/chart"
            def manifestsDir = kubernetesDir + "/manifests"

            def dockerContextDir = WORKSPACE + "/build/package"
            def dockerImage = repoOwner.toLowerCase() + "/" + repoName.toLowerCase()
            def dockerImageVersion = ""

            // Slack variables
            def slackChannel = "${env.SLACK_CHANNEL}"
            def slackWebHookURL = "${env.SLACK_WEBHOOK_URL}"
    
            def utils = new io.fabric8.Utils()        
            def git = new io.stakater.vc.Git()
            def helm = new io.stakater.charts.Helm()
            def templates = new io.stakater.charts.Templates()
            def common = new io.stakater.Common()
            def chartManager = new io.stakater.charts.ChartManager()
            def docker = new io.stakater.containers.Docker()
            def stakaterCommands = new io.stakater.StakaterCommands()
            def slack = new io.stakater.notifications.Slack()
            try {
                stage('Download Dependencies') {
                    sh """
                        cd ${srcDir}
                        make install
                    """ 
                }

                stage('Run Tests') {
                    sh """
                        cd ${srcDir}
                        make test
                    """
                }

                stage('Build Binary') {
                    sh """
                        cd ${srcDir}
                        go build -o ./build/package/${repoName.toLowerCase()}
                    """
                }

                if (utils.isCI()) {
                    stage('CI: Publish Dev Image') {
                        dockerImageVersion = stakaterCommands.getBranchedVersion("${env.BUILD_NUMBER}")
                        docker.buildImageWithTag(dockerContextDir, dockerImage, dockerImageVersion)
                        docker.pushTag(dockerImage, dockerImageVersion)
                    }
                    
                    stage('Notify') {
                        def dockerImageWithTag = "${dockerImage}:${dockerImageVersion}"
                        slack.sendDefaultSuccessNotification(slackWebHookURL, slackChannel, [slack.createDockerImageField(dockerImageWithTag)])

                        def commentMessage = "Image is available for testing. `docker pull ${dockerImageWithTag}`"
                        git.addCommentToPullRequest(commentMessage)
                    }
                } else if (utils.isCD()) {
                    stage('CD: Tag and Push') {
                        print "Generating New Version"
                        def versionFile = ".version"
                        def version = common.shOutput("jx-release-version --gh-owner=${repoOwner} --gh-repository=${repoName} --version-file ${versionFile}")
                        dockerImageVersion = version
                        sh """
                            echo "${version}" > ${versionFile}
                        """

                        // Render chart from templates
                        templates.renderChart(chartTemplatesDir, chartDir, repoName, version, dockerImage)
                        // Generate manifests from chart
                        templates.generateManifests(chartDir, repoName, manifestsDir)

                        git.commitChanges(WORKSPACE, "Bump Version to ${version}")

                        print "Pushing Tag ${version} to Git"
                        git.createTagAndPush(WORKSPACE, version)
                        git.createRelease(version)

                        print "Pushing Tag ${version} to DockerHub"
                        docker.buildImageWithTag(dockerContextDir, dockerImage, "latest")
                        docker.tagImage(dockerImage, "latest", version)
                        docker.pushTag(dockerImage, version)
                        docker.pushTag(dockerImage, "latest")
                    }
                    
                    stage('Chart: Init Helm') {
                        helm.init(true)
                    }

                    stage('Chart: Prepare') {
                        helm.lint(chartDir, repoName)
                        chartPackageName = helm.package(chartDir, repoName)
                    }

                    stage('Chart: Upload') {
                        String cmUsername = common.getEnvValue('CHARTMUSEUM_USERNAME')
                        String cmPassword = common.getEnvValue('CHARTMUSEUM_PASSWORD')
                        chartManager.uploadToChartMuseum(chartDir, repoName, chartPackageName, cmUsername, cmPassword)
                    }

                    stage('Notify') {
                        def dockerImageWithTag = "${dockerImage}:${dockerImageVersion}"
                        slack.sendDefaultSuccessNotification(slackWebHookURL, slackChannel, [slack.createDockerImageField(dockerImageWithTag)])

                        def commentMessage = "Image is available for testing. `docker pull ${dockerImageWithTag}`"
                        git.addCommentToPullRequest(commentMessage)
                    }
                }
            }
            catch(e) {
                slack.sendDefaultFailureNotification(slackWebHookURL, slackChannel, [slack.createErrorField(e)])
            
                def commentMessage = "Yikes! You better fix it before anyone else finds out! [Build ${env.BUILD_NUMBER}](${env.BUILD_URL}) has Failed!"
                git.addCommentToPullRequest(commentMessage)

                throw e
            }
        }
    }
}
