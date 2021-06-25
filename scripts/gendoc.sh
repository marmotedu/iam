# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

for top in pkg internal/pkg
do
    for d in $(find $top -type d)
    do
        if [ ! -f $d/doc.go ]; then
            if ls $d/*.go > /dev/null 2>&1; then
                echo $d/doc.go
                echo "package $(basename $d) // import \"github.com/marmotedu/iam/$d\"" > $d/doc.go
            fi
        fi
    done
done
