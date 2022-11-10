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
	docker run  -p 5433:5432 -d \
		--name pg2 \
		-e POSTGRES_PASSWORD=password \
		-v $(PWD)/dbprepare.sh:/docker-entrypoint-initdb.d/init.sh \
		postgres:14

.PHONY: promote
promote:
	docker exec -it pg2 bash -c "su postgres -c \"/usr/lib/postgresql/11/bin/pg_ctl promote\""

.PHONY: backup
backup:
	docker exec -it pg1 bash -c "pg_basebackup -P -R -X stream -c fast -U user_replication -D ./backup"
	sudo ./recoveryfilesprepare.sh

.PHONY: clear
clear:
	sudo rm -rf $(PWD)/backup
	docker stop pg1 pg2 || true
