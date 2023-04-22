go build . 
if [ $? -ne 0 ];then
    echo "Failed to build"
else
    echo "Builded successfully"
    echo "Add this directory to your path"
fi