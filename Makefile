push:
	docker build . -t pn50:30500/miioctl:slim
	docker push pn50:30500/miioctl:slim