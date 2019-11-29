refresh-crd:
	operator-sdk generate k8s

run:
	operator-sdk up local
