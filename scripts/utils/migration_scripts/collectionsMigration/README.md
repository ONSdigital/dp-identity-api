Zebedee Collections Migration for SecAuth    

## Purpose
The objective of this module is to modify collection records teams entries for SecAuth.   
if migration environment variable is not set then this will revert the actions.

### Environment Variables
**teamsDir** the location of the zebedee collection teams files   
**collectionDir** the location of the zebedee collections   
**collectionCopyDir** the location of where collections to be copied before amending (this directory will be created if not existing)  
**migration** true for migration false for reversion


### How to run on remote environment ###
1) dp remote allow \< environment \>
2) go to .../dp-identity-api/scripts/utils/migration_scripts
    set the \< environment \> in the Makefile
    make all
    (this will copy the compiled code to the environment)
3) dp ssh \< environment \> publishing_mount 1
4)  `
export teamsDir=/var/florence/zebedee/teams/; \
export collectionDir=/var/florence/zebedee/collections/; \
export collectionCopyDir=~/copycollections20221006/; \
export migration=true/false;`
5) './bin-collection-migration/collection-migration'
6) make clean (to clear up afterward)

#### to run locally ####
set the environment variables  appropriately.
'run main.go'
