![go-lambda](http://cl.ly/1w1U1n3w3W2n/go-lamda-alt.png)

## go-lambda ![aws](https://d0.awsstatic.com/logos/aws/AWS_Logo_PoweredBy_127px.png)

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
