mkdir autogen
swagger generate server -f ./petstore-minimal.yaml -t ./autogen -P models.Principal
swagger generate client -f ./petstore-minimal.yaml -t ./autogen --skip-models

