docker build . -t survey

docker run --rm  -p 8081:8080 -p 6061:6060 --name survey-service2  --network mongoCluster survey