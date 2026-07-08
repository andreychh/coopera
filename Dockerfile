# SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
# SPDX-License-Identifier: MIT
FROM gcr.io/distroless/static-debian12:nonroot@sha256:d093aa3e30dbadd3efe1310db061a14da60299baff8450a17fe0ccc514a16639

ARG TARGETPLATFORM
COPY $TARGETPLATFORM/coopera /coopera

EXPOSE 8080

ENTRYPOINT ["/coopera"]
