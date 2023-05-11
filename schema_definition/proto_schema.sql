-- Protobuf -> ZetaSQL

create or replace table sample_table (
  double_field FLOAT64,
  float_field FLOAT64,
  int32_field  INT64,
  int64_field INT64,
  uint32_field INT64,
  sint32_field INT64,
  sint64_field INT64,
  fixed32_field INT64,
  sfixed32_field INT64,
  sfixed64_field INT64,
  bool_field BOOLEAN,
  string_field STRING,
  bytes_field BYTES,
  enum_field INT64,
  message_field STRUCT<hoge STRING>,
  map_field ARRAY<STRUCT<key STRING, value INTEGER>>,
  repeated_int32_field ARRAY<INT64>,
  optional_int32_field INT64,
)

