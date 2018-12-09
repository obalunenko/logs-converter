#!/bin/bash

git pull
#NEWVERSION=$(git tag |  sort -t. -k 1,1n -k 2,2n -k 3,3n -k 4,4n | tail -n1 | awk -F. -v OFS=. 'NF==1{print ++$NF}; NF>1{if(length($NF+1)>length($NF))$(NF-1)++; $NF=sprintf("%0*d", length($NF), ($NF+1)%(10^length($NF))); print}')

NEWVERSION=$(git tag | sed 's/\(.*v\)\([0-9]*\)\.\([0-9]*\)\.\([0-9]*\)/\2;\3;\4;\1/g' | sort  -t';' -k 1,1n  -k 2,2n -k 3,3n | tail -n 1  | awk -F';' '{printf "%s%d.%d.%d", $4, $1,$2,($3 + 1) }')

if [[ "${NEWVERSION}" = "" ]]; then
   NEWVERSION="v1.0.0"
fi

echo ${NEWVERSION}

message="version ${NEWVERSION}"
ADD="version"
read -r  -n 1 -p "y?:" userok
echo ""
if [[ "$userok" = "y" ]]; then
	read -r -n 1 -p "Update commit message?: y/n" userok
	echo ""
	if [[ "$userok" = "y" ]]; then
        read -r -p "Message: " message
        ADD="."
        echo ""
    fi
echo ${NEWVERSION} > version && git add ${ADD} && git commit -m "$message"&& git tag -a ${NEWVERSION} -m ${NEWVERSION} && git push --tags && git push
fi
echo