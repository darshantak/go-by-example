import psycopg2
import uuid
shards = {
    0:{
        "host":"localhost",
        "port":5432,
        "database":"shard0",
        "user":"admin",
        "password":"admin"
    },
    1:{
        "host":"localhost",
        "port":5433,
        "database":"shard1",
        "user":"admin",
        "password":"admin"
    }
}

connections = {}

try:    
    for shard_id, params in shards.items():
        conn = psycopg2.connect(**params)
        connections[shard_id] = conn
        print("Successfully connected to shard",shard_id)
except Exception as e:
    print("Error connecting to database", e)
    exit(1)

def create_tables():
    for conn in connections.values():
        cursor = conn.cursor()
        cursor.execute("""
            CREATE TABLE IF NOT EXISTS users (
                user_id BIGINT PRIMARY KEY,
                name VARCHAR(255),
                email VARCHAR(255)
            );
        """)
        conn.commit()
    print("Tables created")

def generate_records(num_records):
    for i in range(num_records):
        user_id = i+1
        shard_id = user_id%2
        user_name = f"user_{user_id}"
        user_email = f"{user_id}@example.com"
        conn = connections[shard_id]
        cursor = conn.cursor()

        try:
            cursor.execute("INSERT INTO users(user_id,name,email) VALUES(%s,%s,%s)",(user_id,user_name,user_email))
            
        except Exception as e:
            print("Error in inserting data",e)
            exit()
    for conn in connections.values():
        conn.commit()

    print("Data generation successfull.")

if __name__ == "__main__":
    create_tables()
    generate_records(1000)
    for conn in connections.values():
        conn.close()
    
    print("Connections closed")



