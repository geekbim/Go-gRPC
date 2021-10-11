docker run --name postgresql-container -p 5432:5432 -e POSTGRES_PASSWORD=hwhwhwlol -d postgres
docker exec -it postgresql-container psql -U postgres