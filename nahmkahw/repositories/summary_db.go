package repositories

func (r *reportRepoDB) FindReportSummary() (*ReportSummary ,error) {
	sql := `WITH combined AS (
SELECT 1 as rnum, 99999 code,'รายการทั้งหมด' as codename,COUNT(r.RECEIPT_ID) AS codecount
    FROM (
        SELECT f.receipt_id,f.code,o.DATE_SUCCESS
        FROM fees_receipt f
        INNER JOIN fees_order o ON (f.order_id = o.order_id
        AND o.STATUS_SUCCESS = 'SUCCESS'
        AND o.DATE_SUCCESS IS NOT NULL)
        WHERE f.std_code != '6299999991'  and f.code != 40
        and f.STATUS_OPERATE = 'SUCCESS'
    ) r
    
union all

select 2 as rnum, max_code code,max_name codename ,max_count codecount from (
SELECT r.code AS max_code,fee_name max_name,COUNT(r.RECEIPT_ID) AS report_count, MAX(COUNT(r.code)) OVER () AS max_count
    FROM (
        SELECT f.receipt_id,f.code,sh.fee_name, o.DATE_SUCCESS
        FROM fees_receipt f
        INNER JOIN fees_order o ON (f.order_id = o.order_id
        AND o.STATUS_SUCCESS = 'SUCCESS'
        AND o.DATE_SUCCESS IS NOT NULL)
        left join fees_sheet sh on F.CODE = sh.fee_no
        WHERE f.std_code != '6299999991' and f.code != 40
        and f.STATUS_OPERATE = 'SUCCESS'
    ) r
GROUP BY r.code, r.fee_name)
WHERE ROWNUM = 1 and max_count = report_count

union all

select 3 as rnum, min_code code,min_name codename,min_count codecount from (
SELECT r.code AS min_code,fee_name min_name,COUNT(r.RECEIPT_ID) AS report_count, MIN(COUNT(r.code)) OVER () AS min_count
    FROM (
        SELECT f.receipt_id,f.code,sh.fee_name, o.DATE_SUCCESS
        FROM fees_receipt f
        INNER JOIN fees_order o ON (f.order_id = o.order_id
        AND o.STATUS_SUCCESS = 'SUCCESS'
        AND o.DATE_SUCCESS IS NOT NULL)
        left join fees_sheet sh on F.CODE = sh.fee_no
        WHERE f.std_code != '6299999991' and f.code != 40
        and f.STATUS_OPERATE = 'SUCCESS'
    ) r
GROUP BY r.code, r.fee_name)
WHERE ROWNUM = 1 and min_count = report_count
)

SELECT
    MAX(CASE WHEN rnum = 1 THEN code END) AS summary_code,
    MAX(CASE WHEN rnum = 1 THEN codename END) AS summary_name,
    MAX(CASE WHEN rnum = 1 THEN codecount END) AS summary_count,
    
    MAX(CASE WHEN rnum = 2 THEN code END) AS max_code,
    MAX(CASE WHEN rnum = 2 THEN codename END) AS max_name,
    MAX(CASE WHEN rnum = 2 THEN codecount END) AS max_count,
    
    MAX(CASE WHEN rnum = 3 THEN code END) AS min_code,
    MAX(CASE WHEN rnum = 3 THEN codename END) AS min_name,
    MAX(CASE WHEN rnum = 3 THEN codecount END) AS min_count
FROM combined`

	var (
		report ReportSummary
	)

	row := r.oracle_db.QueryRow(sql)
	err := row.Scan(&report.SUMMARY_CODE, &report.SUMMARY_NAME, &report.SUMMARY_COUNT, &report.MAX_CODE,
	&report.MAX_NAME, &report.MAX_COUNT, &report.MIN_CODE, &report.MIN_NAME, &report.MIN_COUNT)

	if err != nil {
		return nil, err
	}

	return &report,nil
}

