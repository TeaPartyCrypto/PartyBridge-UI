image: 
	@make bld
	@cd be/cmd && go mod tidy && GOOS=linux go build -o be main.go
	@cd be && gcloud builds submit --tag gcr.io/mineonlium/partybridgeui . 

bld:
	@cd be/cmd && rm -rf kodata  && cd ../ && yarn && yarn build
	@mv build kodata && mv kodata be/cmd
	@cd be/cmd && go mod tidy && go build -o be main.go

release:
	@make bld
	@make image