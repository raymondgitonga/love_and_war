# LOVE AND WAR : Give Your App Love By Unleashing War
## Simple load testing app to test your http services

![Gladiator](https://github.com/scraly/gophers/blob/main/gladiator-gopher.png)
### Installation

Build docker image: <code> docker build -t love_and_war . 

Run docker image: <code> docker run -t -p 8081:8081 love_and_war 

### Send attack request
```json
{
    "method": "POST",
    "url":"http://127.0.0.1:8080/test",
    "attack_duration":"5",
    "attack_rate":"10",
    "pay_load":"{\"userId\":1}",
    "pass_rate": 100
}
```
#### <code> method: </code> The http method for the endpoint you are testing

#### <code> url:  </code> The http url being attacked

#### <code> attack_duration:  </code> Duration for the attack in seconds

#### <code> attack_rate:  </code> Number of attacks to be performed per seconds

#### <code> pay_load:  </code> String representation of the payload to be sent to the url for Get endpoints send empty string. Convert your json to string [here](https://jsontostring.com/)

#### <code> pass_rate:  </code> Success rate percentage that should be met

### Attack response

```json

{
    "latencies": {
        "total": 0.061594849,
        "mean": 0.000615948
    },
    "duration": 10,
    "wait": 0.000352949,
    "requests": 100,
    "throughput": 10.099509294546035,
    "success": 100,
    "status_codes": {
        "200": 95,
        "500": 5
    },
    "pass": true
}

```

#### <code> latencies: </code> To understand more about latencies read [here](https://bravenewgeek.com/everything-you-know-about-latency-is-wrong/).

#### <code> duration:  </code> How long the attack took in minutes

#### <code> wait:  </code> Wait time before a response in seconds

#### <code> requests:  </code> Total number of requests that were sent

#### <code> throughput:  </code> Transactions per second

#### <code> success:  </code> Percentage success rate

#### <code> status_codes:  </code> Status codes received and their count

#### <code> pass:  </code> If pass rate provided in request was met


[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)





