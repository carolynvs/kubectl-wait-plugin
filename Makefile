V?=0

deploy:
	mkdir -p ~/.kube/plugins/wait
	go build -o ~/.kube/plugins/wait/wait
	cp plugin.yaml ~/.kube/plugins/wait/
	@echo
