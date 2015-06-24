#!/usr/bin/env bash
genqrc assets && go build -ldflags="-H windowsgui" && ./imgconverts.exe
