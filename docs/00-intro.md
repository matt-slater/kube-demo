# Intro

You are a commander in charge of the Rebel Fleet. Your mission is to destroy the Death Star using the power of Kubernetes!

## Mission Prerequisites

- a devshell environment with `klogin` and `kubectl`
- a text editor

## Determine from where you will launch your attack

- Use go/gkpdashboard to find a suitable development cluster
- Log in to the cluster with klogin
  - `klogin -a blah -k`
- Create a namespace in the cluster using your SEAL ID
  - `kubectl create namespace ${SEAL_ID}-rebel-base-dev`
- Since you will be working out of this namespace set your kube context to that namespace
  - `kubectl config set-context --namespace ${YOUR_NAMESPACE_NAME}`

Congratulations! you have successfully set up your environment.

## The First Attack

### Pods

look at the `pod.yaml` file in `resources/00-intro`. A pod is the Kubernetes construct that allows you to run containers. 

This file tells Kubernetes how to run your containers. Customize the values of the `COMMANDER` and `TARGET` fields.

Early recon missions have revealed the location of the Death Star, which will be your target:

`death-star.apps.mt-d1.na-mw-s01.gkp.jpmchase.net`

When you are ready deploy your X Wing pod to the cluster:

`kubectl apply -f resources/00-intro/pod.yaml`

You can view information about your pods like so:

`kubectl get pods`

You will see your pod has failed. Let's investigate. To see the logs for your pod, execute the following:

`kubectl logs x-wing`

There seems to be a problem with our X wing being able to reach the Death Star. This is because of GKP's deny-all security posture. Space can be a dangerous place! Your pod is timing out trying to reach the Death Star. We need to create a NetworkPolicy to allow this communication.

### NetworkPolicy

Look at the `resources/00-intro/networkpolicy.yaml` and edit it to allow access to the Death Star's DNS name. Then apply the file to the cluster:

`kubectl apply -f resources/00-intro/networkpolicy.yaml`

If you run `kubectl get pods` again, you'll see that the pod is stuck in a failed state. If you delete the pod, recreate it and view the logs we should see that the attack was successful.

```shell
kubectl delete pod x-wing
kubectl apply -f resources/00-intro/pod.yaml
kubectl logs x-wing
```




