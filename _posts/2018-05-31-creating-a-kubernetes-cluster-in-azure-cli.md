Creating a Kubernetes cluster in Azure

This was tested on Azure as it was in May 2018. The platform and the CLI are both moving targets, and Azure AKS is still in preview.

## Creating a Kubernetes cluster in Azure

Some variables for my own setup

    $ LOCATION=westeurope
    $ GROUP=Kubernetes
    $ CLUSTER_NAME=cluster01

Log in

    $ az account login

Select the desired account subscription id to use for this session (`az account show`)

    $ az account set --subcription $SUBSCRIPTION

Create a resource group to house the cluster (dependent resources will be created in a different group <https://github.com/Azure/AKS/issues/3>)

    $ az group create --location $LOCATION --resource-group $GROUP

Register provider resources

    $ az provider register -n Microsoft.Compute
    $ az provider register -n Microsoft.Storage
    $ az provider register -n Microsoft.Network

Create service principal credentials.

    $ az ad sp create-for-rbac --role Contributor --scopes "/subscriptions/$SUBSCRIPTION/resourceGroups/$GROUP"

Create a cluster using some of the values returned from the previous command. Use `appId` as `SP_ID` and `password` as `SP_SECRET`.

    $ az aks create --resource-group $GROUP --name $CLUSTER_NAME --node-count 2 --generate-ssh-keys --kubernetes-version 1.9.6 --node-vm-size Standard_DS1_v2 --verbose --location $LOCATION --dns-name-prefix $CLUSTER_NAME --service-principal $SP_ID --client-secret $SP_SECRET

See that the cluster has been successfully created

    $ az aks list

Browse the dashboard

    $ az aks browse --resource-group $GROUP --name $CLUSTER_NAME

## Create Ingress Controller

An Ingress Controller must be created to handle ingress resources in other namespaces.
Using [ingress-nginx](https://kubernetes.github.io/ingress-nginx/) for that.

Get credentials for `kubectl`

    $ az aks get-credentials --recource-group $GROUP --name $CLUSTER_NAME

Select the credentials in `kubectl`. If this doesn't by copy-pasting, the user may be called something else. Find out with `kubectl config view`.

    $ kubectl config set-context ${CLUSTER_NAME}-azure --cluster $CLUSTER_NAME --user=clusterUser_${GROUP}_${CLUSTER_NAME}

Use the newly created context

    $ kubectl config use-context ${CLUSTER_NAME}-azure

> **NOTE** Look at the homepage of ingress-nginx and double-check that these instructions
> are up to date. The Kubernetes eco-system as a whole is a moving target.

Set up all generic resources for nginx 

    $ kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/mandatory.yaml

Set up the outwards-facing main service resource (Azure specific step)

    $ kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/provider/cloud-generic.yaml

Verify the installation. It should show a single controller running.

    $ kubectl get pods --all-namespaces -l app=ingress-nginx --watch

Open the Kubernetes Dashboard to verify that all resources in the `nginx-ingress` namespace are running properly and not emitting errors

    $ az aks browse --resource-group $GROUP --name $CLUSTER_NAME

The ingress controller should have provisioned a single public IP address

    $ az network public-ip list --query '[*].{address: ipAddress, "Type": publicIpAllocationMethod, "Resource Group": resourceGroup, "Tags.service": tags.service}' --output table

## References
- https://docs.microsoft.com/en-us/azure/aks/kubernetes-walkthrough#open-azure-cloud-shell
- https://docs.microsoft.com/en-us/azure/aks/ingress
