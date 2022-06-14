/* Select validator removed events */
CREATE VIEW launch_validator_removed AS
SELECT
	E.id AS event_id,
	trim('"' FROM A1.value ::text) AS address,
	A2.value::numeric AS launch_id,
	E.created_at
FROM
	event AS E
	INNER JOIN attribute AS A1 ON A1.event_id = E.id
	INNER JOIN attribute AS A2 ON A2.event_id = E.id
WHERE
	E.type = 'tendermint.spn.launch.EventValidatorRemoved'
	AND A1.name = 'genesisValidatorAccount'
	AND A2.name = 'launchID'
;

/* Select validator added events */
CREATE VIEW launch_validator_added AS
SELECT
	E.id AS event_id,
	trim('"' FROM A1.value ::text) AS address,
	A2.value::numeric AS launch_id,
	E.created_at
FROM
	event AS E
	INNER JOIN attribute AS A1 ON A1.event_id = E.id
	INNER JOIN attribute AS A2 ON A2.event_id = E.id
WHERE
	E.type = 'tendermint.spn.launch.EventValidatorAdded'
	AND A1.name = 'address'
	AND A2.name = 'launchID'
;

/* Select the added validators that were not removed */
CREATE VIEW launch_validator AS
SELECT
    A.event_id,
    A.address,
    A.launch_id,
    A.created_at
FROM
    launch_validator_added AS A
WHERE
    address NOT IN (
	    SELECT address
	    FROM launch_validator_removed
	    WHERE launch_id = A.launch_id AND created_at >= A.created_at
    )
;

/* Select chain created events */
CREATE VIEW launch_chain_created AS
SELECT
	E.id AS event_id,
	A1.value::numeric AS coordinator_id,
	trim('"' FROM A2.value ::text) AS coordinator_address,
	A3.value::numeric AS launch_id,
	E.created_at
FROM
	event AS E
	INNER JOIN attribute AS A1 ON A1.event_id = E.id
	INNER JOIN attribute AS A2 ON A2.event_id = E.id
	INNER JOIN attribute AS A3 ON A3.event_id = E.id
WHERE
	E.type = 'tendermint.spn.launch.EventChainCreated'
	AND A1.name = 'coordinatorID'
	AND A2.name = 'coordinatorAddress'
	AND A3.name = 'launchID'
;

/* Select campaign created events */
CREATE VIEW campaign_campaign_created AS
SELECT
	E.id AS event_id,
	A1.value::numeric AS coordinator_id,
	trim('"' FROM A2.value ::text) AS coordinator_address,
	A3.value::numeric AS campaign_id,
	E.created_at
FROM
	event AS E
	INNER JOIN attribute AS A1 ON A1.event_id = E.id
	INNER JOIN attribute AS A2 ON A2.event_id = E.id
	INNER JOIN attribute AS A3 ON A3.event_id = E.id
WHERE
	E.type = 'tendermint.spn.campaign.EventCampaignCreated'
	AND A1.name = 'coordinatorID'
	AND A2.name = 'coordinatorAddress'
	AND A3.name = 'campaignID'
;
