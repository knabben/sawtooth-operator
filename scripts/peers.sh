PODS=`kubectl get pods -l app=sawtooth -o custom-columns=NAME:.metadata.name --no-headers`
for i in $PODS
do 
    echo $i;
    kubectl exec -it $i sawtooth peer list; 
done
