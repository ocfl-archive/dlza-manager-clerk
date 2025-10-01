FROM golang:1.25.1 as builder

WORKDIR /dlza-manager-clerk
ARG SSH_PUBLIC_KEY=$SSH_PUBLIC_KEY
ARG SSH_PRIVATE_KEY=$SSH_PRIVATE_KEY

ARG GITLAB_USER=gitlab-ci-token
ARG GITLAB_PASS=$CI_JOB_TOKEN
# ARG SSH_PRIVATE_KEY
# ARG SSH_PUBLIC_KEY

ENV GO111MODULE=on
ENV GOPRIVATE=gitlab.switch.ch/ub-unibas/*
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

COPY . .


# RUN cat go.mod
RUN apt-get update && \
    apt-get install -y \
        git \
        openssh-client \
        ca-certificates \
        protobuf-compiler \
        nodejs \
        npm
RUN npm install -g npm@10.9.1
RUN npm install -g node@22.9.0
# RUN apk add --no-cache ca-certificates git openssh-client
# RUN 'which ssh-agent || ( apt-get update -y && apt-get install openssh-client git -y )'
RUN eval $(ssh-agent -s)
RUN mkdir -p ~/.ssh
RUN chmod 700 ~/.ssh
# #for CI/CD build
RUN echo "$SSH_PRIVATE_KEY" | base64 -d >> ~/.ssh/id_rsa
# #for local build
# RUN echo "$SSH_PRIVATE_KEY" >> ~/.ssh/id_rsa
RUN echo "$SSH_PUBLIC_KEY" | tr -d '\r'   >> ~/.ssh/authorized_keys
# # set chmod 600 else bas permission it fails
RUN chmod 600 ~/.ssh/id_rsa
RUN chmod 644 ~/.ssh/authorized_keys
RUN ssh-keyscan gitlab.switch.ch >> ~/.ssh/known_hosts
RUN chmod 644 ~/.ssh/known_hosts
# RUN git config --global url."ssh://git@gitlab.switch.ch/".insteadOf "https://gitlab.switch.ch/"
RUN git config --global --add url."https://gitlab-ci-token:${CI_JOB_TOKEN}@gitlab.switch.ch".insteadOf "https://gitlab.switch.ch"
# RUN ssh -A -v -l git gitlab.switch.ch

# with DOCKER_BUILDKIT=1 for ssh
# RUN --mount=type=ssh go mod download
RUN go mod download
# RUN git clone https://${GITLAB_USER}:${GITLAB_PASS}@gitlab.switch.ch/ub-unibas/dlza/microservices/pbtypes /pbtypes
# RUN go get google.golang.org/protobuf/protoc-gen-go
# RUN go get google.golang.org/protobuf
# RUN go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
# RUN bash ./build.sh


# RUN git clone -b develop https://gitlab.switch.ch/ub-unibas/dlza/dlza-frontend.git
RUN git clone -b develop git@gitlab.switch.ch:ub-unibas/dlza/dlza-frontend.git
## to override hardcode in frontend that targets "ub-dlza-test" namespace
# RUN sed -i "s|dlza-manager.ub-dlza-test.k8s-001.unibas.ch|dlza-manager.ub-dlza-stage.k8s-001.unibas.ch|g" dlza-frontend/src/client.ts
# RUN sed -i "s|env.PUBLIC_BACKEND_URL|dlza-manager.ub-dlza-prod.k8s-001.unibas.ch|g" dlza-frontend/src/client.ts
RUN cd dlza-frontend && npm i -g vite && npm install husky  && rm package-lock.json && echo "PUBLIC_BACKEND_URL=https://dlza-manager.ub-dlza-stage.k8s-001.unibas.ch/graphql" >> .env && npm install && npm run build
# RUN npm run build dlza-frontend
#RUN cd ..
RUN go build


FROM alpine:latest
RUN apk update && apk add tzdata ca-certificates
WORKDIR /
COPY --from=builder /dlza-manager-clerk /
EXPOSE 8080

ENTRYPOINT ["/dlza-manager-clerk"]


# FROM alpine:latest
# RUN apk update && apk add tzdata ca-certificates
# # COPY --from=builder /dlza-manager-clerk/dlza-manager-clerk /
# # COPY --from=builder ./ub-license .
# COPY --from=builder . /
# EXPOSE 8080
# ENTRYPOINT ["/dlza-manager-clerk/dlza-manager-clerk"]