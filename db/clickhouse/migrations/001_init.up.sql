create table item_service_logs
(
    Id         Int64,
    CampaignId Int64,
    Name       String,
    Description String,
    Priority   Int64,
    Removed    bool,
    EventTime  Int64
)
    engine = NATS
        SETTINGS
        nats_url = 'nats:4222',
        nats_subjects = 'logs',
        nats_format = 'JSONEachRow';

create table item_logs (
                           Id         Int64,
                           CampaignId Int64,
                           Name       String,
                           Description String,
                           Priority   Int64,
                           Removed    bool,
                           EventTime  Int64
) ENGINE = MergeTree() ORDER BY Id;

CREATE MATERIALIZED VIEW consumer TO item_logs
AS SELECT Id,  CampaignId, Name, Description, Priority, Removed, EventTime FROM item_service_logs;

SELECT * FROM item_logs ORDER BY Id;



