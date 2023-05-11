# Sample Custom Schema

Sample definitions to check custom schema for Pub/Sub.  
Data type compativility: https://cloud.google.com/pubsub/docs/bigquery#avro-to-zetasql

## Avro

- avro_schema.sql: Table definition to check BigQuery subscription in Avro.
- test1.avsc: To check data type mapping Avro to ZetaSQL.
  - Avro data definition: https://avro.apache.org/docs/1.11.1/specification/_print/
- test2.avsc: To check custom schema revision.
  - any field must be NULLABLE to change add/remove by revision.
  - Avro needs default field

## ProtocolBuffer

- proto_schema.sql: Table definition to check BigQuery subscription in ProtocolBuffer.
- test.proto: To check data type mapping & custom schema revision.

