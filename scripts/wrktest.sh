#!/bin/bash

# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.    
# Use of this source code is governed by a MIT style    
# license that can be found in the LICENSE file.

: << EOF
API 性能测试脚本，会自动执行 wrk 命令，采集数据、分析数据并调用 gnuplot 画图

使用方式 ( 测试 API 性能)：
1. 启动 apiserver (8080端口)
2. 执行测试脚本: ./wrktest.sh

脚本会生成 _wrk.dat 的数据文件，每列含义为：

并发数 QPS  平均响应时间 成功率

使用方式 (对比2次测试结果)

1. 执行命令： ./wrktest.sh diff apiserver1_wrk.dat http_wrk.dat
 
> Note: 需要确保系统安装了 wrk 和 gnuplot 工具
EOF

t1="apiserver2" # 对比图中红色线条名称
t2="http" # 对比图中粉色线条名称
jobname="apiserver2" # 本次测试名称

## wrk 参数配置
d="300s" 
concurrent="200 500 1000 3000 5000 10000 15000 20000 25000 50000 100000 200000 500000 1000000"
threads=144

if [ "$1" != "" ];then
	url="$1"
else
	url="http://127.0.0.1:8080/sd/health"
fi

cmd="wrk --latency -t$threads -d$d -T30s $url"
apiperformance="${jobname}_performance.png"
apisuccessrate="${jobname}_success_rate.png"
datfile="${jobname}_wrk.dat"

# functions
function convertPlotData() 
{
	echo "$1" | awk -v datfile="$datfile" ' {
		if ($0 ~ "Running") {
			common_time=$2
		}
		if ($0 ~ "connections") {
			connections=$4
			common_threads=$1
		}
		if ($0 ~ "Latency   ") {
			avg_latency=convertLatency($2)
		}
		if ($0 ~ "50%") {
			p50=convertLatency($2)
		}
		if ($0 ~ "75%") {
			p75=convertLatency($2)
		}
		if ($0 ~ "90%") {
			p90=convertLatency($2)
		}
		if ($0 ~ "99%") {
			p99=convertLatency($2)
		}
		if ($0 ~ "Requests/sec") {
			qps=$2
		}
		if ($0 ~ "requests in") {
			allrequest=$1
		}
		if ($0 ~ "Socket errors") {
			err=$4+$6+$8+$10
		}

	}
	END {
		rate=sprintf("%.2f", (allrequest-err)*100/allrequest)
		print connections,qps,avg_latency,rate >> datfile
	}

	function convertLatency(s) {
		if (s ~ "us") {
			sub("us", "", s)
			return s/1000
		}
		if (s ~ "ms") {
			sub("ms", "", s)
			return s
		}
		if (s ~ "s") {
			sub("s", "", s)
			return s * 1000
		}
	}
	
	'
}

function prepare()
{
	rm -f $datfile 
}

function plot() {
	gnuplot <<  EOF
set terminal png enhanced #输出格式为png文件
set output "$apiperformance"  #指定数据文件名称
set title "QPS & TTLB\nRunning: 300s\nThreads: $threads"
set ylabel 'QPS'
set xlabel 'Concurrent'
set y2label 'Average Latency (ms)'
set key top left vertical noreverse spacing 1.2 box
set tics out nomirror
set border 3 front
set style line 1 linecolor rgb '#00ff00' linewidth 2 linetype 3 pointtype 2
set style line 2 linecolor rgb '#ff0000' linewidth 1 linetype 3 pointtype 2
set style data linespoints

set grid #显示网格
set xtics nomirror rotate #by 90#只需要一个x轴
set mxtics 5
set mytics 5 #可以增加分刻度
set ytics nomirror
set y2tics

set autoscale  y
set autoscale y2

plot "$datfile" using 2:xticlabels(1) w lp pt 7 ps 1 lc rgbcolor "#EE0000" axis x1y1 t "QPS","$datfile" using 3:xticlabels(1) w lp pt 5 ps 1 lc rgbcolor "#0000CD" axis x2y2 t "Avg Latency (ms)"

unset y2tics
unset y2label
set ytics nomirror
set yrange[0:100]
set output "$apisuccessrate"  #指定数据文件名称
set title "Success Rate\nRunning: 300s\nThreads: $threads"
plot "$datfile" using 4:xticlabels(1) w lp pt 7 ps 1 lc rgbcolor "#F62817" t "Success Rate"
EOF
}

function plotDiff()
{
	gnuplot <<  EOF
set terminal png enhanced #输出格式为png文件
set output "${t1}_$t2.qps.diff.png"  #指定数据文件名称
set title "QPS & TTLB\nRunning: 300s\nThreads: $threads"
set xlabel 'Concurrent'
set ylabel 'QPS'
set y2label 'Average Latency (ms)'
set key below left vertical noreverse spacing 1.2 box autotitle columnheader
set tics out nomirror
set border 3 front
set style line 1 linecolor rgb '#00ff00' linewidth 2 linetype 3 pointtype 2
set style line 2 linecolor rgb '#ff0000' linewidth 1 linetype 3 pointtype 2
set style data linespoints

#set border 3 lt 3 lw 2   #这会让你的坐标图的border更好看
set grid #显示网格
set xtics nomirror rotate #by 90#只需要一个x轴
set mxtics 5
set mytics 5 #可以增加分刻度
set ytics nomirror
set y2tics

#set pointsize 0.4 #点的像素大小
#set datafile separator '\t' #数据文件的字段用\t分开

set autoscale  y
set autoscale y2

#设置图像的大小 为标准大小的2倍
#set size 2.3,2

plot "/tmp/plot_diff.dat" using 2:xticlabels(1) w lp pt 7 ps 1 lc rgbcolor "#EE0000" axis x1y1 t "$t1 QPS","/tmp/plot_diff.dat" using 5:xticlabels(1) w lp pt 7 ps 1 lc rgbcolor "#EE82EE" axis x1y1 t "$t2 QPS","/tmp/plot_diff.dat" using 3:xticlabels(1) w lp pt 5 ps 1 lc rgbcolor "#0000CD" axis x2y2 t "$t1 Avg Latency (ms)", "/tmp/plot_diff.dat" using 6:xticlabels(1) w lp pt 5 ps 1 lc rgbcolor "#6495ED" axis x2y2 t "$t2 Avg Latency (ms)"

unset y2tics
unset y2label
set ytics nomirror
set yrange[0:100]
set title "Success Rate\nRunning: 300s\nThreads: $threads"
set output "${t1}_$t2.success_rate.diff.png"  #指定数据文件名称
plot "/tmp/plot_diff.dat" using 4:xticlabels(1) w lp pt 7 ps 1 lc rgbcolor "#EE0000" t "$t1 Success Rate","/tmp/plot_diff.dat" using 7:xticlabels(1) w lp pt 7 ps 1 lc rgbcolor "#EE82EE" t "$t2 Success Rate"
EOF
}

if [ "$1" == "diff" ];then
	join $2 $3 > /tmp/plot_diff.dat
	plotDiff `basename $2` `basename $3`
	exit 0
fi


prepare

for c in $concurrent
do
	wrkcmd="$cmd -c $c"
	echo -e "\nRunning wrk command: $wrkcmd"
	result=`eval $wrkcmd`
	convertPlotData "$result"
done

echo -e "\nNow plot according to $datfile"
plot &> /dev/null
echo -e "QPS graphic file is: $apiperformance\nSuccess rate graphic file is: $apisuccessrate"
