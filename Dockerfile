# SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
# SPDX-License-Identifier: MIT
FROM gcr.io/distroless/static-debian12:nonroot@sha256:f5b485ea962d9bd1186b2f6b3a061191539b905b82ec395de78cbfae51f20e35

ARG TARGETPLATFORM
COPY $TARGETPLATFORM/coopera /coopera

EXPOSE 8080

ENTRYPOINT ["/coopera"]
