pipeline{
    agent{
         node {
            label 'master'
            customWorkspace "workspace/${env.BRANCH_NAME}/src/github.com/aditya37/geofence-service/"
        }
    }
    environment {
        SERVICE  = "geofence-service"
        NOTIFDEPLOY = -522638644
    }
    options {
        buildDiscarder(logRotator(daysToKeepStr: env.BRANCH_NAME == 'main' ? '90' : '30'))
    }
    stages{
        stage("Checkout"){
            when {
                anyOf { branch 'main'; branch 'develop'; branch 'staging' }
            }
            // Do clone
            steps {
                echo 'Checking out from git'
                checkout scm
                script {
                    env.GIT_COMMIT_MSG = sh (script: 'git log -1 --pretty=%B ${GIT_COMMIT}', returnStdout:true).trim()
                }
            }
        }
        stage('Build and deploy') {
            environment {
                GOPATH = "${env.JENKINS_HOME}/workspace/${env.BRANCH_NAME}"
                PATH = "${env.GOPATH}/bin:${env.PATH}"
            }
            stages {
                // build to dev
                stage('Deploy to env development') {
                    when {
                        branch 'develop'
                    }
                    environment {
                        NAMESPACE = 'core-development'
			            TAG= '0.0.1'
                    }
                    steps {
                        // get credential file
                        withCredentials([file(credentialsId: 'b2505411-9b44-4bd1-82f8-8c52e938de07', variable: 'config')]) {
                            echo 'Build image'
                            sh "cp $config .env.geofence"
                            sh "chmod 644 .env.geofence"
			                sh 'chmod +x build.sh'
			                sh './build.sh'
                            sh 'chmod +x deploy.sh'
                            sh './deploy.sh'
                            sh 'rm .env.geofence'
                        }
                    }
                }
            }
        }
    }
    post{
        success{
            telegramSend(message:"Application $SERVICE has been [deployed] With Commit Message $GIT_COMMIT_MSG",chatId:"$NOTIFDEPLOY")
        }
        failure{
            telegramSend(message:"Application $SERVICE has been [Failed] With Commit Message $GIT_COMMIT_MSG",chatId:"$NOTIFDEPLOY")
        }
    }
}
