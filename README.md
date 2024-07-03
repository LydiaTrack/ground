# Lydia Base

This is the base repository for the Lydia backend. It has written in Golang and uses the Gin framework. It behaves as a REST API and uses a MongoDB database as remote or local connection.

## Installation

### Prerequisites

- [Golang](https://golang.org/doc/install) (1.22)
- [MongoDB Image](https://hub.docker.com/_/mongo) (4.4.6)
- [Docker](https://docs.docker.com/get-docker/) (20.10.6)

### Steps

1. Clone the repository

```bash
git clone git@github.com:Lydia/lydia-base.git
```

2. Create a `.env` file in the root directory of the project and add the following variables:

```bash
LYDIA_DB_URI=example.mongodb.com
LYDIA_DB_NAME=renoten-db
# It defines the DB type for mongoDB, it can be either CONTAINER or REMOTE
LYDIA_DB_CONNECTION_TYPE=CONTAINER
# DB Container connection port
LYDIA_DB_PORT=27017
JWT_SECRET=youRS3cr3t
JWT_EXPIRES_IN_HOUR=1
JWT_REFRESH_EXPIRES_IN_HOUR=72
DEFAULT_USER_USERNAME=lydia
DEFAULT_USER_PASSWORD=lydia
DEFAULT_ROLE_NAME=STD_USER
DEFAULT_ROLE_TAGS=STD_ROLE
DEFAULT_ROLE_INFO=Standard user role
```

3. Run the following command to start the project

```bash
go build lydia-base/cmd/lydia-base
```

## Abilities

- [x] User management
- [x] Role management
- [x] Authentication
- [x] Authorization
- [x] Refresh token
- [x] Password encryption
- [ ] Password recovery (to be implemented)

> [!NOTE]
> This project is only a base software for actual projects. It is not intended to be used as a standalone software. But can be used as a base for other projects.
