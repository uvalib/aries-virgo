# Aries Virdo

This is an implementation of the Aries API: https://confluence.lib.virginia.edu/display/DCMD/Aries for  Virgo

### System Requirements
* GO version 1.11 or greater (requires go mod)

### Current API

* GET /version : return service version info
* GET /healthcheck : test health of system components; results returned as JSON
* GET /api/aries/[ID] : Get information about objects contained in APTrust that match the ID
