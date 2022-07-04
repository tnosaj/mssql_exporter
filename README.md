https://docs.microsoft.com/en-us/sql/relational-databases/system-dynamic-management-views/sql-server-operating-system-related-dynamic-management-views-transact-sql?view=sql-server-ver16


Keeping a list of improvements:


statement_end_offset from - dm_exec_requests

```
Indicates, in bytes, starting with 0, the ending position of the currently executing statement for the currently executing batch or persisted object. Can be used together with the sql_handle, the statement_start_offset, and the sys.dm_exec_sql_text dynamic management function to retrieve the currently executing statement for the request. Is nullable.
```

Generallys - gather potential plan vs query stats.