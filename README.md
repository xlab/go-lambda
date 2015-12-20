## go-lambda ![aws](https://d0.awsstatic.com/logos/aws/AWS_Logo_PoweredBy_127px.png)

[![Join the chat at https://gitter.im/xlab/go-lambda](https://badges.gitter.im/xlab/go-lambda.svg)](https://gitter.im/xlab/go-lambda?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
![go-lambda](http://cl.ly/3a3V312h102e/go-lamda-gh.png)

### Example

```
$ go-lambda list
$ go-lambda create --role arn:aws:iam::account-id:role/lambda_basic_execution handler github.com/xlab/go-lambda/example

... make changes ...

$ go-lambda update example-handler handler github.com/xlab/go-lambda/example
$ go-lambda list
```

### License

MIT
