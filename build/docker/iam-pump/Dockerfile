# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.    
# Use of this source code is governed by a MIT style    
# license that can be found in the LICENSE file.

FROM BASE_IMAGE
MAINTAINER Lingfei Kong <colin404@foxmail.com>

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime     
RUN echo "Asia/Shanghai" > /etc/timezone 

WORKDIR /opt/iam
COPY iam-pump /opt/iam/bin/

ENTRYPOINT ["/opt/iam/bin/iam-pump"]
