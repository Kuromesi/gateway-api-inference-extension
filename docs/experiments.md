# Experiment Settings

A kubernetes cluster with 5 nodes, and each node has 16 vCPUs, 60 GB of RAM, and an A10 GPU. 10 pods that loaded the Lora Llama2 model are running in this cluster.

# Filter chains

+ baseline：LEAST_REQUEST load balancing algorithm
+ maxium：low latency -> lora-affinity -> least queueing -> least-kv-cache
+ simple-queue：least queueing (with minimum queue size)
+ simple-kvcache：least-kv-cache (with minimum kv-cache size)
+ queue + kvcache：least queueing -> least-kv-cache

# Test Cases
## base-model
Only request llama2 model.



### Performance
| | baseline | maxium | simple-queue | simple-kvcache | queue+kvcache |
| --- | --- | --- | --- | --- | --- |
| requests per minute | 80.09 | 91.68<br/>88.87<br/>93.35<br/>89.70<br/>85.18 | 80.81<br/>80.26 | 89.32<br/>90.10 | 93.21 |
| average time to first token (ttft) | 2.99 | 0.29<br/>0.62<br/>0.87<br/>0.75<br/>0.39 | 1.42<br/>1.10 | 0.85<br/>0.71 | 0.8 |
| ttft P95 | 17.07 | 0.43<br/>1.38<br/>2.21<br/>1.49<br/>0.47 | 5.9<br/>4.52 | 1.78<br/>0.69 | 1.9 |
| Average inter token latency | 0.15 | 0.05<br/>0.05<br/>0.07<br/>0.06<br/>0.05 | 0.09<br/>0.08 | 0.06<br/>0.07 | 0.06 |
| average token throughoutput per second | 16.95 | 18.20<br/>17.33<br/>16.90<br/>17.01<br/>18.04 | 16.25<br/>17.34 | 16.69<br/>18.31 | 16.61 |
| average end to end latency | 25.09 | 21.21<br/>21.85<br/>21.40<br/>21.91<br/>23.27 | 23.69<br/>24.99 | 21.83<br/>22.72 | 20.45 |
| cache utilization |  | 73.6% | 74.7% | 77.1% |  |
| queue time |  | 38.8ms | 236ms | 129ms |  |
| preemptions total |  | 1047 | 1452 | 784 |  |

## multi-lora

Requests to llama2 model with 10 lora, with 40 connections to different lora

| | baseline | maxium | simple-queue | simple-kvcache | queue+kvcache | lora+queue | lora+kvcache |
| --- | --- | --- | --- | --- | --- | --- | --- |
| requests per minute | 54.28<br/>50.554 | 64.29<br/>57.76<br/>57.77 | 49.49 | 46.13<br/>32.68 | 39.14 | 40.84 | 39.88 |
| average time to first token (ttft) | 13.55<br/>11.60 | 4.43<br/>4.49<br/>5.89 | 12.49 | 15.13<br/>14.57 | 11.18 | 5.94 | 10.99 |
| ttft P95 | 56.64<br/>47.87 | 24.60<br/>23.43<br/>33.39 | 43.62 | 96.25<br/>165.83 | 62.17 | 27.39 | 50.76 |
| Average inter token latency | 0.39<br/>0.40<br/> | 0.13<br/>0.13<br/>0.19 | 0.38 | 0.25<br/>0.28 | 0.34 | 0.15 | 0.36 |
| average token throughoutput per second | 14.30<br/>14.63 | 15.78<br/>15.65<br/>15.20 | 13.49 | 16.50<br/>15.87 | 16.24 | 14.89 | 15.98 |
| average end to end latency | 38.65<br/>37.04 | 32.30<br/>30.23<br/>31.58 | 36.39 | 41.81<br/>43.13 | 37.22 | 36.00 | 43.26 |
| average cache utilization | 57.9%<br/>55.5% | 72.8%<br/>61.1%<br/>67.3% | 51.2% | 53.1%<br/>44.8% | 49.9% | 51.0% | 54.9% |
| average queue time | 2.47s<br/>1.91s | 776ms<br/>775ms<br/>457ms | 1.91s | 1.23s<br/>859ms | 863ms | 394ms | 481ms |
| preemptions total |  | 720 |  |  |  | 880 | 824 |

## lite-multi-lora

Requests to llama2 model with 4 lora, with 40 connections to different lora

| | baseline | maxium | simple-queue | simple-kvcache | queue+kvcache | lora+queue | lora+kvcache |
| --- | --- | --- | --- | --- | --- | --- | --- |
| requests per minute | 96.34 | 95.77 | 107.34 | 109.50 | 112.34 | 99.27 | 107.55 |
| average time to first token (ttft) | 3.01 | 1.48 | 0.96 | 0.49 | 0.35 | 2.33 | 1.65 |
| ttft P95 | 18.36 | 8.64 | 16.39 | 0.63 | 0.57 | 13.38 | 9.64 |
| Average inter token latency | 0.12 | 0.11 | 0.08 | 0.05 | 0.05 | 0.13 | 0.09 |
| average token throughoutput per second | 16.75 | 17.09 | 16.93 | 17.94 | 17.89 | 14.74 | 16.37 |
| average end to end latency | 20.19 | 19.46 | 18.33 | 17.01 | 17.15 | 20.54 | 18.58 |
| average cache utilization | 59.2% | 68.3% | 72.4% | 72.7% | 74.7% | 65.9% | 65.8% |
| average queue time | 898ms | 464ms | 279ms | 74.3ms | 39ms | 700ms | 501ms |
| preemptions total | 832 | 959 | 923 | 416 | 524 | 1276 | 928 |




