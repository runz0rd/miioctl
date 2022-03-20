push:
	docker build . -t pn50:30500/miioctl:latest
	docker push pn50:30500/miioctl:latest