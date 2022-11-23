# Learn GO

## Start app

```
docker compose up -d
```

## NSQ

[nsq.io](https://nsq.io/overview/design.html)
[official docker image](https://hub.docker.com/r/nsqio/nsq)
[go-nsq](https://github.com/nsqio/go-nsq)

### NSQ Admin

[nsq-admin](http://localhost:4171/)

### Publish Message

```
curl -d '{"Type":"Message Title","Status":"Created","Txid":"ytugdabd-4e5678yu9i-78q23hd2qdj","Amount":123,"Timestamp":"2022-04-02 23:23"}' 'http://localhost:4151/pub?topic=test_topic'
```
