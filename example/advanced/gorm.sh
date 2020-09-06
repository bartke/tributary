#!/bin/bash

g () {
  sed "s/json:\"$1,omitempty\"/json:\"$1,omitempty\" gorm:\"$2\"/"
}

cat $1 \
| g "uuid" "primaryKey;type:char(36);" \
| g "selections" "foreignKey:BetUuid;references:Uuid;constraint:OnDelete:CASCADE" \
| g "bet_uuid" "type:char(36);" \
| g "id" "primaryKey;autoIncrement:true" \
| g "market" "type:varchar(255)" \
| g "stake" "type:decimal;" \
| g "exchange_rate" "type:decimal;" \
| g "odds" "type:decimal;" \
> $1.tmp && mv $1{.tmp,}

