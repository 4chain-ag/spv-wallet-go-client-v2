# Qucik Guide 

In this directory you can find bunch of examples describing how to use 
the wallet client package during interaction wit the SPV Wallet API. 

1. [Before you run](#before-you-run)
1. How to run example 
1. Proposed order of executing examples 


## Before you run

### Pre-requisites

-   You have access to the `spv-wallet` non-custodial wallet (running locally or remotely).
-   You have installed this package on your machine (`go install` on this project's root directory).

### Concerning the keys

The `ExampleXPub` and `ExampleXPriv` are just placeholders, which won't work. Instead you can:

-  Replace them by newly generated ones using `task generate_keys`
-  Reuse your actual keys if you have them 

> [!CAUTION]
> Don't use the keys which are already added to another wallet.
 
> [!IMPORTANT] 
> Additionally, to make it work properly, you should adjust the `    ExamplePaymail` to align with your `domains` configuration in the `spv-wallet` instance.
