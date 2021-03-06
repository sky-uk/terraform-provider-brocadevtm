#!groovy

import java.util.regex.*

project_name = 'terraform-provider-pulsevtm'
project_owner = 'sky-uk'

project_src_path = "github.com/${project_owner}/${project_name}"
github_token_id = 'svc-paas-github-access-token-as-text'
git_credentials_id = 'svc-paas-github-deploy-key'

version_file = 'VERSION'
major_version = null
minor_version = null
patch_version = null

docker_image = "paas/golang-img:0.10.7a"

// helpers
gitHelper = null
shellHelper = null
goHelper = null
slackHelper = null

slackChannel = '#ott-paas'

loadHelpers()


slackHelper.notificationWrapper(slackChannel, currentBuild, env, true) {
    node {
        wrap([$class: 'TimestamperBuildWrapper']) {
            wrap([$class: 'AnsiColorBuildWrapper']) {
                stage 'checkout'
                deleteDir()
                git_branch = env.BRANCH_NAME
                checkout scm
                gitHelper.prepareGit('svc-paas-github', 'svc-paas-github@jenkins.paas.int.ovp.bskyb.com')

                stage 'version'
                if (autoincVersion()) {
                    writeFile file: version_file, text: version()
                    gitHelper.commit(version_file, "bumping to: ${version()}")
                }

                echo "Starting pipeline for project: [${project_name}], branch: [${git_branch}], version: [${version()}]"

                stage 'lint'
                inContainer {
                    goHelper.goLint(project_src_path)
                }

                stage 'format'
                inContainer {
                    goHelper.goFmt(project_src_path)
                }

                stage 'vet'
                inContainer {
                    goHelper.goVet(project_src_path)
                }

                stage 'build'
                inContainer {
                    goHelper.goBuild(project_src_path)
                }

                stage 'test'
                inContainer {
                    goHelper.goTest(project_src_path)
                }

                stage 'testacc'
                // If the git branch name is prefixed with api5_1_ or api3_8_ we want to use a specific Pulse VTM server. If neither use the default.
                if(git_branch ==~ /^api3_8_.*/) {
                    pulseVTMCredentials="PULSEVTM_3_8_CREDENTIALS"
                    pulseVTMServer=env.PULSEVTM_3_8_SERVER
                    pulseVTMUnverifiedSSL=env.PULSEVTM_ALLOW_UNVERIFIED_SSL
                    pulseVTMAPI="3.8"

                } else {
                    pulseVTMCredentials = "PULSEVTM_5_1_CREDENTIALS"
                    pulseVTMServer = env.PULSEVTM_5_1_SERVER
                    pulseVTMUnverifiedSSL = env.PULSEVTM_ALLOW_UNVERIFIED_SSL
                    pulseVTMAPI = "5.1"
                }

                echo "Running acceptance tests using credentials ID: ${pulseVTMCredentials}, API version: ${pulseVTMAPI} and VTM server: ${pulseVTMServer}"

                inContainer {
                    withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: pulseVTMCredentials, usernameVariable: 'PULSEVTM_USERNAME', passwordVariable: 'PULSEVTM_PASSWORD']]) {
                        env.PULSEVTM_SERVER=pulseVTMServer
                        env.PULSEVTM_ALLOW_UNVERIFIED_SSL=pulseVTMUnverifiedSSL
                        env.PULSEVTM_API_VERSION=pulseVTMAPI
                        goHelper.goTestAcc(project_src_path)
                    }
                }

                stage 'coverage'
                inContainer {
                    goHelper.goCoverage(project_src_path)
                }

            }
        }
    }
    // we only release from master
    if (git_branch == 'master' && !gitHelper.isLastCommitFromUser('svc-paas-github')) {
        stage 'release'
        def approved = true
        timeout(time: 2, unit: 'HOURS') {
            try {
                input message: "Release ${project_name} ${version()} ?"
            } catch (InterruptedException _x) {
                echo "Releasing not approved in time!"
                approved = false
            }
        }

        if (approved) {
            echo "Release has been approved!"
            node() {
                gitHelper.tag(version(), "Jenkins ${env.JOB_NAME} ${env.BUILD_DISPLAY_NAME}")
                gitHelper.push(git_credentials_id, git_branch)

                echo "Creating GitHub Release v${version()}"

                withCredentials([string(credentialsId: github_token_id, variable: 'GITHUB_TOKEN')]) {
                    def github_release_response = gitHelper.createGitHubRelease(env.GITHUB_TOKEN, project_owner, project_name, version(), git_branch)
                    echo "${github_release_response}"
                    // FIXME: this is not working yet

                    echo "Attaching artifacts to GitHub Release v${version()}"
                    try {
                        def upload_response = gitHelper.uploadToGitHubRelease(env.GITHUB_TOKEN, project_owner, project_name, github_release_response.id, "${pwd()}/coverage.html", 'text/html')
                        echo "${upload_response}"
                    } catch (Exception e) {
                        echo "Could not upload the artifact "
                    }


                }

            }
            currentBuild.description = "Released ${version()}"
        }
    }
}

def loadHelpers() {
    fileLoader.withGit('git@github.com:sky-uk/paas-jenkins-pipelines.git', 'master', git_credentials_id, '') {
        this.gitHelper = fileLoader.load('lib/helpers/git')
        this.shellHelper = fileLoader.load('lib/helpers/shell')
        this.goHelper = fileLoader.load('lib/helpers/go')
        this.slackHelper = fileLoader.load('lib/helpers/slack')
    }
}

def autoincVersion() {
    current_version = readFile("${pwd()}/${this.version_file}").trim().tokenize(".")
    setVersion(current_version[0], current_version[1], current_version[2])

    if(gitHelper.checkIfTagExists(version())) {
        this.patch_version++
        if(gitHelper.checkIfTagExists(version())) {
            error "Next patch version (${version()}) already exists!"
        }
        return true
    }
    return false
}

def version() {
    return "${this.major_version}.${this.minor_version}.${this.patch_version}"
}

def setVersion(major, minor, patch) {
    this.major_version = major.toInteger()
    this.minor_version = minor.toInteger()
    this.patch_version = patch.toInteger()
}

def inContainer(Closure body) {
    docker.image(this.docker_image).inside("-v ${pwd()}:/gows/src/${project_src_path} -v ${System.getProperty('java.io.tmpdir')}:${System.getProperty('java.io.tmpdir')}") {
        body()
    }
}
