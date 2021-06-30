#!/bin/env sh
protoc --go_out=../ --go_opt=paths=source_relative --go-grpc_out=../ --go-grpc_opt=paths=source_relative *.proto
protoc --doc_out=resource/custom_markdown.tpl,index.md:./ *.proto