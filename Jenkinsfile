pipeline {
    agent none
    stages {
        stage('Build') {
            agent {
                docker { image 'obraun/vss-protoactor-jenkins' }
            }
            steps {
                sh 'cd BookingService && go build main.go'
                sh 'cd HallService && go build main.go'
                sh 'cd MovieService && go build main.go'
                sh 'cd ShowService && go build main.go'
                sh 'cd UserService && go build main.go'
            }
        }
        stage('Lint') {
            agent {
                docker { image 'obraun/vss-protoactor-jenkins' }
            }
            steps {
                //--deadline 20m --enable-all; --disable-all -E errcheck
                sh 'cd HallService && golangci-lint run --enable-all --skip-dirs proto -D wsl -D lll -D goimports -D golint -D stylecheck'
                sh 'cd MovieService && golangci-lint run --enable-all --skip-dirs proto -D wsl -D lll -D gosimple -D goimports -D golint -D stylecheck'
                sh 'cd BookingService && golangci-lint run --enable-all --skip-dirs proto -D wsl -D lll -D funlen -D unparam -D goimports -D golint -D stylecheck'
                sh 'cd ShowService && golangci-lint run --enable-all --skip-dirs proto -D wsl -D lll -D golint -D goimports -D golint -D stylecheck'
                sh 'cd UserService && golangci-lint run --enable-all --skip-dirs proto -D wsl -D lll'
            }
        }
        stage('Build Docker Image') {
            agent any
            steps {
                sh "docker-build-and-push -b ${BRANCH_NAME} -s HallService -f HallService/dockerfile ."
                sh "docker-build-and-push -b ${BRANCH_NAME} -s MovieService -f MovieService/dockerfile ."
                sh "docker-build-and-push -b ${BRANCH_NAME} -s BookingService -f BookingService/dockerfile ."
                sh "docker-build-and-push -b ${BRANCH_NAME} -s ShowService -f ShowService/dockerfile ."
                sh "docker-build-and-push -b ${BRANCH_NAME} -s UserService -f UserService/dockerfile ."
            }
        }
    }
}