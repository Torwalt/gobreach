SELECT
    email_hash,
    domain,
    breached_info,
    breach_date,
    breach_source
FROM
    Breach
WHERE
    email_hash = :email_hash
