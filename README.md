# Create PagerDuty schedule overrides from the CLI

### How to build
```
# Install into your $GOPATH/bin
go install github.com/jsadowyj/pd-override@latest

# Build locally in directory
git clone https://github.com/jsadowyj/pd-override pd-override

cd $_

make build

./pd-override
```

### Setup
- Generate an API Token on your account page: click top right, "My Profile" ->
  "User Settings" -> "API Access"
- Get the PagerDuty schedule ID for the schedule you would like to create
  overrides for.
  The easiest way to find it is to look at the URL in your browser while editing
  the schedule, e.g. click top menu "People" -> "On-Call Schedules" -> "Your
  schedule"

  **Example:** `https://<your-company>.pagerduty.com/schedules#<schedule_id>`
- Create the following file in your `$HOME` directory
`~/.pd.yml`
```
---
authtoken: <api_token>
schedule_id: <schedule_id>
```
- Alternatively, set the following two environment variables: `PD_API_KEY` and `PD_SCHEDULE_ID`.

Table for weekdays:

| **Weekday** 	| **Letter** 	|
|:-----------:	|:----------:	|
| Sunday      	| U          	|
| Monday      	| M          	|
| Tuesday     	| T          	|
| Wednesday   	| W          	|
| Thursday    	| R          	|
| Friday      	| F          	|
| Saturday    	| S          	|

### Examples
This command would create the following overrides:
- Monday and Wednesday between 09:00am-05:00pm
- Tuesday and Friday between 09:00am-06:00pm
- Thursday between 09:00am-1:00pm

```
pd-override M,W@0900-1700 T,F@0900-1800 R@0900-1300
```
