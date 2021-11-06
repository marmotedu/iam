# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.    
# Use of this source code is governed by a MIT style    
# license that can be found in the LICENSE file.

FROM BASE_IMAGE
LABEL maintainer="<colin404@foxmail.com>"

WORKDIR /opt/iam/scripts/install

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
      echo "Asia/Shanghai" > /etc/timezone

COPY iamctl /usr/bin/
COPY scripts /opt/iam/scripts

ENTRYPOINT ["sleep", "3600"]
