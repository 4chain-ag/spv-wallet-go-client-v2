<div align="center">

# SPV Wallet: Go Client

[![Release](https://img.shields.io/github/release-pre/bitcoin-sv/spv-wallet-go-client.svg?logo=github&style=flat&v=2)](https://github.com/bitcoin-sv/spv-wallet-go-client/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/bitcoin-sv/spv-wallet-go-client/run-tests.yml?branch=main&v=2)](https://github.com/bitcoin-sv/spv-wallet-go-client/actions)
[![Report](https://goreportcard.com/badge/github.com/bitcoin-sv/spv-wallet-go-client?style=flat&v=2)](https://goreportcard.com/report/github.com/bitcoin-sv/spv-wallet-go-client)
[![codecov](https://codecov.io/gh/bitcoin-sv/spv-wallet-go-client/branch/main/graph/badge.svg?v=2)](https://codecov.io/gh/bitcoin-sv/spv-wallet-go-client)
[![Mergify Status](https://img.shields.io/endpoint.svg?url=https://api.mergify.com/v1/badges/bitcoin-sv/spv-wallet-go-client&style=flat&v=2)](https://mergify.io)
<br>

[![Go](https://img.shields.io/github/go-mod/go-version/bitcoin-sv/spv-wallet-go-client?v=2)](https://golang.org/)
[![Gitpod Ready-to-Code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod&v=2)](https://gitpod.io/#https://github.com/bitcoin-sv/spv-wallet-go-client)
[![standard-readme compliant](https://img.shields.io/badge/readme%20style-standard-brightgreen.svg?style=flat&v=2)](https://github.com/RichardLitt/standard-readme)
[![Makefile Included](https://img.shields.io/badge/Makefile-Supported%20-brightgreen?=flat&logo=probot&v=2)](Makefile)


<br/>
</div>

## Table of Contents
1. [Requirements and Compatibility](#requirements-and-compatibility)
1. [Quick start](#quick-start)
1. Documentation
1. Testing
1. [Examples](/examples/README.md)
1. Contributing  
1. License

## Requirements and Compatibility

Instalation: 
```shell script
go get -u github.com/bitcoin-sv/spv-wallet-go-client
```

## Requirements

- **Go Version**: The `spv-wallet-go-client` requires **Go version 1.22.5** or a later supported release of Go. Ensure your Go environment meets this requirement before using the client.


## Compatibility and Support

### Deprecation Notice
The client **does not support** the following:
- **Admin and non-admin old endpoints** of the SPV Wallet API based on the `/v1/` prefix.
- Deprecated methods for building query parameters for HTTP requests.

### Current Compatibility
The client is designed for full compatibility with the newer `/api/v1/` endpoints exposed by the SPV Wallet API. It focuses on aligning with the latest standards and structure provided by the API.
 
### API Admin Endpoints Compatibility

#### Access Keys API
| HTTP Method | Endpoint                     | Action               | Support Status | API Code                                          |   Pagination     |
|-------------|-------------------------------|----------------------|----------------|--------------------------------------------------|----------------- | 
| GET         | /api/v1/admin/users/keys     | Search access keys   | ✅             | [API](/internal/api/v1/admin/accesskeys/access_keys_api.go#L25) | ✅ |

#### Contacts API
| HTTP Method | Endpoint                              | Action               | Support Status | API Code                                       |  Pagination   |
|-------------|---------------------------------------|----------------------|----------------|------------------------------------------------|-------------- |
| GET         | /api/v1/admin/contacts               | Search contacts      | ✅             | [API](/internal/api/v1/admin/contacts/contacts_api.go#L42) | ✅ |
| POST        | /api/v1/admin/contacts/confirmations | Confirm contact      | ✅             | [API](/internal/api/v1/admin/contacts/contacts_api.go#L83) | ❌ |
| PUT         | /api/v1/admin/contacts/{id}          | Update contact       | ✅             | [API](/internal/api/v1/admin/contacts/contacts_api.go#L68) | ❌ |
| DELETE      | /api/v1/admin/contacts/{id}          | Delete contact       | ✅             | [API](/internal/api/v1/admin/contacts/contacts_api.go#L95) | ❌ |
| POST        | /api/v1/admin/contacts/{paymail}     | Create contact       | ✅             | [API](/internal/api/v1/admin/contacts/contacts_api.go#L27) | ❌ |

#### Invitations API
| HTTP Method | Endpoint                              | Action               | Support Status | API Code                                         |   Pagination      |
|-------------|---------------------------------------|----------------------|----------------|--------------------------------------------------|-------------------|
| POST        | /api/v1/admin/invitations/{id}       | Accept invitation    | ✅             | [API](/internal/api/v1/admin/invitations/invitations_api.go#L22) | ❌ |
| DELETE      | /api/v1/admin/invitations/{id}       | Reject invitation    | ✅             | [API](/internal/api/v1/admin/invitations/invitations_api.go#L35) | ❌ |


#### Paymails API
| HTTP Method | Endpoint                              | Action               | Support Status | API Code                                         |  Pagination      |
|-------------|---------------------------------------|----------------------|----------------|--------------------------------------------------|------------------|
| GET         | /api/v1/admin/paymails               | Search paymails      | ✅             | [API](/internal/api/v1/admin/paymails/paymails_api.go#L73) | ✅      |
| POST        | /api/v1/admin/paymails               | Create paymail       | ✅             | [API](/internal/api/v1/admin/paymails/paymails_api.go#L44) | ❌      |
| GET         | /api/v1/admin/paymails/{id}          | Retrieve paymail     | ✅             | [API](/internal/api/v1/admin/paymails/paymails_api.go#L59) | ❌      |
| DELETE      | /api/v1/admin/paymails/{id}          | Delete paymail       | ✅             | [API](/internal/api/v1/admin/paymails/paymails_api.go#L27) | ❌      |

#### Stats API
| HTTP Method | Endpoint                     | Action               | Support Status | API Code                                             |  Pagination   |
|-------------|-------------------------------|----------------------|----------------|-----------------------------------------------------|---------------|
| GET         | /api/v1/admin/stats          | Retrieve stats       | ✅             | [API](/internal/api/v1/admin/stats/stats_api.go#L23) |     ✅        |

#### Status API
| HTTP Method | Endpoint                     | Action               | Support Status | API Code                                               | Pagination      |
|-------------|-------------------------------|----------------------|----------------|-------------------------------------------------------|-----------------|
| GET         | /api/v1/admin/status         | Retrieve status      | ✅             | [API](/internal/api/v1/admin/status/status_api.go#L23) |      ❌         |

#### Transactions API
| HTTP Method | Endpoint                              | Action               | Support Status | API Code                                         |       Pagination      |
|-------------|---------------------------------------|----------------------|----------------|--------------------------------------------------|-----------------------|
| GET         | /api/v1/admin/transactions           | Search transactions | ✅             | [API](/internal/api/v1/admin/transactions/transactions_api.go#L39) | ✅    |
| GET         | /api/v1/admin/transactions/{id}      | Retrieve transaction | ✅             | [API](/internal/api/v1/admin/transactions/transactions_api.go#L26)| ❌    |

#### UTXOs API
| HTTP Method | Endpoint                              | Action               | Support Status | API Code                                            |    Pagination    |
|-------------|---------------------------------------|----------------------|----------------|-----------------------------------------------------| -----------------|
| GET         | /api/v1/admin/utxos                  | Search UTXOs         | ✅             | [API](/internal/api/v1/admin/utxos/utxos_api.go#L25) | ✅               |

#### Webhooks API
| HTTP Method | Endpoint                              | Action               | Support Status | API Code                                          |   Pagination  |
|-------------|---------------------------------------|----------------------|----------------|---------------------------------------------------|---------------|
| GET         | /api/v1/admin/webhooks/subscriptions | Subscribe to webhook | ✅             | [API](/internal/api/v1/admin/webhooks/webhooks_api.go#L23) |  ❌   |
| DELETE      | /api/v1/admin/webhooks/subscriptions | Unsubscribe webhook  | ✅             | [API](/internal/api/v1/admin/webhooks/webhooks_api.go#L36) |  ❌   |

#### XPubs API
| HTTP Method | Endpoint                              | Action               | Support Status | API Code                                            |  Pagination |
|-------------|---------------------------------------|----------------------|----------------|-----------------------------------------------------|-------------|
| GET         | /api/v1/admin/users                  | Search XPubs         | ✅             | [API](/internal/api/v1/admin/xpubs/xpubs_api.go#L41) |  ✅         |
| POST        | /api/v1/admin/users                  | Create XPub          | ✅             | [API](/internal/api/v1/admin/xpubs/xpubs_api.go#L27) |  ❌         |

### API Non-Admin Endpoints Compatibility

#### Access Keys API
| HTTP Method | Endpoint                     | Action               | Support Status | API Code                                          |  Pagination      |
|-------------|-------------------------------|----------------------|----------------|--------------------------------------------------|------------------|
| GET         | /api/v1/users/current/keys   | Search access keys   | ✅             | [API](/internal/api/v1/user/accesskeys/access_key_api.go#L56)   | ✅ |
| POST        | /api/v1/users/current/keys   | Create access key    | ✅             | [API](/internal/api/v1/user/accesskeys/access_key_api.go#L27)   | ❌ |
| GET         | /api/v1/users/current/keys/{id} | Retrieve access key | ✅             | [API](/internal/api/v1/user/accesskeys/access_key_api.go#L42) | ❌ |
| DELETE      | /api/v1/users/current/keys/{id} | Revoke access key   | ✅             | [API](/internal/api/v1/user/accesskeys/access_key_api.go#L82) | ❌ |

#### Contacts API
| HTTP Method | Endpoint                     | Action               | Support Status | API Code                                          |  Pagination  |
|-------------|-------------------------------|----------------------|----------------|--------------------------------------------------|--------------|
| GET         | /api/v1/contacts             | Search contacts      | ✅             | [API](/internal/api/v1/user/contacts/contacts_api.go#L27) | ✅   |
| GET         | /api/v1/contacts/{paymail}   | Retrieve contact     | ✅             | [API](/internal/api/v1/user/contacts/contacts_api.go#L53) | ❌   |
| PUT         | /api/v1/contacts/{paymail}   | Upsert contact       | ✅             | [API](/internal/api/v1/user/contacts/contacts_api.go#L67) | ❌   |
| DELETE      | /api/v1/contacts/{paymail}   | Remove contact       | ✅             | [API](/internal/api/v1/user/contacts/contacts_api.go#L89) | ❌   |
| POST        | /api/v1/contacts/{paymail}   | Confirm contact      | ✅             | [API](/internal/api/v1/user/contacts/contacts_api.go#L101)| ❌   |
| DELETE      | /api/v1/contacts/{paymail}   | Unconfirm contact    | ✅             | [API](/internal/api/v1/user/contacts/contacts_api.go#L113)| ❌   |

#### Invitations API
| HTTP Method | Endpoint                     | Action               | Support Status | API Code                                          |  Pagination               |
|-------------|-------------------------------|----------------------|----------------|--------------------------------------------------|---------------------------|
| POST        | /api/v1/invitations/{paymail}/contacts | Accept invitation   | ✅             | [API](/internal/api/v1/user/invitations/invitations_api.go#L22) | ❌ |
| DELETE      | /api/v1/invitations/{paymail}          | Reject invitation   | ✅             | [API](/internal/api/v1/user/invitations/invitations_api.go#L34) | ❌ |

#### Merkle Roots API
| HTTP Method | Endpoint                     | Action               | Support Status | API Code                                          |  Pagination       |
|-------------|-------------------------------|----------------------|----------------|--------------------------------------------------|-------------------|
| GET         | /api/v1/merkleroots          | Search Merkle roots  | ✅             | [API](/internal/api/v1/user/merkleroots/merkleroots_api.go#L36)| ❌   |

#### Paymails API
| HTTP Method | Endpoint                     | Action               | Support Status | API Code                                          | Pagination       |
|-------------|-------------------------------|----------------------|----------------|--------------------------------------------------|------------------|
| GET         | /api/v1/paymails             | Search paymails      | ✅             | [API](/internal/api/v1/user/paymails/paymails_api.go#L25) | ✅       |

#### Transactions API
| HTTP Method | Endpoint                     | Action               | Support Status | API Code                                          |     Pagination       |
|-------------|-------------------------------|----------------------|----------------|--------------------------------------------------|----------------------|
| GET         | /api/v1/transactions         | Search transactions  | ✅             | [API](/internal/api/v1/user/transactions/transactions_api.go#L137) |✅   |
| POST        | /api/v1/transactions         | Record transaction   | ✅             | [API](/internal/api/v1/user/transactions/transactions_api.go#L93) |❌    |
| POST        | /api/v1/transactions/drafts  | Draft transaction    | ✅             | [API](/internal/api/v1/user/transactions/transactions_api.go#L78) |❌    |
| GET         | /api/v1/transactions/{id}    | Retrieve transaction | ✅             | [API](/internal/api/v1/user/transactions/transactions_api.go#L123) |❌   |
| PATCH       | /api/v1/transactions/{id}    | Update transaction   | ✅             | [API](/internal/api/v1/user/transactions/transactions_api.go#L108) |❌   |

#### UTXOs API
| HTTP Method | Endpoint                     | Action               | Support Status | API Code                                            | Pagination  |
|-------------|-------------------------------|----------------------|----------------|----------------------------------------------------|---------------|
| GET         | /api/v1/utxos                | Search UTXOs         | ✅             | [API](/internal/api/v1/user/utxos/utxos_api.go#L25) |          ❌   |

#### XPubs API
| HTTP Method | Endpoint                     | Action                       | Support Status | API Code                                           |Pagination |
|-------------|-------------------------------|------------------------------|----------------|---------------------------------------------------|-----------|
| GET         | /api/v1/users/current        | Retrieve current user info   | ✅             | [API](/internal/api/v1/user/xpubs/xpub_api.go#L24) |  ❌       |
| PATCH       | /api/v1/users/current        | Update current user info     | ✅             | [API](/internal/api/v1/user/xpubs/xpub_api.go#L24) |  ❌       |



## Feature Updates

While the client strives to support the latest API features, there may be a delay in fully integrating new functionalities. If you encounter any issues or have questions:
- Refer to the official documentation.
- Reach out for support to ensure a smooth development experience.


 
## Quick start

The implementation enforces separation of concerns by isolating admin and non-admin APIs, requiring separate initialization for their respective clients. This ensures clarity and modularity when utilizing the exposed functionality. 
 
### `UserAPI` Initialization Methods:

### 1. [`NewUserAPIWithAccessKey`](/user_api.go#L468)
- **Description:** Initializes a `UserAPI` instance using an access key for authentication.
- **Note:** Requests made with this instance will be securely signed, ensuring integrity and authenticity.

### 2. [`NewUserAPIWithXPriv`](/user_api.go#L449)
- **Description:** Initializes a `UserAPI` instance using an extended private key (xPriv) for authentication.
- **Note:** Requests made with this instance will also be securely signed.
- **Recommendation:** This option offers a high level of security, making it a preferred choice alongside the access key option.

### 3. [`NewUserAPIWithXPub`](/user_api.go#L435)
- **Description:** Initializes a `UserAPI` instance using an extended public key (xPub).
- **Note:** Requests made with this instance will not be signed.
- **Security Advisory:** For enhanced security, it is strongly recommended to use either `NewUserAPIWithAccessKey` or `NewUserAPIWithXPriv` instead, as unsigned requests may be less secure.


### `AdminAPI` Initialization Methods:

### 1. [`NewAdminAPIWithXPriv`](/admin_api.go#L375)
- **Description:** Initializes a `AdminAPI` instance using an extended private key (xPriv) for authentication.
- **Note:** Requests made with this instance will be securely signed, ensuring integrity and authenticity.

### 2. [`NewAdminAPIWithXPub`](/admin_api.go#L390)
- **Description:** Initializes a `AdminAPI` instance using an extended public key (xPub).
- **Note:** Requests made with this instance will not be signed.
- **Security Advisory:** For enhanced security, it is strongly recommended to use either `NewAdminAPIWithXPriv`instead, as unsigned requests may be less secure.

**Code snippets:**
- [AdminAPI example](/examples/admin_add_user/admin_add_user.go)
- [UserAPI example](/examples/list_transactions/list_transactions.go)
