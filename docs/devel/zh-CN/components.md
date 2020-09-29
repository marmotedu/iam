# IAM Components

`/cmd` directory includes every IAM components and is where all binaries and container images are built. For detail about how to launch the IAM cluster see the guide [here](/docs/devel/running-locally.md).

## Overview

IAM contains 12 core components belonging to 6 services, a dependency list generator and a customized installer.

## Core Components
To bootstrap properly, IAM core components need to be run in the order as shown below.

- [`iam-apiserver`](/cmd/iam-apiserver) integrates [dex](https://github.com/dexidp/dex) to provide an OpenID Connect server, which can provide access to third-party authentication systems, and also provides a default local identify.
