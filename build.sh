read -p "POSTGRES_USER [acme]: " pg_user
pg_user=${pg_user:-'acme'}

read -p "POSTGRES_PASSWORD [pass]: " pg_pass
pg_pass=${pg_pass:-'pass'}

read -p "POSTGRES_DB [db]: " pg_db
pg_db=${pg_db:-'db'}

read -p "JWT_TOKEN [un1b0!!Tok3n]: " jwt_token
jwt_token=${jwt_token:-'un1b0!!Tok3n'}


export POSTGRES_USER="$pg_user"
export POSTGRES_PASSWORD="$pg_pass"
export POSTGRES_DB="$pg_db"
export DATABASE_DSN="host=airlineservice-postgres user=$pg_user password=$pg_pass dbname=$pg_db port=5432"
export JWT_TOKEN="$jwt_token"

docker build -t acmesky-airlineservice-api .

docker compose up
