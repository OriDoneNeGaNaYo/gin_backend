import datetime
import os

import pyspark.sql
from pyspark.sql import SparkSession


def main():
    spark = SparkSession.builder \
        .appName("bus-stop_preprocessing") \
        .config("spark.jars", f"{os.getcwd()}/mysql-connector-j-8.0.32/mysql-connector-j-8.0.32.jar") \
        .config("spark.jars.packages", "org.elasticsearch:elasticsearch-hadoop:8.6.2") \
        .getOrCreate()

    directory = f"{os.getcwd()}/analysis/data"

    df = csv_to_df(directory, spark)
    df.show()
    elastic_upsert(df)


def csv_to_df(directory: str, spark: SparkSession) -> pyspark.sql.DataFrame:
    sg = "gyoung.csv"
    bus = "bus.csv"

    sg_data = spark.read.csv(f"{directory}/{sg}", inferSchema=True, header=True)
    bus_data = spark.read.csv(f"{directory}/{bus}", inferSchema=True, header=True)
    sg_data.createOrReplaceTempView("sg_data")
    bus_data.createOrReplaceTempView("bus_data")

    qs = """
    SELECT 
        `정류장아이디` AS id,
        sg_data.`도시명` AS city,
        `정류장 명칭` AS name,
        `위도` AS gps_lati,
        `경도` AS gps_long,
        `수집일시` AS collected_at
    FROM
        bus_data
        JOIN
            sg_data
        ON
            bus_data.`도시명` = sg_data.`도시명`
    """
    return spark.sql(qs)


def mysql_upsert(dataframe: pyspark.sql.DataFrame):
    # Upsert 작업 수행
    upsert_query = """
        INSERT IGNORE INTO bus_stops
        (id, city, name, gps_lati, gps_long, collected_at) -- MySQL 테이블 컬럼 지정
        VALUES (?, ?, ?, ?, ?, ?) -- VALUES 절 지정
        ON DUPLICATE KEY UPDATE
        city = VALUES(city), -- 중복 시 업데이트될 컬럼 지정
        name = VALUES(name),
        gps_lati = VALUES(gps_lati),
        gps_long = VALUES(gps_long),
        collected_at = VALUES(collected_at),
        updated_at = NOW()
    """

    # Write to MySQL Table
    dataframe.write \
        .format("jdbc") \
        .option("driver", "com.mysql.jdbc.Driver") \
        .option("url", "jdbc:mysql://localhost:3306/bus") \
        .option("dbtable", "bus_stops") \
        .option("user", "root") \
        .option("password", "busdb") \
        .option("insertStatement", upsert_query) \
        .mode("overwrite") \
        .save()


def mysql_read(spark: SparkSession) -> pyspark.sql.DataFrame:
    # Read from MySQL Table
    return spark.read \
        .format("jdbc") \
        .option("driver", "com.mysql.jdbc.Driver") \
        .option("url", "jdbc:mysql://localhost:3306/bus") \
        .option("dbtable", "bus_stops") \
        .option("user", "root") \
        .option("password", "busdb") \
        .load()


def elastic_upsert(dataframe: pyspark.sql.DataFrame):
    es_conf = {
        "es.nodes": "localhost",  # Elasticsearch 호스트
        "es.port": "9200",  # Elasticsearch 포트
        "es.resource": "bus_stops/_doc",  # 저장할 인덱스와 타입
        "es.input.json": "yes",  # JSON 형태로 저장
        "es.write.operation": "index",  # 저장 방법 (create, index, update)
        "es.nodes.wan.only": "true"  # 외부 노드와의 통신을 위한 설정
    }

    # elasticsearch upsert 작업 수행
    dataframe.write.format("org.elasticsearch.spark.sql") \
        .options(**es_conf) \
        .mode("overwrite") \
        .save()


main()
