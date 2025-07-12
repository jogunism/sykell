# zsh
#color preset
DEF=`tput sgr0`
RED=`tput setaf 1`
GRN=`tput setaf 2`
YLW=`tput setaf 3`
BLU=`tput setaf 4`
MGT=`tput setaf 5`
CYN=`tput setaf 6`
WHT=`tput setaf 7`
BD=`tput bold`
DM=`tput dim`

echo ""
echo "${GRN}====================================================================="
echo ""
echo "${DEF}${BD} Project Stop "
echo "${DEF}${DM} > docker compose down --volumes --remove-orphans "
echo ""
echo "${GRN}====================================================================="
echo "${DEF}"

sleep 1;

docker compose down --volumes --remove-orphans;


echo ""
echo "${GRN}====================================================================="
echo ""
echo "${DEF}${BD} Remove all launched Docker images "
echo "${DEF}${DM} > docker rmi $(docker images -aq); "
echo ""
echo "${GRN}====================================================================="
echo "${DEF}"

docker rmi $(docker images -aq);


echo ""
echo "${YLW}${BD}ALL DONE!!"
echo "${DEF}${YLW}Thank you so much! :)"
echo "${DEF}"
