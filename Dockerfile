# SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
# SPDX-License-Identifier: MIT
FROM gcr.io/distroless/static-debian12:nonroot

ARG TARGETPLATFORM
COPY $TARGETPLATFORM/coopera /coopera

EXPOSE 8080

ENTRYPOINT ["/coopera"]
