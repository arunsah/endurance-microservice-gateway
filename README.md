# endurance-microservice-gateway
---
endurance-microservice is collection of repository whose objective is to create important of microservice architecture; this repository aims to create load-balancer, service registrar, api gateway. Currently project is for learning purpose and not tested in production environment. Primary language is Golang along with Web front end.

## Primary Functionality and expectation:

- Multiple instances of these services will be added in A record in DNS to work them in DNS-round robin (for !Pv6).
- Will explore the feature of ```anycast``` (for IPv6)
- Will function as load balancer.
- Will function as reverse proxy of other (micro)services.
- Will function as service registrar.
- Will function as API gateway (throttling|limit).

Later these functionality may be broken down in independent applications if the code base increases.
