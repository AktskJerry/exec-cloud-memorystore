# exec-cloud-redis

This tool helps you to execute google cloud memorystore(only redis now) in cloud functions.
Usually we just want build minimum services for testing with redis, in the mean time memorystore  does not allow connnect from local, it's over kill to use GAE,CE,K8s also cost a little bit high in development phase. I suggest you connect with cloud functions and query redis with this tool, it makes you operate as local machine, additionally, it's almost free, serverless fee can be ignored!

Redis features are based on [redigo](https://github.com/gomodule/redigo) project

![image](https://github.com/AktskJerry/exec-cloud-memorystore/blob/master/example.png)

# How to use

  - Follow this [document](https://cloud.google.com/functions/docs/first-go) to create cloud functions with golang
  - Follow this [document](https://cloud.google.com/memorystore/docs/redis/creating-managing-instances?hl=zh-tw) to create redis instance
  - Follow this [document](https://cloud.google.com/memorystore/docs/redis/connect-redis-instance-functions) to connect to redis from cloud functions
  - Paste the code into source in cloud functions and deploy
  - Switch to testing tab, write your query. That's all

# Execution format
Json format with data as key, array is the value.
For example:
```sh
{ 
   "data": [{"do_or_send": "Do", "command": "KEYS", "params":"*"}]
}
```
`Do` is most common used, `command` is the redis execution command, params are the parameters you want to pass to redis. You can also execute without params:
```sh
{ 
   "data": [{"do_or_send": "Do", "command": "FLUSHDB"}]
}
```
More examples:
```sh
{ 
   "data": [{"do_or_send": "Do", "command": "ZREVRANGE", "params":"score,0,-1,withscores"}]
}
```
```sh
{ 
   "data": [{"do_or_send": "Do", "command": "smembers", "params":"friends:12345"}]
}
```
Transaction example:
```sh
{ 
   "data": [{"do_or_send": "Send", "command": "MULTI"},
            {"do_or_send": "Send", "command": "SET", "params":"KeyA,1"},
            {"do_or_send": "Send", "command": "SET", "params":"KeyB,2"},
            {"do_or_send": "Do", "command": "EXEC"}]
}
```

Not tested:
  - Redis Sentinel
  - Redis Cluster

License
----

MIT
