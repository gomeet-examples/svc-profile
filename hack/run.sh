#!/bin/sh

trap killgroup 2

SCRIPT=$(readlink -f "$0")
SCRIPTPATH=$(dirname "$SCRIPT")

GOARCH=$(go env GOARCH)
GOOS=$(go env GOOS)

[ -z $GOMEET_EXEC_TYPE ] && GOMEET_EXEC_TYPE="make" # go, make

[ -z $GOMEET_JWT_SECRET ] && GOMEET_JWT_SECRET=""

if [ -z $GOMEET_SVC_PROFILE_DB_DSN ]
then
	[ -z $GOMEET_SVC_PROFILE_DB_USERNAME ] && GOMEET_SVC_PROFILE_DB_USERNAME="gomeet"
	[ -z $GOMEET_SVC_PROFILE_DB_PASSWORD ] && GOMEET_SVC_PROFILE_DB_PASSWORD="totomysql"
	[ -z $GOMEET_SVC_PROFILE_DB_SERVER ] && GOMEET_SVC_PROFILE_DB_SERVER="localhost"
	[ -z $GOMEET_SVC_PROFILE_DB_PORT ] && GOMEET_SVC_PROFILE_DB_PORT="3306"
	[ -z $GOMEET_SVC_PROFILE_DB_DATABASE ] && GOMEET_SVC_PROFILE_DB_DATABASE="svc_profile"
	GOMEET_SVC_PROFILE_DB_DSN="$GOMEET_SVC_PROFILE_DB_USERNAME:$GOMEET_SVC_PROFILE_DB_PASSWORD@tcp($GOMEET_SVC_PROFILE_DB_SERVER:$GOMEET_SVC_PROFILE_DB_PORT)/$GOMEET_SVC_PROFILE_DB_DATABASE"
fi;

killgroup(){
	echo killing...
	kill 0
}

# SUB-SERVICES DEFINITION : run.sh
# END SUB-SERVICES DEFINITION : run.sh

cd $SCRIPTPATH/../
SERVER_OPTS='serve
				-d
				--jwt-secret "$GOMEET_JWT_SECRET"
				--mysql-migrate
				--mysql-dsn "$GOMEET_SVC_PROFILE_DB_DSN"
				--address ":13000"'

case $GOMEET_EXEC_TYPE in
	"go")
		CMD='CGO_ENABLED=0 go run
			-ldflags "-extldflags \"-lm -lstdc++ -static\""
			-ldflags "-X github.com/gomeet-examples/svc-profile/service.version=$(cat VERSION) -X github.com/gomeet/svc-profile/service.name=svc-profile"
			main.go'
		eval $CMD $SERVER_OPTS
		break
		;;
	"make")
		make
		eval _build/svc-profile $SERVER_OPTS
		break
		;;
	*)
		echo "Error : unknow $GOMEET_EXEC_TYPE value for GOMEET_EXEC_TYPE [go|make] allowed"
		;;
esac

wait

