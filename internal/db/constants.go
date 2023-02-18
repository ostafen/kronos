package db

type Row interface {
	Scan(...any) error
}

const (
	schedTableName = "schedules"
)

const (
	schedTableIdCol             = "id"
	schedTableTitleCol          = "title"
	schedTableActiveCol         = "active"
	schedTableDescriptionCol    = "description"
	schedTableUrlCol            = "url"
	schedTableCronExprCol       = "cron_expr"
	schedTableEmailCol          = "email"
	schedTableNextScheduleAtCol = "next_schedule_at"
	schedTableCreatedAtCol      = "created_at"
	schedTableMetadataCol       = "metadata"
	schedTableFailuresCol       = "failures"
)

var schedTableCols = []string{
	schedTableIdCol,
	schedTableTitleCol,
	schedTableActiveCol,
	schedTableDescriptionCol,
	schedTableUrlCol,
	schedTableCronExprCol,
	schedTableEmailCol,
	schedTableCreatedAtCol,
	schedTableNextScheduleAtCol,
	schedTableMetadataCol,
	schedTableFailuresCol,
}
