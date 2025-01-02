# Quick Guide

In this directory you can find bunch of examples describing how to use
the wallet client package during interaction wit the SPV Wallet API.

1. [Before you run](#before-you-run)
1. [Authorization](#authorization)
1. [How to run example](#how-to-run-an-example)

## Before you run

### Pre-requisites

- You have access to the `spv-wallet` non-custodial wallet (running locally or remotely).
- [Taskfile](https://taskfile.dev/installation/) is installed on your environment.
- SPV Wallet go client instance is properly created and configured.

> [!TIP]
> To verify Taskfile installation run: `task` command in the terminal.

```
task: [default] task --list
task: Available tasks for this project:
* access_key:                        Fetch Access key as User.
* admin_add_user:                    Add user as Admin.
* admin_remove_user:                 Remove user as Admin.
* create_transaction:                Create transaction as User.
* default:                           Display all available tasks.
* generate_keys:                     Generate keys for SPV Wallet API access.
* generate_totp:                     Generate totp.
* get_balance:                       Get balance as User.
* get_handled_paymail_domains:       Get paymail domains as User.
* get_shared_config:                 Get shared config as User.
* list_access_keys:                  Fetch first page of access keys as User.
* list_transactions:                 Fetch first page of transactions as User.
* send_op_return:                    Create draft transaction, finalize transaction and record transaction as User.
* sync_merkleroots:                  Sync Merkle roots as User.
* update_user_xpub_metadata:         Update xPub metadata as User.
* xpriv_from_mnemonic:               Extract xPriv from mnemonic.
* xpub_from_xpriv:                   Extract xPub from xPriv.
```

## Authorization

> [!CAUTION]
> Don't use the keys which are already added to another wallet.

> [!IMPORTANT]
> Additionally, to make it work properly, you should adjust the `ExamplePaymail` to align with your `domains` configuration in the `spv-wallet` instance.

Before interacting with the SPV Wallet API, you must complete the authorization process.

To begin, generate a pair of keys using the `task generate_keys` command, which is included in the dedicated Taskfile.

**Example output:**

```
==================================================================
XPriv:  xprv1d77e47e-452c-453f-bc4c-a42748f8145f
XPub:  xpubd82c277b-0a7e-482f-8ad8-e92958d15acb
Mnemonic:  mnemonic
==================================================================
```

##

> [!TIP]
> Previously generated keys can be used as function parameters.

To verify the connection and authorization, you can either run one of the available code snippets from the examples directory or use the following example. Please note that this is a testable code snippet and should be customized to fit your specific setup.

**Code snippet:**

```
package main

import (
 "context"
 "encoding/json"
 "fmt"
 "log"
 "strings"

 wallet "github.com/bitcoin-sv/spv-wallet-go-client"
)

func main() {
 xPriv := "121d2f43-4261-42ab-813e-3d3fa4d87313"
 cfg := wallet.NewDefaultConfig("http://localhost:3003")
 userAPI, err := wallet.NewUserAPIWithXPriv(cfg, xPriv)
 if err != nil {
  log.Fatal(err)
 }

 xPub, err := userAPI.XPub(context.Background())
 if err != nil {
  log.Fatal(err)
 }
 Print("XPub", xPub)
}

func Print(s string, a any) {
 fmt.Println(strings.Repeat("~", 100))
 fmt.Println(s)
 fmt.Println(strings.Repeat("~", 100))
 res, err := json.MarshalIndent(a, "", " ")
 if err != nil {
  log.Fatal(err)
 }
 fmt.Println(string(res))
}
```

> [!TIP]
> The same principle applies when creating an AdminAPI client instance using one of the available constructors.

**Example output:**

```
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
XPub
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
{
 "createdAt": "2024-10-07T13:39:07.886862Z",
 "updatedAt": "2024-11-20T11:05:22.235832Z",
 "deletedAt": null,
 "metadata": {
  "metadata": {
   "key": "value"
  }
 },
 "id": "c50e4656-75e4-482e-a52d-2b4319919a26",
 "currentBalance": 100,
 "nextInternalNum": 20,
 "nextExternalNum": 2
}
```

## How to run an example

The examples are written in Go and can be run by:

```bash
cd examples
task name_of_the_example
```

> [!TIP]
> To verify Taskfile installation run: `task` command in the terminal.
