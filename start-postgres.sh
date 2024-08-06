docker run --name go-geo \
  -e POSTGRES_USER=$(grep DB_USER .env | cut -d '=' -f2) \
  -e POSTGRES_PASSWORD=$(grep DB_PASSWORD .env | cut -d '=' -f2) \
  -e POSTGRES_DB=$(grep DB_NAME .env | cut -d '=' -f2) \
  -p $(grep DB_PORT .env | cut -d '=' -f2):5432 \
  -d postgis/postgis:13-3.3