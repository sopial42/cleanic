[![Cleanic server CI](https://github.com/sopial42/cleanic/actions/workflows/ci.yml/badge.svg)](https://github.com/sopial42/cleanic/actions/workflows/ci.yml)
[![codecov](https://codecov.io/github/sopial42/cleanic/graph/badge.svg?token=NWA2EYXHAW)](https://codecov.io/github/sopial42/cleanic)

# Cleanic is a patient management server

A simple patient management server written in Go.
- Following Clean Architecture principles.
- Implement JWT auth & RBAC with two roles: admin, doctor.

Details on design choices and implementation can be found later in this README.

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

- In another term, run the integration tests to init the DB with test fixtures
```bash
$ make integration
```


# An implementation of Clean Architecture

## Clean Architecture

Clean Architecture rely on layers, with dependencies always pointing from outer to inner layers.

Layers described from outer to inner:
- **Drivers**: External systems that depends and interact with the **Interface Adapters** layer.
- **Interfaces adapters**: Contains two responsibilities:
    - Input: Adapt incoming data from **Drivers** (e.g., HTTP requests) into a format expected by the **Usecase** layer.
    - Output: Implement interfaces enabling communication with external systems or services, as required and defined by the **Usecase** layer.
- **Usecase**: Contains the core business logic, define and use interfaces for external interactions (which are implemented in the Adapters layer) and manipulating **Entities**.
- **Entities**: Contains business objects. They are independent of application logic and external dependencies.


**The Dependency Rule**: Code in each layer may only depend on inside layers.

## Current implementation

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
| - client/      --> Implements my other services clients |
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


# 🔐 JWT Authentication with Refresh Tokens

A robust authentication flow using JWT `Access Tokens` and `Refresh Tokens` to balance between security and user experience. 

## TLDR;

`Refresh Token`:
- Issued by authentication server upon user login
- Long TTL (days to weeks) so users don’t have to login too often but would open a long window time for attackers if it leaks
- Stored securely server-side in DB, enabling proactive security mechanisms:
    - Rotation
    - Revocation
    - Binding to contextual data (IP address, device fingerprint…)
- Stored client-side in an `HttpOnly` cookie, protected from XSS attacks

Short term token - `Access token`:
- Issued by the authentication server in response to a valid `Refresh Token`
- Short TTL (e.g., 5-15 minutes) to limits the damage in case of token leak
- Used to authenticate requests to services
- Stored in-memory on the client, less safe
- `Access Tokens` are stateless, so there is no way to invalidate them
- Automatically refreshed to keep the user experience away from annoying re-logins


| Token Type    | TTL         | Stored On (Client)     | Revocable | Used For               |
|---------------|-------------|------------------------|-----------|------------------------|
| Refresh Token | Days/Weeks  | HttpOnly Cookie        |    ✅     | Request. `Access Tokens` |
| `Access Token`  | Minutes     | Memory               |    ❌     | Call backend services  |



## Why ? 

`Refresh Tokens` rely on a strong authentication mechanism as user needs to login in order to obtain it from the authentication server.
On the client side, this token is stored in a cookie HttpOnly, which makes it more secure on client side (no XSS attacks).
On the server side, this token is stored in DB. We can enable multiple pro-active security features: it can be rotated, revoked and can even be associated with stronger authentication details (IP, device fingerprint...)
But once the token is emitted, for UX purpose we can suppose the token validity has to be at least for a day to several weeks.

💣  Risk : We avoid using the `Refresh Token` directly for authenticating with services, for several key reasons:
- If a `Refresh Token` leaks undetected, an attacker has a long window to act using this token
- It's stateful and we don't want every microservices rely a single DB/auth service to works
- If it was stateless, we can't revoke it so a leak would be a mess to fix before the end of the TTL

⚕️ Solution : We don't use directly the `Refresh Token` to authenticate over services. We'll use a token with short TTL (5-15 minutes) -> the `Access Token`.

Acces token is issued by the authentication server in response to a valid `Refresh Token` and is only used to authenticate requests to services through the `Authorization` header.
On client side, the `Access Token` is stored in memory to be able to be manipulated by the client (headers, cross-domain requests...)
`Access token` is stateless so their is no way to invalid it pro-actively as seen before(*)... 
But it has only a short-term duration. In case of a leak, attack surface is significantly reduced

⚠️ Caveats & Considerations:
- If an `access token` leaks, it's valid until it expires (5–15 mins).
- If a `Refresh Token` is revoked, all future `access token` requests fail — but existing `Access Tokens` remain valid for their TTL
- *) `Access Tokens` are not revocable unless you make them stateful, which adds complexity and overhead (e.g., token blacklists).

## ✅ Benefits Recap

- Secure long-term authentication with `Refresh Tokens`
- Minimized exposure via short-lived `Access Tokens`
- Seamless UX with automatic token refresh

### Refresh

- Login service retourne AccessToken + RefreshToken
- Stocker RefreshToken en DB    

### Aud feature
- service name + middleware check its own service name 
- gérer audience en dur pour l'instant avec aud AUTH + aud ["user", "patient"]

## Refresh token claims

- sub: identifiant user
- exp: expiration date
- iat: date de création
- aud: audience (facultatif, auth API)
- token_id: à stocker en DB pour revocation si besoin


## Access token 

- sub: identifiant user
- exp: expiration date
- iat: date de création
- aud: audience (auth API)
- roles: 

# TODO: test avec des durées de token de quelques secondes


// Token features
- 1 token per user
- lastUpdate reason (logout / refresh)
