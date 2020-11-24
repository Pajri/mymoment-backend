deploy:
	git pull
	go build -o mymoment
	sudo systemctl restart mymoment.service

ensure:
	dep ensure -v