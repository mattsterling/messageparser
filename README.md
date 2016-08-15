# Message Parser

A simple microservice for parsing a single chat message and providing a response of links (with titles), mentions,
and emoticons that are present in the message. This code was built and tested with Go 1.6.

# Getting Started:
Clone dependencies
`make deps`

# To setup the config for a sample run:
`source setup.sh`

# Run the server:
`make run`

# Install into your Go workspace:
`make install`

# Configuration
Configuration is done via Env Variables. The program expects the following:
* MESSAGE_PARSER_PORT - Port to run the server on.
* MESSAGE_PARSER_MSGSIZE - Max allowed content length of an HTTP body.

# API
* 'POST /v1/message' - Parses a message returning mentions, emojis, and links.
* 'GET  /health' - Shallow health check to indicate service is up.

# Potential Enhancements:
* URLs are expected to end with a white space. A url not ending in a white space but say a delimiter such as a ',' will result in a failed HTTP Get request and the title will not be retrieved.
* Internal parsing utilizes a basic payload scan in N time. This is largely due to not knowing byte sizes preventing a smarter way to segment the payload. Regex were avoided for implementation illustration but may be more useful in a production scenario.
* Parsing is done in one go routine. This could be split amongst more go routines, however under large load the go routines could become very large and induce latency on the machine.
* Retrieval is done via multiple go routines after parsing is complete. No restrictions are in place to the number of links allowed for a message, thus a high N of links will results in N go routines and larger API Latency.
* Links should be cached to avoid extra calls to find link titles.