func (r *reportRepoDB) FindReportSuccess() (*ReportSuccess ,error) {
	sql := `WITH combined AS (
SELECT 1 as rnum, COUNT(r.RECEIPT_ID) AS value
FROM (
    SELECT f.receipt_id, f.code, o.DATE_SUCCESS
    FROM fees_receipt f
    INNER JOIN fees_order o 
        ON f.order_id = o.order_id
        AND o.STATUS_SUCCESS = 'SUCCESS'
        AND o.DATE_SUCCESS IS NOT NULL
    WHERE f.std_code != '6299999991'
      AND f.code != 40
      AND f.STATUS_OPERATE = 'SUCCESS'
      AND TRUNC(o.DATE_SUCCESS) = TRUNC(SYSDATE)
) r

union all

SELECT 2 as rnum, COUNT(r.RECEIPT_ID) AS value
FROM (
    SELECT f.receipt_id, f.code, o.DATE_SUCCESS
    FROM fees_receipt f
    INNER JOIN fees_order o 
        ON f.order_id = o.order_id
        AND o.STATUS_SUCCESS = 'SUCCESS'
        AND o.DATE_SUCCESS IS NOT NULL
    WHERE f.std_code != '6299999991'
      AND f.code != 40
      AND f.STATUS_OPERATE = 'SUCCESS'
      AND TRUNC(o.DATE_SUCCESS, 'MM') = TRUNC(SYSDATE, 'MM')
) r

union all

SELECT 3 as rnum, COUNT(r.RECEIPT_ID) AS value
FROM (
    SELECT f.receipt_id, f.code, o.DATE_SUCCESS
    FROM fees_receipt f
    INNER JOIN fees_order o 
        ON f.order_id = o.order_id
        AND o.STATUS_SUCCESS = 'SUCCESS'
        AND o.DATE_SUCCESS IS NOT NULL
    WHERE f.std_code != '6299999991'
      AND f.code != 40
      AND f.STATUS_OPERATE = 'SUCCESS'
      AND TRUNC(o.DATE_SUCCESS, 'YYYY') = TRUNC(SYSDATE, 'YYYY')
) r
)

SELECT
    MAX(CASE WHEN RNUM = 1 THEN VALUE END) AS SUCCESS_TODAY,
    MAX(CASE WHEN RNUM = 2 THEN VALUE END) AS SUCCESS_THIS_MONTH,
    MAX(CASE WHEN RNUM = 3 THEN VALUE END) AS SUCCESS_THIS_YEAR
FROM combined`

	var (
		report ReportSuccess
	)

	row := r.oracle_db.QueryRow(sql)
	err := row.Scan(&report.SUCCESS_TODAY, &report.SUCCESS_THIS_MONTH, &report.SUCCESS_THIS_YEAR)

	if err != nil {
		return nil, err
	}

	return &report,nil
}

func (r *reportRepoDB) FindReportCancel() (*ReportCancel ,error) {
	sql := `WITH combined AS (
SELECT 1 as rnum , COUNT(r.RECEIPT_ID) AS value
FROM (
    SELECT f.receipt_id, f.code, o.DATE_SUCCESS
    FROM fees_receipt f
    INNER JOIN fees_order o 
        ON f.order_id = o.order_id
        AND o.STATUS_SUCCESS = 'SUCCESS'
        AND o.DATE_SUCCESS IS NOT NULL
    WHERE f.std_code != '6299999991'
      AND f.code != 40
      AND f.STATUS_OPERATE = 'CANCEL'
      AND TRUNC(o.DATE_SUCCESS) = TRUNC(SYSDATE-4)
) r

union all

SELECT 2 as rnum , COUNT(r.RECEIPT_ID) AS value
FROM (
    SELECT f.receipt_id, f.code, o.DATE_SUCCESS
    FROM fees_receipt f
    INNER JOIN fees_order o 
        ON f.order_id = o.order_id
        AND o.STATUS_SUCCESS = 'SUCCESS'
        AND o.DATE_SUCCESS IS NOT NULL
    WHERE f.std_code != '6299999991'
      AND f.code != 40
      AND f.STATUS_OPERATE = 'CANCEL'
      AND TRUNC(o.DATE_SUCCESS, 'MM') = TRUNC(SYSDATE, 'MM')
) r

union all

SELECT 3 as rnum , COUNT(r.RECEIPT_ID) AS value
FROM (
    SELECT f.receipt_id, f.code, o.DATE_SUCCESS
    FROM fees_receipt f
    INNER JOIN fees_order o 
        ON f.order_id = o.order_id
        AND o.STATUS_SUCCESS = 'SUCCESS'
        AND o.DATE_SUCCESS IS NOT NULL
    WHERE f.std_code != '6299999991'
      AND f.code != 40
      AND f.STATUS_OPERATE = 'CANCEL'
      AND TRUNC(o.DATE_SUCCESS, 'YYYY') = TRUNC(SYSDATE, 'YYYY')
) r
)

SELECT
    MAX(CASE WHEN RNUM = 1 THEN VALUE END) AS CANCEL_TODAY,
    MAX(CASE WHEN RNUM = 2 THEN VALUE END) AS CANCEL_THIS_MONTH,
    MAX(CASE WHEN RNUM = 3 THEN VALUE END) AS CANCEL_THIS_YEAR
FROM combined
`

	var (
		report ReportCancel
	)

	row := r.oracle_db.QueryRow(sql)
	err := row.Scan(&report.CANCEL_TODAY, &report.CANCEL_THIS_MONTH, &report.CANCEL_THIS_YEAR)

	if err != nil {
		return nil, err
	}

	return &report,nil
}