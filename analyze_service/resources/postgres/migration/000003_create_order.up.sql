create table public.order
(
	id           uuid      not null,
	company_name text,
	status       text,
	tariff_id    uuid,
	company_info json,
	users        json,
	inn          text,
	created_at   timestamp not null,
	rejected_at  timestamp,
	success_at   timestamp,
	CONSTRAINT order_id_pk PRIMARY KEY ("id"),
	CONSTRAINT tariff_kind CHECK ( status in ('waiting', 'success', 'rejected'))
);