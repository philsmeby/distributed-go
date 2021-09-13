FROM golang
WORKDIR /go/src/app
COPY sensors .
ENTRYPOINT ./sensors --name $NAME --freq $FREQ --min $MIN --max $MAX --step $STEP