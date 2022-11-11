PWD=$(shell pwd)

.PHONY: master
master:
	docker run --rm -p 5432:5432 -d  \
		--name pg1 \
		-e POSTGRES_PASSWORD=password \
		-v $(PWD)/dbprepare.sh:/docker-entrypoint-initdb.d/init.sh \
		postgres:14

.PHONY: standby
standby:
	sudo rm -rf $(PWD)/backup/*
	docker run -p 5433:5432 -d \
		--name pg2 \
		-e POSTGRES_PASSWORD=password \
		-v $(PWD)/dbprepare.sh:/docker-entrypoint-initdb.d/init.sh \
		-v $(PWD)/standbydata:/var/lib/postgresql/data \
		-v $(PWD)/backup:/var/lib/postgresql/backup \
		postgres:14
	sleep 1
	docker exec -it pg2 bash -c "PGPASSWORD=repsecret pg_basebackup -D /var/lib/postgresql/backup -Xstream -R -v -P -h 10.7.14.40 -p 5432 -U user_replication"
	docker stop pg2
	sudo rm -rf $(PWD)/standbydata/*
	sudo cp -r $(PWD)/backup/* $(PWD)/standbydata/
	docker restart pg2
	# PGPASSWORD=repsecret pg_basebackup -D /var/lib/postgresql/data -Xstream -R -v -P -h 10.7.14.40 -p 5432 -U user_replication
	# pg_rewind -D /var/lib/postgresql/data -P -R --source-server='host=10.7.14.40 port=5432 dbname=haha user=user_replication password=repsecret'


.PHONY: promote
promote:
	docker exec -it pg2 bash -c "su postgres -c \"pg_ctl promote\""

#.PHONY: backup
#backup:
#	docker exec -it pg1 bash -c "pg_basebackup -P -R -X stream -c fast -U user_replication -D ./backup"

.PHONY: clear
clear:
	sudo rm -rf $(PWD)/backup
	sudo rm -rf $(PWD)/standbydata
	docker stop pg1 pg2 || true
	docker rm pg2
