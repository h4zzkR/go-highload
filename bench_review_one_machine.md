## case1
- 2 go servers (8080,8081)
- default nginx lb params (no params)
- cache: proxy_cache_path /var/cache/nginx levels=1:2 keys_zone=my_cache:10m;
         proxy_cache_valid 200 60m;

`
Server Software:        nginx/1.23.3
Server Hostname:        127.0.0.1
Server Port:            80

Document Path:          /date
Document Length:        320001 bytes

Concurrency Level:      350
Time taken for tests:   6.516 seconds
Complete requests:      20000
Failed requests:        0
Keep-Alive requests:    0
Total transferred:      6402780000 bytes
HTML transferred:       6400020000 bytes
Requests per second:    3069.42 [#/sec] (mean)
Time per request:       114.028 [ms] (mean)
Time per request:       0.326 [ms] (mean, across all concurrent requests)
Transfer rate:          959611.87 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.8      0       9
Processing:    39  113  13.2    110     256
Waiting:        5  112  12.1    110     256
Total:         48  113  13.7    111     260

Percentage of the requests served within a certain time (ms)
  50%    111
  66%    112
  75%    113
  80%    113
  90%    118
  95%    127
  98%    144
  99%    184
 100%    260 (longest request)
`

1. Requests per second:    3069.42 [#/sec] (mean)
1. Time per request:       114.028 [ms] (mean)
1. Time per request:       0.326 [ms] (mean, across all concurrent requests)
1. Transfer rate:          959611.87 [Kbytes/sec] received

## case2
- 2 go servers (8080,8081)
- default nginx lb params (no params)
- cache: no cache

`
Server Software:        nginx/1.23.3
Server Hostname:        127.0.0.1
Server Port:            80

Document Path:          /date
Document Length:        320001 bytes

Concurrency Level:      350
Time taken for tests:   21.336 seconds
Complete requests:      20000
Failed requests:        0
Keep-Alive requests:    0
Total transferred:      6402780000 bytes
HTML transferred:       6400020000 bytes
Requests per second:    937.38 [#/sec] (mean)
Time per request:       373.382 [ms] (mean)
Time per request:       1.067 [ms] (mean, across all concurrent requests)
Transfer rate:          293057.70 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.9      0      10
Processing:    21  370  77.7    362     798
Waiting:        7  346  73.5    343     763
Total:         21  370  77.5    362     798

Percentage of the requests served within a certain time (ms)
  50%    362
  66%    387
  75%    406
  80%    419
  90%    464
  95%    509
  98%    565
  99%    611
 100%    798 (longest request)
`

1. Requests per second:    937.38 [#/sec] (mean)
1. Time per request:       373.382 [ms] (mean)
1. Time per request:       1.067 [ms] (mean, across all concurrent requests)
1. Transfer rate:          293057.70 [Kbytes/sec] received



## case3
- 1 go server (8080)
- default nginx lb params (no params)
- cache: proxy_cache_path /var/cache/nginx levels=1:2 keys_zone=my_cache:10m;
         proxy_cache_valid 200 60m;

`
Server Software:        nginx/1.23.3
Server Hostname:        127.0.0.1
Server Port:            80

Document Path:          /date
Document Length:        320001 bytes

Concurrency Level:      350
Time taken for tests:   6.373 seconds
Complete requests:      20000
Failed requests:        0
Keep-Alive requests:    0
Total transferred:      6402780000 bytes
HTML transferred:       6400020000 bytes
Requests per second:    3138.16 [#/sec] (mean)
Time per request:       111.530 [ms] (mean)
Time per request:       0.319 [ms] (mean, across all concurrent requests)
Transfer rate:          981100.26 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   1.9      0      22
Processing:    41  110  13.5    109     255
Waiting:        8  109  10.4    108     234
Total:         41  110  14.8    109     263

Percentage of the requests served within a certain time (ms)
  50%    109
  66%    109
  75%    110
  80%    110
  90%    112
  95%    112
  98%    115
  99%    198
 100%    263 (longest request)
`

1. Requests per second:    3138.16 [#/sec] (mean)
1. Time per request:       111.530 [ms] (mean)
1. Time per request:       0.319 [ms] (mean, across all concurrent requests)
1. Transfer rate:          981100.26 [Kbytes/sec] received


## case4
- 3 go servers (8080,8081,8082)
- default nginx lb params (no params)
- cache: proxy_cache_path /var/cache/nginx levels=1:2 keys_zone=my_cache:10m;
         proxy_cache_valid 200 60m;

`
Server Software:        nginx/1.23.3
Server Hostname:        127.0.0.1
Server Port:            80

Document Path:          /date
Document Length:        320001 bytes

Concurrency Level:      350
Time taken for tests:   6.673 seconds
Complete requests:      20000
Failed requests:        0
Keep-Alive requests:    0
Total transferred:      6402780000 bytes
HTML transferred:       6400020000 bytes
Requests per second:    2997.00 [#/sec] (mean)
Time per request:       116.783 [ms] (mean)
Time per request:       0.334 [ms] (mean, across all concurrent requests)
Transfer rate:          936970.86 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   1.8      0      22
Processing:     6  115  13.4    112     210
Waiting:        6  114  12.6    111     186
Total:         27  115  13.7    112     217

Percentage of the requests served within a certain time (ms)
  50%    112
  66%    114
  75%    117
  80%    119
  90%    133
  95%    142
  98%    154
  99%    165
 100%    217 (longest request)
`

1. Requests per second:    2997.00 [#/sec] (mean)
1. Time per request:       116.783 [ms] (mean)
1. Time per request:       0.334 [ms] (mean, across all concurrent requests)
1. Transfer rate:          936970.86 [Kbytes/sec] received

