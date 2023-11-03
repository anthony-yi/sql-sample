# Sample SQS consumer and producer

## AWS Authentication

The code will use your session's active profile for AWA Authentication. So make sure to have your `aws configure` before running below code.

## Producer

### Code
```
producer\main.go
```

### Build
```
go -o produce ./producer/main.go
```

### Run

3 command line parameters:
* -queue-url = SQS Url (default: "")
* -num = Number of messages to produce/send to SQS (default: 10)
* -randomize-ids = Randomize both MessageGroupId and MessageDeduplicationId (default: false)

To generate 10 messages (not randomize-ids), run:
```
produce -queue-url <url>
```

To generate randomized 5 messages, run:
```
produce -queue-url <url> -num 5 -randomize-ids
```

## Consumer

### Code
```
consumer\main.go
```

### Build
```
go -o consume ./consumer/main.go
```

### Run

4 command line parameters:
* -queue-url = SQS Url (default: "")
* -threads = Number of consumers (default: 2)
* -max-messages = Max number of messages to receive (MaxNumberOfMessages in AWS Documenation) (default: 1)
* -wait-time-seconds = Wait time for message to arrive (default: 10)

To start with 2 consumers, run:
```
consume -queue-url <url> -threads 2
```

To start with 5 consumers and get 10 messages per consumer, run:
```
consume -queue-url <url> -threads 5 -max-messages 10
```

To start with default and wait 1 minute for message to arrive, run:
```
consume -queue-url <url> -wait-time-seconds 60
```
