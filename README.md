[![Go](https://github.com/kotai-tech/server/actions/workflows/ci.yml/badge.svg)](https://github.com/kotai-tech/server/actions/workflows/ci.yml)
[![codecov](https://codecov.io/github/kotai-tech/server/graph/badge.svg?token=NWA2EYXHAW)](https://codecov.io/github/kotai-tech/server)

# Patient management server

A simple patient management server written in Go, following Clean Architecture principles.

## An implementation of Clean Architecture

### Clean Architecture

Clean Architecture principles rely on layers, with dependencies always pointing from outer to inner layers.

Layers described from outer to inner:
- **Drivers**: External systems that depends and interact with the **Interface Adapters** layer.
- **Interfaces adapters**: Contains two responsibilities:
    - Input: Adapt incoming data from **Drivers** (e.g., HTTP requests) into a format expected by the **Usecase** layer.
    - Output: Implement interfaces enabling communication with external systems, as required and defined by the **Usecase** layer.
- **Usecase**: Contains the core business logic, define and use interfaces for external interactions (which are implemented in the Adapters layer) and manipulating **Entities**.
- **Entities**: Contains business objects. They are independent of application logic and external dependencies.


**The Dependency Rule**: Code in each layer may only depend on inside layers.

### Current implementation

```
+---------------------------------------------------------+
|                      Drivers                            |
|---------------------------------------------------------|
| - PostgreSQL (via docker-compose)                       |
| - Echo HTTP server (started in main.go)                 |
+----------------------------▲----------------------------+
                             |
+----------------------------|----------------------------+
|               Interface Adapters (adapters/)            |
|---------------------------------------------------------|
| - rest/         --> Handles HTTP requests               |
| - persistence/  --> Implements PersistenceInterface     |
+----------------------------▲----------------------------+
                             |
+----------------------------|----------------------------+
|           Usecase (services/)                           |
|---------------------------------------------------------|
| - ServiceInterface      --> Describes service behavior  |
| - PersistenceInterface  --> Describes external needs    |
| - serviceImpl struct    --> Implements ServiceInterface |
|                         --> Dep on PersistenceInterface |
+----------------------------▲----------------------------+
                             |
+----------------------------|----------------------------+
|               Domain (domains/)                         |
|---------------------------------------------------------|
| - Business Entities (e.g., Patient, User, etc.)         |
+---------------------------------------------------------+

      <-- All dependencies point inward (Dependency Rule) -->
```

- **Drivers**: A PostgreSQL (managed by docker-compose for dev) and an Echo HTTP server started in the main.go entry point.
- **Interface Adapters** (in the `adapters/` directory): Contains the code required for the Usecases to communicate with external systems:
    - `rest/` Handles incoming HTTP requests.
    - `persistence/` Handles interactions with the database or other storage layers.
- **Usecase** (in the `services/` directory): Each service (or Usecase) defines two key interfaces:
    - `ServiceInterface` Describes how the service can be used.
    - `PersistenceInterface` Describes the external dependencies required by the service (typically storage access).
- **Domain**: Contains the core business entities

## Play with it !

### Run the tests

- Start the dependencies
```bash
$ make dependencies
```

- Run the tests (Also run temporary the server)
```bash
$ make ci-integration
```

### Start the server

- Start the dependencies
```bash
$ make dependencies
```

- Run the server
```bash
$ make run
```
