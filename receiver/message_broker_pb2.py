# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: message_broker.proto
# Protobuf Python Version: 5.27.2
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    5,
    27,
    2,
    '',
    'message_broker.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x14message_broker.proto\x12\x06\x62roker\"\x1a\n\x07Message\x12\x0f\n\x07\x63ontent\x18\x01 \x01(\t\"S\n\x0fMessageMetadata\x12 \n\x07message\x18\x01 \x01(\x0b\x32\x0f.broker.Message\x12\x0f\n\x07\x63ommand\x18\x02 \x01(\t\x12\r\n\x05topic\x18\x03 \x01(\t\"\x1d\n\x0cTopicRequest\x12\r\n\x05topic\x18\x01 \x01(\t\",\n\x08Response\x12\x0f\n\x07success\x18\x01 \x01(\x08\x12\x0f\n\x07message\x18\x02 \x01(\t2\xb2\x01\n\rMessageBroker\x12\x34\n\x07Publish\x12\x17.broker.MessageMetadata\x1a\x10.broker.Response\x12\x34\n\tSubscribe\x12\x14.broker.TopicRequest\x1a\x0f.broker.Message0\x01\x12\x35\n\x0bUnsubscribe\x12\x14.broker.TopicRequest\x1a\x10.broker.ResponseB\x0fZ\r/proto;brokerb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'message_broker_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z\r/proto;broker'
  _globals['_MESSAGE']._serialized_start=32
  _globals['_MESSAGE']._serialized_end=58
  _globals['_MESSAGEMETADATA']._serialized_start=60
  _globals['_MESSAGEMETADATA']._serialized_end=143
  _globals['_TOPICREQUEST']._serialized_start=145
  _globals['_TOPICREQUEST']._serialized_end=174
  _globals['_RESPONSE']._serialized_start=176
  _globals['_RESPONSE']._serialized_end=220
  _globals['_MESSAGEBROKER']._serialized_start=223
  _globals['_MESSAGEBROKER']._serialized_end=401
# @@protoc_insertion_point(module_scope)
