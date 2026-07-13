# SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
# SPDX-License-Identifier: MIT
FROM gcr.io/distroless/static-debian12:nonroot@sha256:b7bb25d9f7c31d2bdd1982feb4dafcaf137703c7075dbe2febb41c24212b946f

ARG TARGETPLATFORM
COPY $TARGETPLATFORM/coopera /coopera

EXPOSE 8080

ENTRYPOINT ["/coopera"]
