# airstrike
Load-testing or stress-testing an API can be imagined as an air strike, which
consists of:

    1. a Mission that describes when and what ordnance will be deployed
    2. a Squadron of planes that can simultaneously deploy their Arsenals
    3. an Arsenal on each plane consisting of Bombs and Missiles

The default mission (1) causes the configured squadron to simultaneously deploy
their configured arsenals every 5 seconds.

Some ordnance will hit its target and result in reportable damage (response time) sooner than others, and the output logging will reflect this, as the reports will typically arrive "out of order" as the concurrent weapon deployments finish.

**Bombs** are used for API transactions where all inputs can be known before
runtime. If you know the HTTP verb, URL, and (optional) payload ahead of time, you can use a Bomb.

**Missiles** are used for API transactions that depend on the current state of your
account, such as deleting the most-recently-created object. You provide them with a function that will be executed when the Missile is deployed.