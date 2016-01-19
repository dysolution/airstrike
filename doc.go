/*
Package airstrike conducts workflow and/or stress-testing against an API.

Load-testing or stress-testing an API can be imagined as an air strike, which
consists of:

    1. a Mission that describes when and what ordnance will be deployed
    2. a Squadron of planes that can simultaneously deploy their Arsenals
    3. an Arsenal on each plane consisting of Bombs and Missiles

The default Mission simultaneously commands each plane in the configured
Squadron to deploy its configured arsenal every 5 seconds.  Some ordnance will
hit its target and result in reportable damage (response time) sooner than
others, and the output logging will reflect this, as the reports will typically
arrive "out of order" as the concurrent weapon deployments finish.

Bombs are used for API transactions where all inputs can be known before
runtime. If you know the HTTP verb, URL, and (optional) payload ahead of time,
you can use a Bomb. Bombs assume the state of objects associated with your
account and as such must be trusted to hit their target. If they "miss" with a
404 or other 4xx, this will be logged as an error.

Missiles are used for API transactions that depend on the current state of your
account, such as deleting the most-recently-created object. You provide them
with a function that will be executed when the Missile is deployed. Unlike a
Bomb, Missiles need to be "guided" in this fashion to ensure they will hit
their target. As with Bombs, 4xx responses are reported as errors.
*/
package airstrike
