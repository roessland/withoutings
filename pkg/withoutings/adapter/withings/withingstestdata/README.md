# Withings test data

Responses from Withings.

## Sleep v2 - Getsummary

Webhook: 
```
userid=13371337&startdate=1709336580&enddate=1709368500&appli=44
```

withings-resp-sleepv2-getsummary-0.json
```shell
curl -v -X POST -H "Authorization: Bearer $WITHINGS_ACCESS_TOKEN" -d 'action=getsummary&data_fields=nb_rem_episodes%2Csleep_efficiency%2Csleep_latency%2Ctotal_sleep_time%2Ctotal_timeinbed%2Cwakeup_latency%2Cwaso%2Capnea_hypopnea_index%2Cbreathing_disturbances_intensity%2Casleepduration%2Cdeepsleepduration%2Cdurationtosleep%2Cdurationtowakeup%2Chr_average%2Chr_max%2Chr_min%2Clightsleepduration%2Cnight_events%2Cout_of_bed_count%2Cremsleepduration%2Crr_average%2Crr_max%2Crr_min%2Csleep_score%2Csnoring%2Csnoringepisodecount%2Cwakeupcount%2Cwakeupduration&enddateymd=2023-05-23&startdateymd=2023-05-22' https://wbsapi.withings.net/v2/sleep | jq > withings-resp-sleepv2-getsummary-0.json
```

## Sleep v2 - Get

Webhook (same as Getsummary): 
```
userid=13371337&startdate=1709336580&enddate=1709368500&appli=44
```

```
curl -v -X POST -H "Authorization: Bearer $WITHINGS_ACCESS_TOKEN" -d 'action=get&data_fields=hr,rr,snoring,sdnn_1,rmssd,mvt_score&enddate=1709368500&startdate=1709336580' https://wbsapi.withings.net/v2/sleep | jq > withdata.json
```

## Weight - Getmeas

Webhook:
```
userid=13371337&startdate=1706376729&enddate=1706376730&appli=1
```

```
curl -v -X POST -H "Authorization: Bearer $WITHINGS_ACCESS_TOKEN" -d 'action=get&data_fields=hr,rr,snoring,sdnn_1,rmssd,mvt_score&enddate=1709368500&startdate=1709336580' https://wbsapi.withings.net/v2/sleep | jq > withdata.json
```