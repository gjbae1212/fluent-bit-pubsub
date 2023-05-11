-- Avro -> ZetaSQL

create or replace table sample_table (
  BooleanField  BOOLEAN,
  IntField INTEGER,
  LongField INTEGER,
  FloatField FLOAT64,
  DoubleField FLOAT64,
  BytesField BYTES,
  StringField STRING,
  RecordField STRUCT<one INTEGER, two FLOAT64, three STRING>,
  ArrayField ARRAY<string>,
  MapField ARRAY<STRUCT<myMapType FLOAT64>>,
  NullUnionField STRING,
  FixedField BYTES,
  EnumField INT64,
  TimestampField TIMESTAMP,
  DateField DATE,
  TimeField TIME,
)
