pipeline {
    agent { label 'cpu1' }
    options { timestamps() }
    environment {
        CLUSTER_NAME = "${env.BRANCH_NAME.replaceAll("_", "-")}-$BUILD_NUMBER"
        NODE_COUNT = 1
        NODE_VM_SIZE = "Standard_D8s_v3"
    }
    stages {
        stage('cleanup workspace') {
            steps {
                script{
                    if (skipAll()) {
                        echo 'Aborting Build'
                        currentBuild.result = 'ABORTED'
                        throw new FlowInterruptedException(Result.ABORTED)
                    }
                    cleanWs()
                }
            }
        }
        stage('checkout') {
            steps {
                checkout scm
            }
        }
        stage('set globals') {
            steps {
                script {
                    if (env.BRANCH_NAME == "develop") {
                        def currentRC = sh(script: 'git fetch --tags && git tag -l --sort -version:refname | head -n 1 | tr "-" " " | awk  \'{print  $2}\' | tr -d rc', returnStdout: true).trim()
                        def nextRC = currentRC.toInteger() + 1
                        def nextVersion = sh(script: 'git fetch && git tag -l --sort -version:refname  | sed \'s/-.*$//g\' | sort --version-sort | tail -n1', returnStdout: true).trim()
                        NEXT_VERSION = "${nextVersion}-rc${nextRC}"
                    } else if (env.BRANCH_NAME == "main") {
                        NEXT_VERSION = sh(script: 'git fetch && git tag -l --sort -version:refname  | sed \'s/-.*$//g\' | sort --version-sort | tail -n1', returnStdout: true).trim()
                    } else {
                        NEXT_VERSION = "${env.BRANCH_NAME}-$BUILD_NUMBER"
                    }
                    echo "NEXT VERSION: ${NEXT_VERSION}"
                }
            }
        }
        stage('build linux') {
            steps {
                script {
                    sh """
                        go build -v -ldflags="-X 'main.buildVersion=${NEXT_VERSION}' -X 'main.commit=${GIT_COMMIT}'" -o cnvrgctl cmd/cnvrgctl/*.go
                       """
                }
            }
        }
        stage('setup test cluster') {
            when {
                expression { !skipTests()  }
            }
            steps {
                script {
                    withCredentials([azureServicePrincipal('jenkins-cicd-azure-new')]) {
//                        sh 'az login --service-principal -u $AZURE_CLIENT_ID -p $AZURE_CLIENT_SECRET -t $AZURE_TENANT_ID'
//                        sh 'az account set -s $AZURE_SUBSCRIPTION_ID'
//                        sh "az group create --location ${CLUSTER_LOCATION} --name ${CLUSTER_NAME}"
//                        sh "az aks create --resource-group  ${CLUSTER_NAME} --name ${CLUSTER_NAME} --location ${CLUSTER_LOCATION} --node-count ${NODE_COUNT} --node-vm-size ${NODE_VM_SIZE} --service-principal ${AZURE_CLIENT_ID} --client-secret ${AZURE_CLIENT_SECRET}"
//                        sh "az aks get-credentials --resource-group ${CLUSTER_NAME} --name ${CLUSTER_NAME} --file kubeconfig --subscription $AZURE_SUBSCRIPTION_ID"
//                        // sleep for one minute, just to make sure AKS cluster is completely ready
//                        sh "sleep 60"
//                        // deploy nginx ingress
//                        sh "KUBECONFIG=${workspace}/kubeconfig kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.40.2/deploy/static/provider/cloud/deploy.yaml"
                    }
                }
            }
        }
        stage('bump version') {
            when {
                expression { return ((env.BRANCH_NAME == "develop" || env.BRANCH_NAME == "master")) }
            }
            steps {
                script {
                    withCredentials([usernamePassword(credentialsId: '9e673d23-974c-460c-ba67-1188333cf4b4', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD')]) {
                        def url = sh(returnStdout: true, script: 'git config remote.origin.url').trim().replaceAll("https://", "")
                        sh """
                            git tag -a ${NEXT_VERSION} -m "${env.BRANCH_NAME}-${env.BUILD_NUMBER}"
                            git push https://${USERNAME}:${PASSWORD}@${url} --tags -f
                        """
                        if (env.BRANCH_NAME == "master") {
                            url = sh(returnStdout: true, script: 'git config remote.origin.url').trim().replaceAll("https://", "")
                            def nextRC = sh(script: "scripts/semver.sh bump minor ${NEXT_VERSION}", returnStdout: true).trim()
                            echo "next version gonna be: ${nextRC}-rc0"
                            sh """
                            git tag -a ${nextRC}-rc0 -m "${env.BRANCH_NAME}-${env.BUILD_NUMBER}"
                            git push https://${USERNAME}:${PASSWORD}@${url} --tags -f
                            """
                        }
                    }
                }
            }
        }
    }
    post {
        success {
            script {
                echo "Success!"
            }
        }
        failure {
            script {
                echo 'Failed!'
            }
        }
        always {
            script {
                withCredentials([azureServicePrincipal('jenkins-cicd-azure-new')]) {
                    sh 'az login --service-principal -u $AZURE_CLIENT_ID -p $AZURE_CLIENT_SECRET -t $AZURE_TENANT_ID'
                    sh 'az account set -s $AZURE_SUBSCRIPTION_ID'
                    sh """
                    if [ \$(az group list -o table  | grep ^${CLUSTER_NAME} | wc -l)  -gt 0 ]
                    then
                        echo "deleting aks cluster..."
                        az group delete --name ${CLUSTER_NAME} --no-wait -y
                    else
                        echo "cluster not found, skipping cluster delete"
                    fi 
                    """
                }
            }
        }
    }
}
