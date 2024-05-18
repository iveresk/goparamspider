# goparamspider
**!Disclaimer**: this is not `paramspider` or `sqlmap` wrapper it's a standalone Red Team enumeration tool, made due to unsuitable results for enterprise projects.  

## Usage
1. Fill JSON by keywords in `./assets/payloads.json` file. Consider two modes `day` (critical checks, fuzzing, could be built-into PR pipeline). And a `night` mode for a scheduled pipelines. Now it takes 16-18 hours to check the default ones I pushed to the repo.  
```
{
  "mode" : [{
    "day" : [{
      "GET" : [{
        "routes" : [
          "api/auth/",
          "api/auth/guest/",
          "api/auth/login/",
          "api/grafana/"
        ],
        "parameters" : [
          "v",
          "query",
          "log",
          "search",
          "execute"
        ],
        "payloads" : [
          "%3D%3D",
          "%3D",
          "%27",
          "%27%20%2D%2D",
          "%27%20%23",
          "%27%20â",
          "%27%2D%2D",
          "%27%2F%2A",
          "%27%23",
          "%22%20%2D%2D",
          "%20or%201%3D1â"
        ]
... POST, and other methods
```  
2. Enumeration run  
So, the necessity in my own enumeration tool was held due to unsetisfied results of other tools (see the Disclaimer), and poor functionality on uploading my own payloads, formats, parameters.  
So, as a necessity to automate my work and experience, and make it in a quality way you can deliver to a customer there was made a "hackathon" to create this enumeration tool in 4 days.  
```
 Flags:
 -m Mode of the tool usage defining the if it is day (PR) mode or night (Full) scan.
 -u Target Domain.
 -d - The delay between requests not to be blocked by WAF. Default value is 1000ms
 -l - The count of params to be tested combined in line.
 - f - Flag to set output to the logging file /var/log/syslog.
 - v - Flag to set verbose flag and record all debugging and rejected requests.
Example:
 ./goparamspider -m day -u domain.com
```
## Plans for the nearest future
 * Golang 1.22.3 log/v2 adaptation.

