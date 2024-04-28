pipeline {
    agent any

    environment {
        DATE_TAG = "${new Date().format('yyyy-MM-dd')}"
        DOCKER_CREDENTIALS_ID = 'github-credentials-fevzi'
    }
    stages {
        stage('SCM') {
            steps {
                checkout scm
            }
        }
        stage('SonarQube Analysis') {
            steps {
                script {
                    def scannerHome = tool 'sonar-scanner';
                    withSonarQubeEnv('sonarqube-server') {
                        sh "${scannerHome}/bin/sonar-scanner -Dsonar.projectKey=SysSyncRunner"
                    }
                }
            }
        }
        stage('Build and Push Image') {
            steps {
                script {
                    //docker.withRegistry('https://ghcr.io', DOCKER_CREDENTIALS_ID) {
                        def imageTag = "ghcr.io/fevzisahinler/sysflowrunner:${env.DATE_TAG}-${env.BUILD_ID}"
                        def dockerImage = docker.build(imageTag)
                        //dockerImage.push()
                    //}
                }
            }
        }
    }
}
