#!/usr/bin/env bash

#protoc --go_out=. proto/*.proto

mkdir -p agent/logic/agentpb
protoc --go_out=. agent/logic/agentpb/agent.proto

mkdir -p game/core/gamepb
protoc --go_out=. game/core/gamepb/game.proto

protoc --go_out=plugins=grpc:. pb/*.proto
