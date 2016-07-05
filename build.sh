#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

echo "Building Docker image microservice2"
docker build -t microservice2 -f Dockerfile .

echo "Generating ./run.sh"
rm -f ./run.sh
echo "#!/bin/bash" >> ./run.sh
echo "docker run --net=codefresh_test --rm --name ms2 -it microservice2" >> ./run.sh
chmod +x ./run.sh

echo "You can now run ./run.sh to start microservice2."



