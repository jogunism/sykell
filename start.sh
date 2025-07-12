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
echo "${DEF}${BD} Project Start "
echo "${DEF}${DM} > docker compose up -d "
echo ""
echo "${GRN}====================================================================="
echo "${DEF}"

docker compose up -d;


echo ""
echo "${YLW}${BD} Frontend started via http://localhost:5173 "
echo ""
echo "${YLW}${BD} Backend starting... "
echo "${DEF}"
echo "${YLW} -- Please wait unil DB Connected --"
echo "${DM} I konw.. waiting is really annoying but, it taking a little more time. I using AWS MySQL and it takes a while to connect."
echo "${DM} You'll see the the message about DB Connection Successfull very soon." 
echo ""
echo "${DEF}"

docker logs -f sykell-backend-1

