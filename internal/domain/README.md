# Domain Subdirectory

This directory contains the domain definitions for the internal domains. 

Each domain defines a set of entities and their relationships. A domain can have service layer or repository layer definitions, but it is not required.

For every domain, directory structure must be defined as follows:

```
domain
├── README.md
├── {domain_name}
│   ├── {domain_name}.go
│   ├── commands
│   │   └── {command_name}.go
│   ├── responses
│   │   └── {response_name}.go
```

## Domain

The {domain_name}.go file defines the domain entity and its relationships.

## Commands

The commands directory contains the command definitions for the domain. A command is a service layer definition that defines the business logic for a domain.

## Responses

The responses directory contains the response definitions for the domain. A response is a repository layer definition that defines the data structure for a domain.